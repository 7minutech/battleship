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

func repl() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		fmt.Println(scanner.Text())
		words := getWords(scanner.Text())
		for _, w := range words {
			fmt.Println(w)
		}
		runCommand(words)
	}
}
