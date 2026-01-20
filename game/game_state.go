package game

const BOARD_SIZE = 10

const STARTING_SHIP_COUNT = 5

type PlayerTurn int

const (
	player1 PlayerTurn = 0
	player2 PlayerTurn = 1
)

type gameState struct {
	gameBoard board
	player1   player
	player2   player
	turn      PlayerTurn
}

type board struct {
	sqaures [BOARD_SIZE][BOARD_SIZE]*ship
}

type player struct {
	userName  string
	shipCount int
}

type ship struct {
	length   int
	start    string
	end      string
	occupies map[string]bool // a1[false], a2[false], a3[false], a4[true]
	hp       int
	alive    bool
}
