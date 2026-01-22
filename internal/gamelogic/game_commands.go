package gamelogic

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func Help() {
	fmt.Println("help: prints possible commands")
	fmt.Println("quit: exits the program")
}

func Quit() {
	fmt.Println("exiting program...")
	os.Exit(0)
}

func Welcome() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Wecome to battle ship please enter your username!")
	fmt.Print("username: ")
	scanner.Scan()
	words := GetWords(scanner.Text())
	if len(words) == 0 {
		return "", errors.New("error: did not provide input for user name")
	}
	return words[0], nil
}
