package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/config"
	"github.com/hasib-003/newsLetter/controllers"
	"github.com/hasib-003/newsLetter/services"
)

func RegisterNewsRoutes(r *gin.Engine) {
	db := config.GetDB() 
	newsService := services.NewNewsService(db)
	newsController := controllers.NewNewsController(newsService)

	r.GET("/news", newsController.FetchNewsByTopic)
	r.GET("/send-email", newsController.SendEmails)
	r.GET("/topics/:id", newsController.GetNewsByTopicID)
	r.GET("/subscribed-topics/:user_id", newsController.GetSubscribedTopics)

}
