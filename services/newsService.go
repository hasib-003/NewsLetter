package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/hasib-003/newsLetter/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/joho/godotenv"
)

type NewsService struct{}

func NewNewsService() *NewsService {
	return &NewsService{}
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

	var newsResponse models.NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsResponse); err != nil {
		return nil, errors.New("failed to parse news response")
	}

	return newsResponse.Articles, nil
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
				body += fmt.Sprintf("Title: %s\nDescription: %s\nURL: %s\n\n", article.Title, article.Description, article.URL)

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
