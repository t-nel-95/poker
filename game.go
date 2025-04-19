package poker

import "fmt"

// Status of a game
type GameStatus int

// Game State enums
const (
	Init                GameStatus = iota 	// Create deck, shuffle
	WaitingForPlayers               		// Players join lobby, at least 2 req to proceed
	StartGame                        		// Choose dealer position, deal cards, exec small and big blind
	PreFlop                          		// First round of betting, start with player after big blind, ends when everyone is called/folded
	Flop                              		// Deal three community cards, start with first active player after dealer position, ends when everyone is called/folded
	Turn                              		// Deal one community card, start with first active player after dealer position, ends when everyone is called/folded
	River                             		// Deal one community card, start with first active player after dealer position, ends when everyone is called/folded
	DetermineWinner							// Determine who the winner(s) are
)

// Game structure
type Game struct {
	players []Player
	GameStatus
}

// NewGame creates a new game instance with initial values
func NewGame () *Game {
	fmt.Println("Starting a new game!")
	return &Game{[]Player{}, Init}
}

// AddPlayer to the game instance
func (g *Game) AddPlayer (p Player) {
	g.players = append(g.players, p)
	fmt.Printf("Player %s joined the game\n", p.Name)
}