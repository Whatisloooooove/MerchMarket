package main

import (
	"merch_service/internal/db"
	"merch_service/internal/handlers"
	"merch_service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db.Init()

	// Create a Gin router
	r := gin.Default()

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/login", handlers.LoginHandler)
	}

	// Protected routes (JWT authentication required)
	auth := api.Group("/")
	auth.Use(middlewares.JWTAuth())
	{
		auth.GET("/me", handlers.MeHandler) // Returns current user info
	}

	// Start the server on port 8080
	r.Run(":8080")
}
