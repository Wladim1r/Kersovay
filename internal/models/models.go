// Package models
package models

type User struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Books    []Book `json:"books"`
}

type Book struct {
	ID     uint   `json:"-"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Price  int    `json:"price"`
	UserID uint   `json:"-"`
	User   User   `json:"-"`
}

func NewBook(title, author string, year, price int) Book {
	return Book{
		Title:  title,
		Author: author,
		Year:   year,
		Price:  price,
	}
}
