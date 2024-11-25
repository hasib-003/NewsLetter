package models

import (
	"database/sql"
	"log"

	"github.com/hasib-003/newsLetter/config"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreateUser(email, name, password string) (User, error) {

	query := `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id, email, name, password`
	var user User
	err := config.DB.QueryRow(query, email, name, password).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		log.Println("Error inserting user:", err)
		return user, err
	}
	return user, nil
}
func GetAllUsers() ([]User, error) {
	query := `SELECT id, email, name, password FROM users`
	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Password)
		if err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func GetAUser(email string)(User,error){
	query:=`SELECT id,name,password FROM users WHERE email=$1`
	var user User
	err := config.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Password)
	if err!=nil{
		if err == sql.ErrNoRows {
			log.Printf("No user found with email: %s", email)
			return User{}, nil 
		}
		log.Printf("Error Fetching User:%v",err)
	}
	return user,nil
}
