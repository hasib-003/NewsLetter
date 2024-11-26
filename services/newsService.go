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

type NewsService struct{ DB *gorm.DB }

func NewNewsService(db *gorm.DB) *NewsService {
	return &NewsService{DB: db}
}

func (ns *NewsService) FetchNewsByTopic(topic string) ([]models.Article, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("failed to load environment variables")
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		return nil, errors.New("API key is missing")
	}

	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&apiKey=%s", topic, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to fetch news")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response: %v", err)
		return nil, errors.New("failed to read API response")
	}

	var newsResponse models.NewsResponse
	if err := json.Unmarshal(body, &newsResponse); err != nil {
		log.Printf("Error parsing news response: %v", err)
		return nil, errors.New("failed to parse news response")
	}

	for i, article := range newsResponse.Articles {
		if i > 5 {
			break
		}
		topicEntry := models.Topic{
			Name:        topic,
			Description: article.Description,
		}
		if err := ns.DB.Create(&topicEntry).Error; err != nil {
			log.Printf("Error saving topic to database: %v", err)
			return nil, err
		}

	}
	

	return newsResponse.Articles, nil
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

func (ns *NewsService) SubscribeUserToTopic(userID uint, topicName string) error {

	var topics []models.Topic
	err := ns.DB.Where("name = ?", topicName).Find(&topics).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error fetching topics by name: %v", err)
	}

	if len(topics) == 0 {
		return fmt.Errorf("no topics found with the name: %s", topicName)
	}

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
func (ns *NewsService) GetSubscribedTopics(userID uint) ([]models.Topic, error) {
	
	var subscriptions []models.Subscription
	err := ns.DB.Where("user_id = ?", userID).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching subscriptions: %v", err)
	}

	var topics []models.Topic
	for _, sub := range subscriptions {
		var topic models.Topic
		err := ns.DB.First(&topic, sub.TopicID).Error
		if err != nil {
			continue
		}
		topics = append(topics, topic)
	}

	return topics, nil
}


func (ns *NewsService) SendEmails(users []models.User) map[string]string {
	status := make(map[string]string)
	statusCh := make(chan map[string]string, len(users))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user models.User) {
			defer wg.Done()

			var subscriptions []models.Subscription
			err := ns.DB.Where("user_id = ?", user.ID).Find(&subscriptions).Error
			if err != nil {
				log.Printf("Error fetching subscriptions for user %s: %v", user.Email, err)
				statusCh <- map[string]string{user.Email: "Failed to fetch subscriptions"}
				return
			}

			var userNews []models.Topic
			for _, subscription := range subscriptions {
				var topic models.Topic
				err := ns.DB.First(&topic, subscription.TopicID).Error
				if err != nil {
					log.Printf("Error fetching topic for subscription: %v", err)
					continue
				}
				userNews = append(userNews, topic)
				// for _, article := range userNews {
				// 	if article.Name == topic.Name {
				// 		userNews = append(userNews, article)
				// 	}
				// }
			}
			// 	if len(userNews) > 5 {
			// 	userNews = userNews[:5]
			// }
			log.Printf("User %s subscribed to these topics: %v", user.Email, userNews)

			body := "Here are the top news articles:\n\n"
			for i, article := range userNews {
				if i>20{
					break
				}
				body += fmt.Sprintf("Title: %s\nDescription: %s\n", article.Name, article.Description)
				

			}
			err = SendEmail(user.Email, "Weekly Newsletter", body)
			statusUpdate := make(map[string]string)
			if err != nil {
				log.Printf("Error sending email to %s: %v", user.Email, err)
				statusUpdate[user.Email] = "Failed"
			} else {
				statusUpdate[user.Email] = "Success"
			}
			statusCh <- statusUpdate
		}(user)
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
