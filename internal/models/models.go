// Package models
package models

type User struct {
	ID       int    `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string
	Books    []Book `gorm:"foreignKey:UserID"`
}

type Book struct {
	ID     uint `gorm:"primaryKey"`
	Title  string
	Author string
	Year   int
	Price  int

	UserID uint
	User   User `gorm:"constraint:OnDelete:CASCADE;"`
}

func NewBook(title, author string, year, price int) Book {
	return Book{
		Title:  title,
		Author: author,
		Year:   year,
		Price:  price,
	}
}
