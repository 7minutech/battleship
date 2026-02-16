package gamelogic

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ClientHelp() {
	fmt.Println("help: prints possible commands")
	fmt.Println("place: place a ship on the board; usage: place <ship type> <x> <y> <orientation>")
	fmt.Println("show: shows the current state of the board")
	fmt.Println("ships: cruiser, destroyer, submarine, battleship, carrier")
	fmt.Println("quit: exits the program")
	// make this a map of command to description and loop through it instead of hardcoding the print statements
}

func ServerHelp() {
	fmt.Println("help: prints possible commands")
	fmt.Println("quit: exits the program")
	fmt.Println("pause: sends a message to all clients that the server is paused")
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
