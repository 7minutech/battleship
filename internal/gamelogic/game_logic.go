package gamelogic

import "strings"

func GetWords(input string) []string {
	words := strings.Split(input, " ")
	return words
}
