package utils

import (
	"fmt"
)

func Wait() {
	fmt.Printf("\nДля продолжения нажмите q... ")
	for {
		var input string
		fmt.Scanln(&input)
		if input == "q" {
			break
		}
	}
}
