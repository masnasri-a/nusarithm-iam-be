package handlers

import (
	"net/http"
	"strings"

	"backend/internal/application/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      struct {
			ID          string                 `json:"id"`
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			Claims      map[string]interface{} `json:"claims"`
		} `json:"role"`
		Domain struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"domain"`
	} `json:"user"`
}

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			X-NRM-DID	header		string				true	"Domain ID"
//	@Param			credentials	body		LoginRequest		true	"Login credentials"
//	@Success		200			{object}	AuthResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	domainIdStr := c.GetHeader("X-NRM-DID")
	if domainIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-NRM-DID header is required"})
		return
	}

	domainID, err := uuid.Parse(domainIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain UUID in X-NRM-DID header"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResp, err := h.authService.Login(domainID, req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	response := AuthResponse{
		Token: loginResp.AccessToken,
	}
	response.User.ID = loginResp.User.ID.String()
	response.User.Username = loginResp.User.Username
	response.User.Email = loginResp.User.Email
	response.User.FirstName = loginResp.User.FirstName
	response.User.LastName = loginResp.User.LastName
	response.User.Role.ID = loginResp.User.Role.ID.String()
	response.User.Role.Name = loginResp.User.Role.Name
	response.User.Role.Description = loginResp.User.Role.Description
	response.User.Role.Claims = loginResp.User.Role.Claims
	response.User.Domain.ID = loginResp.User.Domain.ID.String()
	response.User.Domain.Name = loginResp.User.Domain.Name
	response.User.Domain.Description = loginResp.User.Domain.Description

	c.JSON(http.StatusOK, response)
}

// ValidateToken godoc
//
//	@Summary		Validate JWT token
//	@Description	Validate JWT token and return user information
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer token"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		401				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/auth/validate [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":  true,
		"claims": claims,
	})
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get authenticated user's profile information
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer token"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		401				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	// Validate token and get claims
	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Get user profile using user ID from token
	user, err := h.authService.GetProfile(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	profile := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role": map[string]interface{}{
			"id":          user.Role.ID,
			"name":        user.Role.Name,
			"description": user.Role.Description,
			"claims":      user.Role.Claims,
		},
		"domain": map[string]interface{}{
			"id":          user.Domain.ID,
			"name":        user.Domain.Name,
			"description": user.Domain.Description,
		},
	}

	c.JSON(http.StatusOK, profile)
}
