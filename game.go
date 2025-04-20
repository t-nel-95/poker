package poker

import "fmt"

// Status of a game
type GameStatus int

// Game State enums
const (
	Init              GameStatus = iota // Create deck, shuffle
	WaitingForPlayers                   // Players join lobby, at least 2 req to proceed
	StartGame                           // Choose dealer position, deal cards, exec small and big blind
	PreFlop                             // First round of betting, Initia with player after big blind, ends when everyone is called/folded
	Flop                                // Deal three community cards, Initia with first active player after dealer position, ends when everyone is called/folded
	Turn                                // Deal one community card, Initia with first active player after dealer position, ends when everyone is called/folded
	River                               // Deal one community card, Initia with first active player after dealer position, ends when everyone is called/folded
	DetermineWinner                     // Determine who the winner(s) are
)

// Game structure
type Game struct {
	Players []Player
	GameStatus
	StartingMoney int
	BigBlind      int
	DealerIndex   int
	Deck          *Deck
}

// NewGame creates a new game instance with initial values
func NewGame(startingMoney, bigBlind int) *Game {
	fmt.Println("Starting a new game!")
	deck := NewDeck()
	deck.Shuffle()
	return &Game{
		Players:       []Player{},
		GameStatus:    Init,
		StartingMoney: startingMoney,
		BigBlind:      bigBlind,
		DealerIndex:   -1,
		Deck:          deck,
	}
}

// AddPlayer to the game instance
func (g *Game) AddPlayer(name string) {
	p := NewPlayer(name, g.StartingMoney)
	g.Players = append(g.Players, *p)
	fmt.Printf("Player %s joined the game\n", p.Name)
}

// Initialise transitions the game from Init to WaitingForPlayers state
func (g *Game) Initialise() {
	if g.GameStatus != Init {
		fmt.Println("Game cannot be initialised. Current state:", g.GameStatus)
		return
	}
	g.GameStatus = WaitingForPlayers
	fmt.Println("Game has been initialised. Waiting for Players to join.")
}

// StartGame transitions the game from WaitingForPlayers to // StartGame transitions the game from WaitingForPlayers to InitiaGame

func (g *Game) StartGame() {
	if g.GameStatus != WaitingForPlayers {
		fmt.Println("Game cannot Start. Current state:", g.GameStatus)
		return
	}

	// Ensure all Players are ready
	for _, player := range g.Players {
		if !player.IsReady {
			fmt.Printf("Player %s is not ready. Cannot Start the game.\n", player.Name)
			return
		}
	}

	// Transition to StartGame state
	g.GameStatus = StartGame
	fmt.Println("Game has Started. Setting up the game...")

	// Set dealer position
	g.DealerIndex = 0
	fmt.Printf("Player %s is the dealer.\n", g.Players[g.DealerIndex].Name)

	// Deal two cards to each player
	deck := NewDeck()
	deck.Shuffle()
	for i := range g.Players {
		g.Players[i].Deal(deck)
		g.Players[i].Deal(deck)
	}

	// Handle blinds
	smallBlindIndex := -1
	bigBlindIndex := -1

	// Special case for two Players: small blind and big blind are the same player
	if len(g.Players) == 2 {
		smallBlindIndex = g.DealerIndex
		bigBlindIndex = (g.DealerIndex + 1) % len(g.Players)
	} else {
		smallBlindIndex = (g.DealerIndex + 1) % len(g.Players)
		bigBlindIndex = (g.DealerIndex + 2) % len(g.Players)
	}

	g.Players[smallBlindIndex].Raise(g.BigBlind / 2)
	fmt.Printf("Player %s posts the small blind of $%d.\n", g.Players[smallBlindIndex].Name, g.BigBlind/2)

	g.Players[bigBlindIndex].Raise(g.BigBlind)
	fmt.Printf("Player %s posts the big blind of $%d.\n", g.Players[bigBlindIndex].Name, g.BigBlind)
}
