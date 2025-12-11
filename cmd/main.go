package main

import (
	"fmt"
	authHandlers "library/internal/auth/handlers"
	authRepo "library/internal/auth/repository"
	"library/internal/books/handlers"
	"library/internal/books/repository"
	"library/internal/db"
	"library/internal/migration"
	"library/internal/models"
	"library/utils"
	"os"
	"os/exec"
)

func main() {
	database := db.MustLoad()
	defer database.Close()

	authRepository := authRepo.NewAuthRepo(database)
	bookRepo := repository.NewBookRepo(database)

	if err := authRepository.CreateUserTable(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := bookRepo.CreateTable(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	authHandler := authHandlers.NewAuthHandler(authRepository)
	bookHandler := handlers.NewHandler(bookRepo)

	var currentUserID int
	var currentUsername string

authMenu:
	for {
		clearScreen()
		utils.ShowMainMenu()
		choice := utils.GetMenuItem(5)

		clearScreen()

		switch choice {
		case 1:
			fmt.Printf("\n🔑 ВХОД В СИСТЕМУ\n\n")
			fmt.Print("Имя пользователя: ")
			username := utils.GetString(false)
			fmt.Print("Пароль: ")
			password := utils.GetString(false)

			userID, err := authHandler.Login(username, password)
			if err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}

			currentUserID = userID
			currentUsername = username
			goto mainMenu

		case 2:
			fmt.Printf("\n📝 РЕГИСТРАЦИЯ\n\n")
			fmt.Print("Имя пользователя (минимум 3 символа): ")
			username := utils.GetString(false)
			fmt.Print("Пароль (минимум 5 символов): ")
			password := utils.GetString(false)

			if err := authHandler.Register(username, password); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}

			utils.Wait()
			continue

		case 3:
			fmt.Println("Всего доброго! 👋")
			return
		case 4:
			if err := migration.NewMigration(authRepository, bookRepo).ExportData(); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}
			fmt.Println("Данные успешно экспортированы в migration.json")
			utils.Wait()
		case 5:
			if err := migration.NewMigration(authRepository, bookRepo).ImportData(); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}
			fmt.Println("Данные успешно импортированы из migration.json")
			utils.Wait()

		default:
			fmt.Println("❌ Неверный выбор")
			utils.Wait()
			continue
		}
	}

mainMenu:
	fmt.Printf("\n🎉 Добро пожаловать, %s! Просьба не шуметь\n", currentUsername)

	for {
		clearScreen()

		utils.ShowBookMenu()
		numberOption := utils.GetMenuItem(6)

		clearScreen()

		switch numberOption {
		case 1:
			if err := bookHandler.ShowAllBooks(currentUserID); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}
			utils.Wait()
		case 2:
			fmt.Printf("🔍 Укажите название той книги, которая вас интересует\n")
			fmt.Printf("\n")
			title := utils.ChooseTitleBook()
			if err := bookHandler.ShowOneBook(title, currentUserID); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}
			utils.Wait()
		case 3:
			fmt.Printf("\n✨-------------------------------------------✨\n")
			fmt.Printf("✨ Этап добавления новой книги в библиотеку  ✨\n")
			fmt.Printf("✨-------------------------------------------✨\n\n")

			fmt.Printf("➡️ Введите название книги: ")
			title := utils.GetString(false)

			fmt.Printf("➡️ Введите имя автора книги: ")
			autor := utils.GetString(false)

			fmt.Printf("➡️ Введите год издания книги: ")
			year := utils.GetInt("year", false)

			fmt.Printf("➡️ Введите цену книги (в рублях): ")
			price := utils.GetInt("price", false)

			newBook := models.NewBook(title, autor, year, price)
			if err := bookHandler.CreateBook(newBook, currentUserID); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}

			fmt.Print('\n')
		case 4:
			fmt.Printf("\nЭтап удаления книги из списка 🚮\n")
			fmt.Print("🔍 Введите название той книги, которую хотите удалить из списка: ")
			title := utils.GetString(false)

			if err := bookHandler.ShowOneBook(title, currentUserID); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}
			if err := bookHandler.DeleteBook(title, currentUserID); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}
		case 5:
			fmt.Printf("\nЭтап обновления книги из списка\n")
			fmt.Print("Введите название той книги, которую хотите обновить: ")
			title := utils.GetString(false)

			if err := bookHandler.ShowOneBook(title, currentUserID); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}

			withoutChange := utils.WithoutChange

			fmt.Printf("➡️ Введите новое название книги\n(%s): ", withoutChange)
			newTitle := utils.GetString(true)

			fmt.Printf("➡️ Введите нового имя автора книги\n(%s): ", withoutChange)
			newAutor := utils.GetString(true)

			fmt.Printf("➡️ Введите новый год издания книги\n(%s): ", withoutChange)
			newYear := utils.GetInt("year", true)

			fmt.Printf("➡️ Введите новую цену книги (в рублях)\n(%s): ", withoutChange)
			newPrice := utils.GetInt("price", true)

			updatedBook := models.NewBook(newTitle, newAutor, newYear, newPrice)
			if err := bookHandler.UpdateBook(title, updatedBook, currentUserID); err != nil {
				fmt.Println(err.Error())
				utils.Wait()
				continue
			}
		case 6:
			fmt.Println("До свидания, " + currentUsername + "! 👋")
			utils.Wait()
			goto authMenu
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
