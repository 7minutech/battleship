package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getWords(input string) []string {
	words := strings.Split(input, " ")
	return words
}

func repl(handler func([]string)) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := getWords(scanner.Text())
		handler(words)
	}
}
