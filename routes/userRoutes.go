package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/services"
	"github.com/hasib-003/newsLetter/controllers"
)

func RegisterUserRoutes(r *gin.Engine) {

	userService := services.NewUserService()
	userController := controllers.NewUserController(userService)

	r.POST("/register", userController.RegisterUser)
}
