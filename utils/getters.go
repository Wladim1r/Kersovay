// Package utils
package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	WithoutChange = "оставьте поле пустым, если не хотите ничего менять"
	MinimalYear   = 800
)

func GetInt(someone string, empty bool) int {
	for {
		ui := bufio.NewReader(os.Stdin)
		str, err := ui.ReadString(NewLine)
		if err != nil {
			fmt.Printf("Ошибка ввода\nПопробуйте еще раз\n")
			continue
		}
		str = strings.TrimSpace(str)

		if empty && str == "" {
			return 0
		}

		number, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Ошибка ввода\nПопробуйте еще раз: ")
			continue
		}

		if number < 1 {
			fmt.Printf("Число должно быть положительным\nПопробуйте еще раз: ")
			continue
		}

		if someone == "year" {
			if number < MinimalYear {
				fmt.Printf("Слишком старая книга, наша библиотека такие не принимает 🤷")
				continue
			}
			if number > time.Now().Year() {
				fmt.Printf("Указанный вами год превышает текущий\nПопробуйте еще раз: ")
				continue
			}
		}

		return number
	}
}

func GetString(empty bool) string {
	for {
		ui := bufio.NewReader(os.Stdin)
		str, err := ui.ReadString(NewLine)
		if err != nil {
			fmt.Printf("Ошибка ввода\n")
			continue
		}
		str = strings.TrimSpace(str)

		if len(str) == 0 {
			if empty {
				return ""
			}

			fmt.Printf("Поле не может оставаться пустым\nПопробуйте еще раз: ")
			continue
		}

		return str
	}
}
