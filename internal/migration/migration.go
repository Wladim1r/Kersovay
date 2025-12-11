package migration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"library/internal/auth/repository"
	booksRepo "library/internal/books/repository"
	"library/internal/models"
)

type MigrationData struct {
	Users []models.User `json:"users"`
}

type Migration struct {
	authRepo repository.AuthRepository
	bookRepo booksRepo.BookRepository
}

func NewMigration(
	authRepo repository.AuthRepository,
	bookRepo booksRepo.BookRepository,
) *Migration {
	return &Migration{authRepo: authRepo, bookRepo: bookRepo}
}

func (m *Migration) ExportData() error {
	users, err := m.authRepo.GetAllUsers()
	if err != nil {
		return err
	}

	for i := range users {
		rows, err := m.bookRepo.ShowAllBooks(users[i].ID)
		if err != nil {
			return err
		}

		var books []models.Book
		for rows.Next() {
			var book models.Book
			if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Price, &book.UserID); err != nil {
				return err
			}
			books = append(books, book)
		}
		users[i].Books = books
	}

	data := MigrationData{Users: users}

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("migration.json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (m *Migration) ImportData() error {
	file, err := ioutil.ReadFile("migration.json")
	if err != nil {
		return fmt.Errorf("ошибка чтения файла миграции: %w", err)
	}

	var rawData map[string]json.RawMessage
	err = json.Unmarshal(file, &rawData)
	if err != nil {
		return fmt.Errorf("ошибка десериализации сырых данных: %w", err)
	}

	var usersRaw []json.RawMessage
	err = json.Unmarshal(rawData["users"], &usersRaw)
	if err != nil {
		return fmt.Errorf("ошибка десериализации сырых данных пользователей: %w", err)
	}

	if err := m.authRepo.TruncateUsers(); err != nil {
		return fmt.Errorf("ошибка очистки таблицы пользователей: %w", err)
	}
	if err := m.bookRepo.TruncateBooks(); err != nil {
		return fmt.Errorf("ошибка очистки таблицы книг: %w", err)
	}

	// Вспомогательная структура для раздельной десериализации пользователя и его книг
	type UserWithBooksRaw struct {
		models.User
		Books []json.RawMessage `json:"books"`
	}

	for i, userRaw := range usersRaw {
		fmt.Printf("Обработка пользователя %d\n", i+1)
		var userWithBooks UserWithBooksRaw
		if err := json.Unmarshal(userRaw, &userWithBooks); err != nil {
			fmt.Printf("  - ошибка десериализации пользователя: %v\n", err)
			continue
		}

		user := userWithBooks.User
		if user.Username == "" || user.Password == "" {
			fmt.Printf("  - неверные данные пользователя: %+v\n", user)
			continue
		}

		fmt.Printf("  - Создание пользователя: %s\n", user.Username)
		if err := m.authRepo.CreateUser(user.Username, user.Password); err != nil {
			fmt.Printf("  - ошибка создания пользователя %s: %v\n", user.Username, err)
			continue
		}
		fmt.Printf("  - Пользователь %s успешно создан\n", user.Username)

		id, _, err := m.authRepo.GetUserByUsername(user.Username)
		if err != nil {
			fmt.Printf("  - ошибка получения пользователя %s: %v\n", user.Username, err)
			continue
		}

		if len(userWithBooks.Books) == 0 {
			fmt.Printf("  - Пользователь %s не имеет книг для импорта\n", user.Username)
			continue
		}

		fmt.Printf("  - Импорт книг для пользователя %s\n", user.Username)
		for _, bookRaw := range userWithBooks.Books {
			var book models.Book
			if err := json.Unmarshal(bookRaw, &book); err != nil {
				fmt.Printf("    - ошибка десериализации книги: %v\n", err)
				continue
			}

			fmt.Printf("    - Создание книги: %s\n", book.Title)
			if err := m.bookRepo.CreateBook(book, id); err != nil {
				fmt.Printf(
					"    - ошибка создания книги '%s' для пользователя %s: %v\n",
					book.Title,
					user.Username,
					err,
				)
				continue
			}
			fmt.Printf("    - Книга %s успешно создана\n", book.Title)
		}
		fmt.Printf("  - Завершено импорт книг для пользователя %s\n", user.Username)
	}

	fmt.Println("Импорт данных завершен.")
	return nil
}
