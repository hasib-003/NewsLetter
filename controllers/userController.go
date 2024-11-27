package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hasib-003/newsLetter/services"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	var userInput struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.UserService.RegisterUser(userInput.Email, userInput.Name, userInput.Password)
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

func (uc *UserController) GetAllUsers(c *gin.Context){
	users,err:=uc.UserService.GetAllUsers()
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch users"})
	}
	c.JSON(http.StatusOK,users)
}
func (uc *UserController) GetAUser(c *gin.Context){
	email:=c.Query("email")
	if email==""{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Email is required"})
	}
	user,err:=uc.UserService.GetAUser(email)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch User "})
	}
	c.JSON(http.StatusOK,user)
}
func (uc *UserController) SubscribeToTopic(c *gin.Context) {
	var request struct {
		UserID    uint   `json:"user_id"`
		TopicName string `json:"topic_name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := uc.UserService.SubscribeUserToTopic(request.UserID, request.TopicName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully subscribed to topic"})
}
