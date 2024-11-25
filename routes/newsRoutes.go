package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/controllers"
	"github.com/hasib-003/newsLetter/services"
)

func RegisterNewsRoutes(r *gin.Engine) {

	newsService := services.NewNewsService()
	NewsController := controllers.NewNewsController(newsService)

	r.GET("/news", NewsController.GetNews)
	r.GET("/send-email", NewsController.SendEmails)
}
