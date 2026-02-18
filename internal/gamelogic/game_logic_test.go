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
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 0, col: 3}, ship: createCruiser(), orientation: horizontal},
			expected: nil,
		},
		{input: shipPlacement{
			start: boardMove{row: 3, col: 4}, end: boardMove{row: 6, col: 4}, ship: createBattleship(), orientation: vertical},
			expected: nil,
		},
		{input: shipPlacement{
			start: boardMove{row: 0, col: 3}, end: boardMove{row: 0, col: 1}, ship: createCruiser(), orientation: horizontal},
			expected: ErrInvalidOrientation,
		},
		{input: shipPlacement{
			start: boardMove{row: 6, col: 4}, end: boardMove{row: 3, col: 4}, ship: createCruiser(), orientation: vertical},
			expected: ErrInvalidOrientation,
		},
		{input: shipPlacement{
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 0, col: 4}, ship: createCruiser(), orientation: horizontal},
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
	gs := NewGameState()

	var testShip ship = createCruiser()
	var gameBoard board
	// ship from a1 - a3
	gameBoard.sqaures[0][0] = &testShip
	gameBoard.sqaures[0][1] = &testShip
	gameBoard.sqaures[0][2] = &testShip

	type input struct {
		shipPlacement shipPlacement
		board         board
	}

	cases := []struct {
		input    input
		expected error
	}{
		{input: input{shipPlacement: shipPlacement{
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 0, col: 3}, ship: createCruiser(), orientation: horizontal}, board: gameBoard},
			expected: ErrInvalidOccupiedSqaure,
		},
		{input: input{shipPlacement: shipPlacement{
			start: boardMove{row: 3, col: 4}, end: boardMove{row: 6, col: 4}, ship: createBattleship(), orientation: vertical}, board: gameBoard},
			expected: nil,
		},
		{input: input{shipPlacement: shipPlacement{
			start: boardMove{row: 0, col: 1}, end: boardMove{row: 2, col: 1}, ship: createCruiser(), orientation: vertical}, board: gameBoard},
			expected: ErrInvalidOccupiedSqaure,
		},
		{input: input{shipPlacement: shipPlacement{
			start: boardMove{row: 1, col: 1}, end: boardMove{row: 1, col: 3}, ship: createCarrier(), orientation: horizontal}, board: gameBoard},
			expected: nil,
		},
		{input: input{shipPlacement: shipPlacement{
			start: boardMove{row: 0, col: 9}, end: boardMove{row: 2, col: 9}, ship: createCruiser(), orientation: vertical}, board: gameBoard},
			expected: nil,
		},
	}

	for _, c := range cases {
		err := gs.shipsOccupyRange(c.input.shipPlacement, c.input.board)
		if !(errors.Is(err, c.expected)) {
			t.Errorf("actual: %v; expected: %v", err, c.expected)
		}
	}
}

func TestGetShip(t *testing.T) {
	player := CreatePlayer("tester")
	gs := NewGameState()

	type input struct {
		shipName string
		player   Player
	}

	cases := []struct {
		input    input
		expected error
	}{
		{input: input{shipName: "cruiser", player: player}, expected: nil},
		{input: input{shipName: "battleship", player: player}, expected: nil},
		{input: input{shipName: "destroyer", player: player}, expected: nil},
		{input: input{shipName: "submarine", player: player}, expected: nil},
		{input: input{shipName: "carrier", player: player}, expected: nil},
		{input: input{shipName: "frigate", player: player}, expected: ErrShipNotFound},
	}

	for _, c := range cases {
		_, err := gs.getShip(c.input.shipName, c.input.player)
		if !(errors.Is(err, c.expected)) {
			t.Errorf("GetShip(%s): actual: %v; expected: %v", c.input.shipName, err, c.expected)
		}
	}
}

func TestPickRandomSquare(t *testing.T) {
	// this test is not deterministic but it should be very unlikely to fail unless there is a bug in the function
	for i := 0; i < 100; i++ {
		move := PickRandomSquare()
		if move.row < 0 || move.row >= BOARD_SIZE {
			t.Errorf("PickRandomSquare(): row = %d; expected row between 0 and %d", move.row, BOARD_SIZE-1)
		}
		if move.col < 0 || move.col >= BOARD_SIZE {
			t.Errorf("PickRandomSquare(): col = %d; expected col between 0 and %d", move.col, BOARD_SIZE-1)
		}
	}
}
