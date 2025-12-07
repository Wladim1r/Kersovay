// Package repository
package repository

import (
	"database/sql"
	"fmt"
	"library/internal/models"
	"strings"
)

type BookRepository interface {
	CreateTable() error
	CreateBook(book models.Book) error
	ShowAllBooks() (*sql.Rows, error)
	ShowOneBook(title string) (*models.Book, error)
	UpdateBook(title string, book models.Book) error
	DeleteBook(title string) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

func (br *bookRepository) CreateTable() error {
	_, err := br.db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			year INTEGER,
			price INTEGER,
			UNIQUE(title, author)
		);`)

	if err != nil {
		return fmt.Errorf("❌ не удалось создать базу данных: %w", err)
	}

	return nil
}

func (br *bookRepository) CreateBook(book models.Book) error {
	_, err := br.db.Exec(
		"INSERT INTO books (title, author, year, price) VALUES (?, ?, ?, ?)",
		book.Title,
		book.Author,
		book.Year,
		book.Price,
	)

	return err
}

func (br *bookRepository) ShowAllBooks() (*sql.Rows, error) {
	query := "SELECT * FROM books ORDER BY title ASC"

	rows, err := br.db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (br *bookRepository) ShowOneBook(title string) (*models.Book, error) {
	book := new(models.Book)

	query := "%" + title + "%"

	row := br.db.QueryRow("SELECT * FROM books WHERE title LIKE ?", query)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Price)

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (br *bookRepository) UpdateBook(title string, book models.Book) error {
	var updates []string
	var args []interface{}

	if book.Title != "" {
		updates = append(updates, "title = ?")
		args = append(args, book.Title)
	}
	if book.Author != "" {
		updates = append(updates, "author = ?")
		args = append(args, book.Author)
	}
	if book.Year != 0 {
		updates = append(updates, "year = ?")
		args = append(args, book.Year)
	}
	if book.Price != 0 {
		updates = append(updates, "price = ?")
		args = append(args, book.Price)
	}

	if len(updates) == 0 {
		return nil // нет полей для обновления
	}

	query := fmt.Sprintf("UPDATE books SET %s WHERE title = ?", strings.Join(updates, ", "))
	args = append(args, title)

	_, err := br.db.Exec(query, args...)
	return err
}

func (br *bookRepository) DeleteBook(title string) error {
	_, err := br.db.Exec("DELETE FROM books WHERE title = ?", title)

	return err
}
