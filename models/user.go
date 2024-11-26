package models

import (
	"log"

	"github.com/hasib-003/newsLetter/config"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Name     string
	Password string
}

type Topic struct {
	gorm.Model
	Name        string 
	Description string
}

type Subscription struct {
	gorm.Model
	UserID  uint
	TopicID uint
	 User    User   `gorm:"foreignKey:UserID"`
	Topic   Topic  `gorm:"foreignKey:TopicID"`
}

func CreateUser(email, name, password string) (User, error) {

	user := User{
		Email:    email,
		Name:     name,
		Password: password,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		log.Println("Error inserting user:", err)
		return user, err
	}
	return user, nil
}
func GetAllUsers() ([]User, error) {
	var users []User
	if err := config.DB.Find(&users).Error; err != nil {
		log.Println("Error fetching users:", err)
		return nil, err
	}
	return users, nil
}
func GetAUser(email string)(User,error){
	var user User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("No user found with email: %s", email)
			return user, nil
		}
		log.Printf("Error fetching user: %v", err)
		return user, err
	}
	return user, nil
}
