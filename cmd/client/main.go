package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/7minutech/battleship/internal/gamelogic"
)

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := gamelogic.GetWords(scanner.Text())
		fmt.Print(words)
	}
}
