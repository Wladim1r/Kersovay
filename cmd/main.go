package main

import (
	"fmt"
	"library/internal/books/handlers"
	"library/internal/books/repository"
	"library/internal/db"
	"library/internal/models"
	"library/utils"
	"os"
	"os/exec"
	"time"
)

func main() {
	db := db.MustLoad()
	repo := repository.NewBookRepo(db)
	if err := repo.CreateTable(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	handler := handlers.NewHandler(repo)

	fmt.Printf("\nДобро пожаловать в библиотеку. Просьба не шуметь\n")
	fmt.Printf("Выберите один из предложенных вариантов:\n")

	for {
		utils.ShowMenu()
		numberOption := utils.ChooseOption()

		// отчистка экрана
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		switch numberOption {
		case 1:
			if err := handler.ShowAllBooks(); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case 2:
			fmt.Printf("🔍 Укажите название той книги, которая вас интересует\n")
			fmt.Printf("\n")
			title := utils.ChooseTitleBook()
			if err := handler.ShowOneBook(title); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case 3:
			fmt.Printf("\n✨-------------------------------------------✨\n")
			fmt.Printf("✨ Этап добавления новой книги в библиотеку  ✨\n")
			fmt.Printf("✨-------------------------------------------✨\n\n")

			fmt.Printf("➡️ Введите название книги: ")
			title := utils.GetString(false)
			fmt.Printf("✅ Название успешно сохранено!\n\n")

			fmt.Printf("➡️ Введите имя автора книги: ")
			autor := utils.GetString(false)
			fmt.Printf("✅ Автор книги успешно сохранен!\n\n")

			fmt.Printf("➡️ Введите год издания книги: ")
			year := utils.GetInt("year", false)
			fmt.Printf("✅ Дата успешно сохранена!\n\n")

			fmt.Printf("➡️ Введите цену книги (в рублях): ")
			price := utils.GetInt("price", false)
			fmt.Printf("✅ Цена успешно сохранена!\n\n")

			newBook := models.NewBook(title, autor, year, price)
			if err := handler.CreateBook(newBook); err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Print('\n')
		case 4:
			fmt.Printf("\nЭтап удаления книги из списка 🚮\n")
			fmt.Print("🔍 Введите название той книги, которую хотите удалить из списка: ")
			title := utils.GetString(false)
			if err := handler.DeleteBook(title); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case 5:
			fmt.Printf("\nЭтап обновления книги из списка\n")
			fmt.Print(
				"Введите название той книги, которую хотите обновить: ",
			)
			title := utils.GetString(false)

			if err := handler.ShowOneBook(title); err != nil {
				fmt.Println(err.Error())
				continue
			}

			withoutChange := utils.WithoutChange

			fmt.Printf("➡️ Введите новое название книги\n(%s): ", withoutChange)
			newTitle := utils.GetString(true)
			fmt.Printf("✅ Название успешно сохранено!\n\n")

			fmt.Printf("➡️ Введите нового имя автора книги\n(%s): ", withoutChange)
			newAutor := utils.GetString(true)
			fmt.Printf("✅ Автор книги успешно сохранен!\n\n")

			fmt.Printf("➡️ Введите новый год издания книги\n(%s): ", withoutChange)
			newYear := utils.GetInt("year", true)
			fmt.Printf("✅ Дата успешно сохранена!\n\n")

			fmt.Printf("➡️ Введите новую цену книги (в рублях)\n(%s): ", withoutChange)
			newPrice := utils.GetInt("price", true)
			fmt.Printf("✅ Цена успешно сохранена!\n\n")

			updatedBook := models.NewBook(newTitle, newAutor, newYear, newPrice)
			if err := handler.UpdateBook(title, updatedBook); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case 6:
			fmt.Println("Bye-bye 👋")
			time.Sleep(3 * time.Second)
			return
		}
	}
}
