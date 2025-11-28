package main

import (
	"fmt"
	"library/internal/database"
	"library/internal/models"
	"library/internal/repository"
	"library/utils"
)

func main() {

	db := database.MustLoad()
	repo := repository.NewUserDB(db)
	if err := repo.CreateTable(); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("\nДобро пожаловать в библиотеку. Просьба не шуметь\n")
	fmt.Printf("Выберите один из предложенных вариантов:\n")

	for {
		utils.ShowMenu()
		numberOption := utils.ChooseOption()

		switch numberOption {
		case 1:
			if err := repo.ShowAllBooks(); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case 2:
			fmt.Printf("Укажите название той книги, которая вас интересует\n")
			fmt.Printf("\n")
			title := utils.ChooseTitleBook()
			if err := repo.ShowOneBook(title); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case 3:
			fmt.Printf("\n✨------------------------------------✨\n")
			fmt.Printf("✨ Этап добавления новой книги в библиотеку✨\n")
			fmt.Printf("✨------------------------------------✨\n\n")

			fmt.Printf("Введите название книги: ")
			title := utils.GetString()
			fmt.Printf("✅ Название успешно сохранено!\n\n")

			fmt.Printf("Введите имя автора книги: ")
			autor := utils.GetString()
			fmt.Printf("✅ Автор книги успешно сохранен!\n\n")

			fmt.Printf("Введите год издания книги: ")
			year := utils.GetInt("year")
			fmt.Printf("✅ Дата успешно сохранена!\n\n")

			fmt.Printf("Введите цену книги (в рублях): ")
			price := utils.GetInt("price")
			fmt.Printf("✅ Цена успешно сохранена!\n\n")

			newBook := models.NewBook(title, autor, year, price)
			if err := repo.CreateBook(newBook); err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Print('\n')
		case 4:
			fmt.Printf("\nЭтап удаления книги из списка\n")
			fmt.Print("Введите название той книги, которую хотите удалить из списка: ")
			title := utils.GetString()
			if err := repo.DeleteBook(title); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case 5:
			return
		}
	}
}
