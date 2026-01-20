package main

import (
	"fmt"
)

func repl() {
	for {
		var input string
		fmt.Print(">>> ")
		fmt.Scanln(&input)
		runCommand(input)
	}
}
