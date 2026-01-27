package gamelogic

import (
	"errors"
	"testing"
)

func TestGetWords(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{"move", []string{"move"}},
		{"move a1", []string{"move", "a1"}},
		{"place b2", []string{"place", "b2"}},
		{"a b c d", []string{"a", "b", "c", "d"}},
		{"", []string{""}},
	}

	for _, c := range cases {
		actual := GetWords(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("len(%s) = %d; not %d)", c.input, len(actual), len(c.expected))
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("%s[%d] = %s; not %s", c.input, i, actual[i], c.expected[i])
			}
		}
	}

}

func TestConvertMove(t *testing.T) {
	type expected struct {
		move boardMove
	}
	cases := []struct {
		input    string
		expected boardMove
	}{
		{input: "a1", expected: boardMove{row: 0, col: 0}},
		{input: "!2", expected: boardMove{}},
		{input: "b10", expected: boardMove{row: 1, col: 9}},
		{input: "word", expected: boardMove{}},
		{input: "c4", expected: boardMove{row: 2, col: 3}},
		{input: "j8", expected: boardMove{row: 9, col: 7}},
	}

	for _, c := range cases {
		mv, _ := convertMove(c.input)
		if mv.row != c.expected.row {
			t.Errorf("ConvertMove(%s): boardMove.row = %d; expected boardMove.row = %d", c.input, mv.row, c.expected.row)
		}
		if mv.col != c.expected.col {
			t.Errorf("ConvertMove(%s): boardMove.col = %d; expected boardMove.col = %d", c.input, mv.col, c.expected.col)
		}
	}
}

func TestValidShipRange(t *testing.T) {
	cases := []struct {
		input    shipPlacement
		expected error
	}{
		{input: shipPlacement{
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 0, col: 3}, ship: startCruiser, orientation: horizantal},
			expected: nil,
		},
		{input: shipPlacement{
			start: boardMove{row: 3, col: 4}, end: boardMove{row: 6, col: 4}, ship: startBattleship, orientation: vertical},
			expected: nil,
		},
		{input: shipPlacement{
			start: boardMove{row: 0, col: 3}, end: boardMove{row: 0, col: 1}, ship: startCruiser, orientation: horizantal},
			expected: ErrInvalidOrientation,
		},
		{input: shipPlacement{
			start: boardMove{row: 6, col: 4}, end: boardMove{row: 3, col: 4}, ship: startCruiser, orientation: vertical},
			expected: ErrInvalidOrientation,
		},
		{input: shipPlacement{
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 0, col: 3}, ship: startCarrier, orientation: horizantal},
			expected: ErrInvalidShipSpace,
		},
	}

	for _, c := range cases {
		err := validShipRange(c.input)
		if !(errors.Is(err, c.expected)) {
			t.Errorf("actual: %v; expected: %v", err, c.expected)
		}
	}
}

func TestShipsOccupyRange(t *testing.T) {
	gs := NewGameState(Player{})

	var testShip ship = startCruiser
	// ship from a1 - a3
	gs.gameBoard.sqaures[0][0] = &testShip
	gs.gameBoard.sqaures[0][1] = &testShip
	gs.gameBoard.sqaures[0][2] = &testShip

	cases := []struct {
		input    shipPlacement
		expected error
	}{
		{input: shipPlacement{
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 0, col: 3}, ship: startCruiser, orientation: horizantal},
			expected: ErrInvalidOccupiedSqaure,
		},
		{input: shipPlacement{
			start: boardMove{row: 3, col: 4}, end: boardMove{row: 6, col: 4}, ship: startBattleship, orientation: vertical},
			expected: nil,
		},
		{input: shipPlacement{
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 2, col: 1}, ship: startCruiser, orientation: vertical},
			expected: ErrInvalidOccupiedSqaure,
		},
		{input: shipPlacement{
			start: boardMove{row: 1, col: 1}, end: boardMove{row: 1, col: 3}, ship: startCarrier, orientation: horizantal},
			expected: nil,
		},
		{input: shipPlacement{
			start: boardMove{row: 0, col: 9}, end: boardMove{row: 2, col: 9}, ship: startCruiser, orientation: vertical},
			expected: nil,
		},
	}

	for _, c := range cases {
		err := gs.shipsOccupyRange(c.input)
		if !(errors.Is(err, c.expected)) {
			t.Errorf("actual: %v; expected: %v", err, c.expected)
		}
	}
}
