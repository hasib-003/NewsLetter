package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/config"
	"github.com/hasib-003/newsLetter/controllers"
	"github.com/hasib-003/newsLetter/services"
)

func RegisterUserRoutes(r *gin.Engine) {
	db := config.GetDB() 
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	r.POST("/register", userController.RegisterUser)
	r.GET("/getAllUsers",userController.GetAllUsers)
	r.GET("/getAUser",userController.GetAUser)
	r.POST("/subscribe",userController.SubscribeToTopic)
	r.POST("/unsubscribe", userController.UnsubscribeFromTopic)

}
