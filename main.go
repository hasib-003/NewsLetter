package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/config"
	"github.com/hasib-003/newsLetter/routes"
	"github.com/hasib-003/newsLetter/scheduler"
)

func main() {

	config.ConnectDB()

	r := gin.Default()

	routes.RegisterUserRoutes(r)
	routes.RegisterNewsRoutes(r)

	scheduler.StartScheduler()

	r.Run(":8080")
}
