// Package handlers
package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"library/internal/books/repository"
	"library/internal/models"
	"strings"
)

type bookHandler struct {
	r repository.BookRepository
}

func NewHandler(r repository.BookRepository) *bookHandler {
	return &bookHandler{
		r: r,
	}
}

func (h *bookHandler) CreateBook(book models.Book) error {
	if err := h.r.CreateBook(book); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("❌ такая книга уже существует: %w", err)
		}

		return fmt.Errorf("❌ не удалось создать книгу: %w", err)
	}

	fmt.Printf("🎉 Книга успешно добавлена в общий список!\n\n")
	return nil
}

func (h *bookHandler) ShowAllBooks() error {
	rows, err := h.r.ShowAllBooks()
	if err != nil {
		return fmt.Errorf("❌ ошибка при запросе к бд: %w", err)
	}
	defer rows.Close()

	var k int
	for rows.Next() {
		k++
		book := new(models.Book)

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Price)
		if err != nil {
			return fmt.Errorf("❌ ошибка при чтении записи из бд: %w", err)
		}

		localShow(*book)
	}
	if k == 0 {
		return fmt.Errorf("❌ книг нет")
	}

	return nil
}

func (h *bookHandler) ShowOneBook(title string) error {
	book, err := h.r.ShowOneBook(title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("❌🔍 книга с таким названием не найдена\n\n")
		}
		return fmt.Errorf("❌ ошибка при чтении записи из бд: %w\n\n", err)
	}

	localShow(*book)

	return nil
}

func (h *bookHandler) UpdateBook(title string, book models.Book) error {
	if err := h.r.UpdateBook(title, book); err != nil {
		return fmt.Errorf("❌ не удалось обновить книгу: %w", err)
	}

	fmt.Printf("✅ Книга успешно обновлена!\n\n")
	return nil
}

func (h *bookHandler) DeleteBook(title string) error {
	if err := h.r.DeleteBook(title); err != nil {
		return fmt.Errorf("❌ не удалось удалить книгу: %w", err)
	}

	fmt.Printf("✅ Книга успешно удалена!\n\n")
	return nil
}

func localShow(book models.Book) {
	fmt.Printf("📚------------------📚\n")
	fmt.Printf("|    📖 КНИГА №%d    \n", book.ID)
	fmt.Printf("+-------------------+\n")
	fmt.Printf("| 📝 Название: %s\n", book.Title)
	fmt.Printf("| ✍️ Автор: %s\n", book.Author)
	fmt.Printf("| 🗓️ Год издания: %d\n", book.Year)
	fmt.Printf("| 💰 Цена (в рублях): %d\n\n\n", book.Price)
	fmt.Printf("📚-------------------+\n")
}
