package main

import (
	"fmt"
	"strings"
)

func getWords(input string) []string {
	words := strings.Split(input, " ")
	return words
}

func repl() {
	for {
		var input string
		fmt.Print(">>> ")
		fmt.Scanln(&input)
		runCommand(input)
	}
}
