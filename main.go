// Package main Nusarithm IAM Backend API
//
//	@title			Nusarithm IAM API
//	@version		1.0
//	@description	This is the API for Nusarithm IAM Backend
//	@host			localhost:8080
//	@BasePath		/
package main

import (
	"log"

	"backend/internal/infrastructure/config"
	"backend/internal/presentation/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	// Initialize database config
	dbConfig := config.NewDatabaseConfig()
	db, err := dbConfig.OpenDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Setup router
	r := routes.SetupRouter(db)

	log.Fatal(r.Run(":8080"))
}
