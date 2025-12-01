package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const NewLine = '\n'

func ShowMenu() {
	fmt.Printf("+--------------------------------------------------+\n")
	fmt.Printf("|             -+-       МЕНЮ       -+-             |\n")
	fmt.Printf("+--------------------------------------------------+\n")
	fmt.Printf("| 📚 1 - ознакомиться со всеми книгами в библотеке |\n")
	fmt.Printf("| 📘 2 - ознакомиться с определенной книгой        |\n")
	fmt.Printf("| ➕ 3 - добавить новую книгу                      |\n")
	fmt.Printf("| 🗑️ 4 - удалить книгу                             |\n")
	fmt.Printf("| 🔄 5 - обновить книгу                            |\n")
	fmt.Printf("| 🚪 6 - уйти                                      |\n")
	fmt.Printf("+--------------------------------------------------+\n\n")
}

func ChooseOption() int {
	var numberOption int

	for {
		fmt.Print("Поле для ввода действия над библиотекой: ")

		ui := bufio.NewReader(os.Stdin)
		str, err := ui.ReadString(NewLine)
		if err != nil {
			fmt.Printf("Ошибка ввода\nПопробуйте еще раз\n\n")
			continue
		}

		str = strings.TrimSpace(str)

		if len(str) == 0 {
			fmt.Printf("Нельзя оставлять поле пустым\n\n")
			continue
		}

		numberOption, err = strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Ввод не может содержать какие-либо символы, кроме цифр\n\n")
			continue
		}

		if numberOption > 6 {
			fmt.Printf("Число слишком большое\n\n")
			continue
		}
		if numberOption < 1 {
			fmt.Printf("Число слишком маленькое\n\n")
			continue
		}
		break
	}
	return numberOption
}

func ChooseTitleBook() string {
	for {
		fmt.Print("Поле для ввода названия книги: ")

		ui := bufio.NewReader(os.Stdin)
		str, err := ui.ReadString(NewLine)
		if err != nil {
			fmt.Printf("Ошибка ввода\nПопробуйте еще раз\n\n")
			continue
		}
		str = strings.TrimSpace(str)

		if len(str) == 0 {
			fmt.Printf("Ввод не должен быть пустым\n\n")
			continue
		}
		return str
	}
}
