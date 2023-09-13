package models

type User struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	LastName	string  `json:"lastName"`
	Email		string	`json:"email"`
	Password	string 	`json:"password"`
	IsAdmin     bool    `json:"isAdmin"`
}