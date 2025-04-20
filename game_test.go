package poker

import (
	"testing"
)

func TestInitialise(t *testing.T) {
	game := NewGame(1000, 50)
	game.Initialise()

	if game.GameStatus != WaitingForPlayers {
		t.Errorf("Expected game status to be WaitingForPlayers, got %v", game.GameStatus)
	}
}

func TestAddPlayer(t *testing.T) {
	game := NewGame(1000, 50)
	game.AddPlayer("Alice")

	if len(game.Players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(game.Players))
	}

	if game.Players[0].Name != "Alice" {
		t.Errorf("Expected player name to be Alice, got %s", game.Players[0].Name)
	}
}

func TestStartGame(t *testing.T) {
	game := NewGame(1000, 50)

	game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	game.Initialise()
	game.Players[0].IsReady = true
	game.Players[1].IsReady = true
	game.StartGame()

	if game.GameStatus != StartGame {
		t.Errorf("Expected game status to be StartGame, got %v", game.GameStatus)
	}

	if len(game.Players[0].CardStack.cards) != 2 || len(game.Players[1].CardStack.cards) != 2 {
		t.Error("Expected each player to have 2 cards dealt")
	}
	if game.Players[0].bet != 25 && game.Players[1].bet != 50 {
		t.Error("Expected small blind and big blind to be posted correctly")
	}
}

func TestPreFlop(t *testing.T) {
	// Setup
	game := NewGame(1000, 50)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.Initialise()

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game
	game.StartGame()

	// PreFlop phase
	game.PreFlop()

	// Assertions
	if game.GameStatus != PreFlop {
		t.Errorf("Expected game status to be PreFlop, got %v", game.GameStatus)
	}

	if game.Community.Count() != 3 { // Updated to use Count()
		t.Errorf("Expected 3 community cards, got %d", game.Community.Count())
	}

	// Ensure community cards are not nil
	for _, card := range game.Community.cards {
		if card == (Card{}) {
			t.Error("Community card is uninitialized")
		}
	}
}

func TestTurn(t *testing.T) {
	// Setup
	game := NewGame(1000, 50)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.Initialise()

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game and transition to PreFlop
	game.StartGame()
	game.PreFlop()

	// Turn phase
	game.Turn()

	// Assertions
	if game.GameStatus != Turn {
		t.Errorf("Expected game status to be Turn, got %v", game.GameStatus)
	}

	if game.Community.Count() != 4 { // Updated to use Count()
		t.Errorf("Expected 4 community cards, got %d", game.Community.Count())
	}

	// Ensure the additional community card is not nil
	lastCard := game.Community.cards[game.Community.Count()-1]
	if lastCard == (Card{}) {
		t.Error("Last community card is uninitialized")
	}
}

func TestRiver(t *testing.T) {
	// Setup
	game := NewGame(1000, 50)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.Initialise()

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game and transition to Turn
	game.StartGame()
	game.PreFlop()
	game.Turn()

	// River phase
	game.River()

	// Assertions
	if game.GameStatus != River {
		t.Errorf("Expected game status to be River, got %v", game.GameStatus)
	}

	if game.Community.Count() != 5 { // Updated to use Count()
		t.Errorf("Expected 5 community cards, got %d", game.Community.Count())
	}

	// Ensure the additional community card is not nil
	lastCard := game.Community.cards[game.Community.Count()-1]
	if lastCard == (Card{}) {
		t.Error("Last community card is uninitialized")
	}
}

func TestAddBetsToPots(t *testing.T) {
	// Setup
	game := NewGame(1000, 50)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.AddPlayer("Charlie")
	game.Initialise()

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game
	game.StartGame()
	game.PreFlop()

	// Simulate bets using Raise method
	game.Players[0].Raise(100, game)
	game.Players[1].Call(game)
	game.Players[2].Call(game)

	// Call AddBetsToPots
	game.AddBetsToPots()

	// Assertions
	if len(game.Pots) != 1 {
		t.Errorf("Expected 1 pot, got %d", len(game.Pots))
	}

	expectedPotValue := 300
	if game.Pots[0].Amount != expectedPotValue {
		t.Errorf("Expected pot value to be %d, got %d", expectedPotValue, game.Pots[0].Amount)
	}

	// Ensure bets are cleared
	for i, player := range game.Players {
		if player.bet != 0 {
			t.Errorf("Expected player %d bet to be 0, got %d", i, player.bet)
		}
	}
}
