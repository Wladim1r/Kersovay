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

func GetInt(someone string) int {
	for {
		ui := bufio.NewReader(os.Stdin)
		str, err := ui.ReadString(NewLine)
		if err != nil {
			fmt.Printf("Ошибка ввода\nПопробуйте еще раз\n")
			continue
		}
		str = strings.TrimSpace(str)

		number, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Ошибка ввода\nПопробуйте еще раз: ")
			continue
		}

		if number < 1 {
			fmt.Printf("Число должно быть положительным\nПопробуйте еще раз: ")
			continue
		}

		if someone == "year" && number > time.Now().Year() {
			fmt.Printf("Указанный вами год не совпадает с текущим\nПопробуйте еще раз: ")
			continue
		}

		return number
	}
}

func GetString() string {
	for {
		ui := bufio.NewReader(os.Stdin)
		str, err := ui.ReadString(NewLine)
		if err != nil {
			fmt.Printf("Ошибка ввода\n")
			continue
		}
		str = strings.TrimSpace(str)

		if len(str) == 0 {
			fmt.Printf("Поле не может оставаться пустым\nПопробуйте еще раз: ")
			continue
		}

		return str
	}
}
