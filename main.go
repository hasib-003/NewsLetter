package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/config"
	"github.com/hasib-003/newsLetter/routes"
)

func main() {
	// Connect to PostgreSQL
	config.ConnectDB()

	// Set up Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterUserRoutes(r)
	routes.RegisterNewsRoutes(r)

	// Start the server
	r.Run(":8080")
}
