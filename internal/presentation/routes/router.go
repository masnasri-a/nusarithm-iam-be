package routes

import (
	"database/sql"

	"backend/internal/application/services"
	"backend/internal/infrastructure/repositories"
	"backend/internal/presentation/handlers"

	_ "backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	// Initialize repositories
	domainRepo := repositories.NewDomainRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	domainService := services.NewDomainService(domainRepo)
	roleService := services.NewRoleService(roleRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, roleRepo, domainRepo, "your-secret-key") // TODO: Use environment variable for secret

	// Initialize handlers
	domainHandler := handlers.NewDomainHandler(domainService)
	roleHandler := handlers.NewRoleHandler(roleService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware - allow all origins, support credentials
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "X-NRM-DID", "X-Nrm-Did"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,     // Credentials cannot be used with AllowOrigins: ["*"]
		MaxAge:           12 * 3600, // 12 hours
	}))

	// Ping endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Handle OPTIONS requests for all routes
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-NRM-DID")
		c.Header("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours
		c.Status(200)
	})

	// Role routes (must come before domain routes to avoid path conflicts)
	r.GET("/roles", roleHandler.ListRoles)
	r.GET("/roles/:id", roleHandler.GetRole)
	r.GET("/domains/:domainId/roles", roleHandler.GetRolesByDomain)
	r.POST("/domains/:domainId/roles", roleHandler.CreateRole)
	r.PUT("/roles/:id", roleHandler.UpdateRole)
	r.DELETE("/roles/:id", roleHandler.DeleteRole)

	// User routes
	r.GET("/users", userHandler.ListUsers)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/users/:id/reset-password", userHandler.ResetUserPassword)
	r.GET("/domains/:domainId/users", userHandler.GetUsersByDomain)
	r.POST("/users", userHandler.CreateUser)
	r.PUT("/users/:id", userHandler.UpdateUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	// Auth routes
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/validate", authHandler.ValidateToken)
	r.GET("/auth/profile", authHandler.GetProfile)

	// Domain routes
	r.GET("/domains", domainHandler.ListDomains)
	r.GET("/domains/:domainId", domainHandler.GetDomain)
	r.POST("/domains", domainHandler.CreateDomain)
	r.PUT("/domains/:domainId", domainHandler.UpdateDomain)
	r.DELETE("/domains/:domainId", domainHandler.DeleteDomain)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
