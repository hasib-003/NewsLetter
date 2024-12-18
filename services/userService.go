package services

import (
	"fmt"
	"log"

	"github.com/hasib-003/newsLetter/models"
	"gorm.io/gorm"
)

type UserService struct{DB *gorm.DB}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (us *UserService) RegisterUser(email, name, password string) (models.User, error) {

	user, err := models.CreateUser(email, name, password)
	if err != nil {
		log.Println("Error in service layer:", err)
		return user, err
	}
	return user, nil
}
func (us *UserService) GetAllUsers() ([]models.User, error) {
	users, err := models.GetAllUsers()
	if err != nil {
		log.Println("Error in service layer:", err)
		return nil, err
	}
	return users, nil
}
func (us *UserService) GetAUser(email string)(models.User,error){
	user,err :=models.GetAUser(email)
	if err!=nil{
		log.Println("Error in service layer ",err)
		return user,err
	}
	return user,nil
}
func (us *UserService) SubscribeUserToTopic(userID uint, topicName string) error {

	var topics []models.Topic
	err := us.DB.Where("name = ?", topicName).Find(&topics).Error
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
		if err := us.DB.Create(&subscription).Error; err != nil {
			return fmt.Errorf("error creating subscription for topic %s: %v", topic.Name, err)
		}
	}

	return nil
}
func (us *UserService) UnsubscribeUserFromTopic(userID uint, topicName string) error {

	var topics []models.Topic
	err := us.DB.Where("name = ?", topicName).Find(&topics).Error
	if err != nil {
		return fmt.Errorf("error fetching topics by name: %v", err)
	}

	if len(topics) == 0 {
		return fmt.Errorf("no topics found with the name: %s", topicName)
	}

	for _, topic := range topics {
		err := us.DB.Where("user_id = ? AND topic_id = ?", userID, topic.ID).Delete(&models.Subscription{}).Error
		if err != nil {
			return fmt.Errorf("error unsubscribing from topic %s: %v", topic.Name, err)
		}
	}

	return nil
}

