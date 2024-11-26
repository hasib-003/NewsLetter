package controllers

import (
	"net/http"
	"strconv"

	"github.com/hasib-003/newsLetter/models"
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

func (nc *NewsController) FetchNewsByTopic(c *gin.Context) {
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
func (nc *NewsController) GetNewsByTopicID(c *gin.Context) {
	topicID := c.Param("id")

	var id uint
	if parsedID, err := strconv.ParseUint(topicID, 10, 32); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	} else {
		id = uint(parsedID)
	}

	topics, err := nc.NewsService.GetNewsByTopicID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
		return
	}

	if len(topics) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No news found for the given topic ID"})
		return
	}

	c.JSON(http.StatusOK, topics)
}
func (nc *NewsController) SubscribeToTopic(c *gin.Context) {
	var request struct {
		UserID    uint   `json:"user_id"`
		TopicName string `json:"topic_name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := nc.NewsService.SubscribeUserToTopic(request.UserID, request.TopicName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully subscribed to topic"})
}
func (nc *NewsController) GetSubscribedTopics(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var subscriptions []models.Subscription
	err := nc.NewsService.DB.Where("user_id = ?", request.UserID).Find(&subscriptions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching subscriptions"})
		return
	}

	var topics []models.Topic
	for _, sub := range subscriptions {
		var topic models.Topic
		err := nc.NewsService.DB.First(&topic, sub.TopicID).Error
		if err != nil {
			continue
		}
		topics = append(topics, topic)
	}

	c.JSON(http.StatusOK, gin.H{"subscribed_topics": topics})
}

func (nc *NewsController) SendEmails(c *gin.Context) {
	users, err := nc.NewsService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
	}
	news, err := nc.NewsService.FetchNewsByTopic("technology")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
	}

	status := nc.NewsService.SendEmails(users, news)
	c.JSON(http.StatusOK, gin.H{"message": "Emails are being sent", "status": status})
}
