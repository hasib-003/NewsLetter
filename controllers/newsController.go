package controllers

import (

	"net/http"

	"github.com/hasib-003/newsLetter/services"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	NewsService *services.NewsService
}

func NewNewsController(newsService *services.NewsService) *NewsController {
	return &NewsController{
		NewsService: newsService,
	}
}

func (nc *NewsController) GetNews(c *gin.Context) {
	topic := c.Query("topic") 

	if topic == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Topic is required"})
		return
	}

	articles, err := nc.NewsService.FetchNewsByTopic(topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"topic":    topic,
		"articles": articles,
	})
}

func (nc *NewsController) SendEmails(c *gin.Context){
	users,err :=nc.NewsService.GetUsers()
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch users"})
	}
	news,err :=nc.NewsService.FetchNewsByTopic("technology")
		if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch news"})
	}
	// for i := 0; i < 5; i++ {
	// 	log.Printf("Title: %s\nDescription: %s\nURL: %s\n\n", news[i].Title, news[i].Description, news[i].URL)
	// }

	status:= nc.NewsService.SendEmails(users,news)
	c.JSON(http.StatusOK, gin.H{"message": "Emails are being sent","status":status})
}
