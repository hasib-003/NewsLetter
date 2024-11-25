package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/services"
	"net/http"
)

type UserController struct {
	UserService *services.UserService
}


func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (ctrl *UserController) RegisterUser(c *gin.Context) {
	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.UserService.RegisterUser(userInput.Email, userInput.Name, userInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}
