package gamelogic

const BOARD_SIZE = 10

const STARTING_SHIP_COUNT = 5

const CRUISER_LENGTH, SUBMARINE_LENGTH = 3, 3
const BATTLESHIP_LENGTH = 4
const CARRIER_LENGTH = 5
const DESTROYER_LENGTH = 2

var startCruiser = ship{
	name:    "cruiser",
	length:  CRUISER_LENGTH,
	start:   boardMove{},
	end:     boardMove{},
	modules: map[boardMove]bool{},
	hp:      CRUISER_LENGTH,
	alive:   true,
}

var startBattleship = ship{
	name:    "battleship",
	length:  BATTLESHIP_LENGTH,
	start:   boardMove{},
	end:     boardMove{},
	modules: map[boardMove]bool{},
	hp:      BATTLESHIP_LENGTH,
	alive:   true,
}

var startCarrier = ship{
	name:    "carrier",
	length:  CARRIER_LENGTH,
	start:   boardMove{},
	end:     boardMove{},
	modules: map[boardMove]bool{},
	hp:      CARRIER_LENGTH,
	alive:   true,
}

var startSubmarine = ship{
	name:    "submarine",
	length:  SUBMARINE_LENGTH,
	start:   boardMove{},
	end:     boardMove{},
	modules: map[boardMove]bool{},
	hp:      SUBMARINE_LENGTH,
	alive:   true,
}

var startDestroyer = ship{
	name:    "destroyer",
	length:  DESTROYER_LENGTH,
	start:   boardMove{},
	end:     boardMove{},
	modules: map[boardMove]bool{},
	hp:      DESTROYER_LENGTH,
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
