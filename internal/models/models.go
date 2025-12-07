// Package models
package models

type Book struct {
	ID     int
	Title  string
	Author string
	Year   int
	Price  int
}

func NewBook(title, author string, year, price int) Book {
	return Book{
		Title:  title,
		Author: author,
		Year:   year,
		Price:  price,
	}
}
