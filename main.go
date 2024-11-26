package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/config"
	"github.com/hasib-003/newsLetter/models"
	"github.com/hasib-003/newsLetter/routes"
	"github.com/hasib-003/newsLetter/scheduler"
)

func main() {

	config.ConnectDB()
	err := config.DB.AutoMigrate(&models.User{}, &models.Topic{}, &models.Subscription{})
	if err !=nil{
		log.Println("Failed to migrate ")
	}

	r := gin.Default()

	routes.RegisterUserRoutes(r)
	routes.RegisterNewsRoutes(r)

	scheduler.StartScheduler()

	r.Run(":8080")
}
