package gamelogic

const BOARD_SIZE = 10

const STARTING_SHIP_COUNT = 5

var startCruiser = ship{
	name:    "cruiser",
	length:  3,
	start:   boardMove{},
	end:     boardMove{},
	modules: map[boardMove]bool{},
	hp:      3,
	alive:   true,
}

type PlayerTurn int

const (
	player1 PlayerTurn = 0
	player2 PlayerTurn = 1
)

type orientation int

const (
	horizantal orientation = 0
	vertical   orientation = 1
)

type gameState struct {
	player    Player
	turn      PlayerTurn
	gameBoard board
}

type board struct {
	sqaures [BOARD_SIZE][BOARD_SIZE]*ship
}

type Player struct {
	userName  string
	shipCount int
	ships     []ship
}

type ship struct {
	name    string
	length  int
	start   boardMove
	end     boardMove
	modules map[boardMove]bool // a1[false], a2[false], a3[false], a4[true]
	hp      int
	alive   bool
}

type boardMove struct {
	row int
	col int
}

type shipPlacement struct {
	start       boardMove
	end         boardMove
	ship        ship
	orientation orientation
}
