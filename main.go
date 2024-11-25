package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/config"
	"github.com/hasib-003/newsLetter/routes"
)

func main() {

	config.ConnectDB()

	r := gin.Default()

	routes.RegisterUserRoutes(r)
	routes.RegisterNewsRoutes(r)

	r.Run(":8080")
}
