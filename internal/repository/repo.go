// Package repository
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"library/internal/models"
)

type userDB struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *userDB {
	return &userDB{
		db: db,
	}
}

func (db *userDB) CreateTable() error {
	_, err := db.db.Exec(
		"CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, title TEXT, author TEXT, year INTEGER, price INTEGER)",
	)

	if err != nil {
		return fmt.Errorf("не удалось создать базу данных: %w", err)
	}

	return nil
}

func (db *userDB) CreateBook(book models.Book) error {
	_, err := db.db.Exec(
		"INSERT INTO books (title, author, year, price) VALUES (?, ?, ?, ?)",
		book.Title,
		book.Author,
		book.Year,
		book.Price,
	)

	if err != nil {
		return fmt.Errorf("не удалось создать книгу: %w", err)
	}

	fmt.Printf("Книга успешно добавлена в общий список!\n\n")

	return nil
}

func (db *userDB) DeleteBook(title string) error {
	_, err := db.db.Exec("DELETE FROM books WHERE title = ?", title)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("книга с таким названием не найдена: %w", err)
		default:
			return fmt.Errorf("не удалось удалить книгу: %w", err)
		}
	}

	return nil
}

func (db *userDB) ShowAllBooks() error {
	query := "SELECT * FROM books"

	rows, err := db.db.Query(query)
	if err != nil {
		return fmt.Errorf("ошибка при запросе к бд: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		book := new(models.Book)

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Price)
		if err != nil {
			return fmt.Errorf("ошибка при чтении записи из бд: %w", err)
		}

		localShow(*book)
	}

	return nil
}

func (db *userDB) ShowOneBook(title string) error {
	book := new(models.Book)

	row := db.db.QueryRow("SELECT * FROM books WHERE title = ?", title)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Price)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("книга с таким названием не найдена: %w", err)
		default:
			return fmt.Errorf("ошибка при чтении записи из бд: %w", err)
		}
	}

	localShow(*book)

	return nil
}

func localShow(book models.Book) {
	fmt.Println("📚------------------------------------📚")
	fmt.Printf("📖 КНИГА №%d\n", book.ID)
	fmt.Printf("📝 Название: %s\n", book.Title)
	fmt.Printf(" ✍️ Автор: %s\n", book.Author)
	fmt.Printf("🗓️ Год издания: %d\n", book.Year)
	fmt.Printf("💰 Цена (в рублях): %d\n", book.Price)
	fmt.Printf("📚------------------------------------📚\n\n")
}
