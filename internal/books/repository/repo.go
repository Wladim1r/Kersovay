package repository

import (
	"database/sql"
	"fmt"
	"library/internal/models"
	"strings"
)

type BookRepository interface {
	CreateTable() error
	CreateBook(book models.Book, userID int) error
	ShowAllBooks(userID int) (*sql.Rows, error)
	ShowOneBook(title string, userID int) (*models.Book, error)
	UpdateBook(title string, book models.Book, userID int) error
	DeleteBook(title string, userID int) error
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		year INTEGER,
		price INTEGER,
		user_id INTEGER NOT NULL,
		UNIQUE(title, author, user_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	)

	if err != nil {
		return fmt.Errorf("❌ не удалось создать базу данных: %w", err)
	}

	return nil
}

func (br *bookRepository) CreateBook(book models.Book, userID int) error {
	_, err := br.db.Exec(
		"INSERT INTO books (title, author, year, price, user_id) VALUES (?, ?, ?, ?, ?)",
		book.Title,
		book.Author,
		book.Year,
		book.Price,
		userID,
	)

	return err
}

func (br *bookRepository) ShowAllBooks(userID int) (*sql.Rows, error) {
	query := "SELECT id, title, author, year, price, user_id FROM books WHERE user_id = ? ORDER BY title ASC"

	rows, err := br.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (br *bookRepository) ShowOneBook(title string, userID int) (*models.Book, error) {
	book := new(models.Book)

	query := "%" + title + "%"

	row := br.db.QueryRow(
		"SELECT id, title, author, year, price, user_id FROM books WHERE title LIKE ? AND user_id = ?",
		query,
		userID,
	)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Price, &book.UserID)

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (br *bookRepository) UpdateBook(title string, book models.Book, userID int) error {
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
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE books SET %s WHERE title = ? AND user_id = ?",
		strings.Join(updates, ", "),
	)
	args = append(args, title, userID)

	_, err := br.db.Exec(query, args...)
	return err
}

func (br *bookRepository) DeleteBook(title string, userID int) error {
	_, err := br.db.Exec("DELETE FROM books WHERE title = ? AND user_id = ?", title, userID)

	return err
}
