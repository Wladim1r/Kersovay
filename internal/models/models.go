// Package models
package models

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Price  int    `json:"price"`
}

func NewBook(title, author string, year, price int) Book {
	return Book{
		Title:  title,
		Author: author,
		Year:   year,
		Price:  price,
	}
}
