package routing

const EXCHANGE_BATTLESHIP_DIRECT = "battleship_direct"
const PAUSE_KEY = "pause_started"
const NEW_PLAYER_KEY = "player.created"

type PauseMessage struct {
	Content string
}
