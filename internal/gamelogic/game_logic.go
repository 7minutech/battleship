package gamelogic

import (
	"fmt"
	"strconv"
	"strings"
)

func GetWords(input string) []string {
	words := strings.Split(input, " ")
	return words
}

type boardMove struct {
	row int
	col int
}

func convertMove(mv string) (boardMove, error) {
	if len(mv) != 2 {
		return boardMove{}, fmt.Errorf("error: move was not exactly 2: %s", mv)
	}

	letter := mv[0]
	num := string(mv[1])

	numInt, err := strconv.Atoi(num)
	if err != nil {
		return boardMove{}, fmt.Errorf("error: could not convert number in move to int: %s", mv)
	}

	if !(int('a') <= int(letter) && int(letter) <= int('z')) {
		return boardMove{}, fmt.Errorf("error: move did not contain a letter that was between a-z: %s", mv)
	}

	if 1 <= numInt && numInt <= 9 {
		return boardMove{}, fmt.Errorf("error: move did not contain a number that was between 1-9: %s", mv)
	}

	row := (int('a') - int(letter))
	col := numInt - 1

	return boardMove{row: row, col: col}, nil
}
