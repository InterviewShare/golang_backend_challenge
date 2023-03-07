package models

type Article struct {
	ID		string		`json:"id"`
	Title	string	 	`json:"title"`
	Content	string	 	`json:"content"`
	Author	string	 	`json:"author"`
}