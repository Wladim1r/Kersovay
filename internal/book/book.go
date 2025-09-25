// Package book
package book

import (
	"encoding/json"
	"fmt"
	"os"
)

type Book struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Price  int    `json:"price"`
}

type BookStore struct {
	Books []Book `json:"books"`
	MaxID int    `json:"max_id"`
}

func NewBook(name, author string, year, price int) Book {
	return Book{
		Name:   name,
		Author: author,
		Year:   year,
		Price:  price,
	}
}

func NewBookStore() BookStore {
	return BookStore{}
}

func (books *BookStore) CreateBook(fileName string, book Book) {
	ReadJSONFile(fileName, books)

	books.MaxID = 0

	if len(books.Books) == 0 {
		books.MaxID = 1
	}

	for _, book := range books.Books {
		if book.ID >= books.MaxID {
			books.MaxID = book.ID + 1
		}
	}
	book.ID = books.MaxID

	books.Books = append(books.Books, book)

	dataJSON, err := json.Marshal(*books)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при записи в JSON формат")
		return
	}
	err = os.WriteFile(fileName, dataJSON, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при записи данных в файл")
	}

	fmt.Fprintf(
		os.Stdout,
		"Книга успешно сохранена под номером %d и добавлена в общий список!\n",
		book.ID,
	)
}

func (books *BookStore) RemoveBook(fileName string, title string) {
	ReadJSONFile(fileName, books)

	for i, val := range books.Books {
		if val.Name == title {
			books.Books = append(books.Books[:i], books.Books[i+1:]...)

			for i := 0; i < len(books.Books); i++ {
				books.Books[i].ID = i + 1
			}

			dataJSON, err := json.Marshal(*books)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при записи в JSON формат")
			}
			err = os.WriteFile(fileName, dataJSON, 0666)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при записи данных в файл")
			}

			fmt.Fprintf(os.Stdout, "Книга \"%s\" успешно удалена из списка!\n", val.Name)
			fmt.Fprintln(os.Stdout)

			return
		}
	}
	fmt.Fprintln(os.Stdout, "Книга по данному названию не найдена")
	fmt.Fprintln(os.Stdout)
}

func CreatFile(filePath string) (*os.File, error) {
	var file *os.File

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка при создании файла")
			return nil, err
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_RDWR, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка при открытии файла")
			return nil, err
		}
	}
	defer file.Close()

	return file, nil
}

func ReadJSONFile(fileName string, bookStore *BookStore) {
	dataJSON, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при чтении файла")
		return
	}

	if len(dataJSON) == 0 {
		return
	}

	err = json.Unmarshal(dataJSON, bookStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при чтении файла: %v\n", err)
		return
	}
}
