package services

import (
	"log"

	"github.com/hasib-003/newsLetter/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
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
