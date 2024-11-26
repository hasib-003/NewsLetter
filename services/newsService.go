package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/hasib-003/newsLetter/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type NewsService struct{DB *gorm.DB}

func NewNewsService(db *gorm.DB) *NewsService {
	return &NewsService{DB: db}
}

func (ns *NewsService) FetchNewsByTopic(topic string)  ([]models.Article, error) {

	err := godotenv.Load()
	if err != nil {
		return nil,errors.New("failed to load environment variables")
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		return  nil,errors.New("API key is missing")
	}

	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&apiKey=%s", topic, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return  nil,errors.New("failed to fetch news")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response: %v", err)
		return nil,errors.New("failed to read API response")
	}

	var newsResponse models.NewsResponse
	if err := json.Unmarshal(body, &newsResponse); err != nil {
		log.Printf("Error parsing news response: %v", err)
		return nil,errors.New("failed to parse news response")
	}

	for i, article := range newsResponse.Articles {
		if i>5{
			break
		}
		topicEntry := models.Topic{
			Name:        topic,
			Description: article.Description,
		}
		if err := ns.DB.Create(&topicEntry).Error; err != nil {
			log.Printf("Error saving topic to database: %v", err)
			return nil,err
		}
		
	}

	return  newsResponse.Articles,nil
}
func (ns *NewsService) GetNewsByTopicID(topicID uint) ([]models.Topic, error) {
	var topics []models.Topic
	err := ns.DB.Where("id = ?", topicID).Find(&topics).Error
	if err != nil {
		log.Printf("Error fetching news by topic ID: %v", err)
		return nil, err
	}
	return topics, nil
}
func (ns *NewsService) GetTopicsByName(topicName string) ([]models.Topic, error) {
	var topics []models.Topic
	err := ns.DB.Where("name = ?", topicName).Find(&topics).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching topics: %v", err)
	}
	return topics, nil
}
// func (ns *NewsService) SubscribeUserToTopic(userID uint, topicName string) error {
// 	var topic models.Topic
// 	err := ns.DB.Where("name = ?", topicName).First(&topic).Error
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		return fmt.Errorf("error checking for topic: %v", err)
// 	}

// 	subscription := models.Subscription{
// 		UserID:  userID,
// 		TopicID: topic.ID,
// 	}
// 	if err := ns.DB.Create(&subscription).Error; err != nil {
// 		return fmt.Errorf("error creating subscription: %v", err)
// 	}
// 	return nil
// }
func (ns *NewsService) SubscribeUserToTopic(userID uint, topicName string) error {
	// Find all topics with the same name
	var topics []models.Topic
	err := ns.DB.Where("name = ?", topicName).Find(&topics).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error fetching topics by name: %v", err)
	}

	// If no topics found, return an error
	if len(topics) == 0 {
		return fmt.Errorf("no topics found with the name: %s", topicName)
	}

	// Create subscription for each topic found
	for _, topic := range topics {
		subscription := models.Subscription{
			UserID:  userID,
			TopicID: topic.ID,
		}
		if err := ns.DB.Create(&subscription).Error; err != nil {
			return fmt.Errorf("error creating subscription for topic %s: %v", topic.Name, err)
		}
	}

	return nil
}

func (ns *NewsService) SendEmails(users []models.User, news []models.Article) map[string]string {
	status := make(map[string]string)
	statusCh := make(chan map[string]string, len(users))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()

			body := "Here are the top news articles:\n\n"
			for i, article := range news {
				if i >= 5 {
					break
				}
				body += fmt.Sprintf("Title: %s\nDescription: %s\n", article.Title, article.Description)

			}
			err := SendEmail(email, "Weekly Newsletter", body)
			statusUpdate := make(map[string]string)
			if err != nil {
				log.Printf("Error sending email to %s: %v", email, err)
				statusUpdate[email] = "Failed"
			} else {
				statusUpdate[email] = "Success"
			}
			statusCh <- statusUpdate
		}(user.Email)
	}
	go func() {

		wg.Wait()
		close(statusCh)
	}()
	for update := range statusCh {
		mu.Lock()
		for k, v := range update {
			status[k] = v
		}
		mu.Unlock()
	}

	log.Println("Email sending status:")
	for email, emailStatus := range status {
		log.Printf("Email: %s, Status: %s\n", email, emailStatus)
	}
	return status
}

func SendEmail(to, subject, body string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	from := os.Getenv("SENDER_EMAIL")
	client := sendgrid.NewSendClient(apiKey)
	message := mail.NewSingleEmail(
		mail.NewEmail("Newsletter", from),
		subject,
		mail.NewEmail("User", to),
		body,
		body,
	)

	log.Printf("Sending email to: %s\nSubject: %s\nBody: %s\nFrom: %s\n", to, subject, body, from)
	_, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}
	return nil
}
func (ns *NewsService) GetUsers() ([]models.User, error) {
	userService := NewUserService()
	return userService.GetAllUsers()
}
