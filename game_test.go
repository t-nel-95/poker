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

	// Ensure no community cards are dealt
	if game.Community.Count() != 0 {
		t.Errorf("Expected 0 community cards, got %d", game.Community.Count())
	}
}

func TestTurn(t *testing.T) {
	// Setup
	game := NewGame(1000, 50)
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	p1 := game.getPlayer("Alice", 0) // Small blind
	p2 := game.getPlayer("Bob", 1)   // Big blind
	game.Initialise()

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game and transition to PreFlop
	game.StartGame()
	game.PreFlop()
	p1.Call(game)  // Call the big blind
	p2.Check(game) // Check the big blind
	game.Flop()

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
	p1 := game.getPlayer("Alice", 0) // Small blind
	p2 := game.getPlayer("Bob", 1)   // Big blind
	game.Initialise()

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game and transition to Turn
	game.StartGame()
	game.PreFlop()
	p1.Call(game)  // Call the big blind
	p2.Check(game) // Check the big blind
	game.Flop()
	p1.Check(game) // Check the flop
	p2.Check(game) // Check the flop
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
	game.Initialise()
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")
	game.AddPlayer("Charlie")
	p1 := game.getPlayer("Alice", 0)   // Dealer
	p2 := game.getPlayer("Bob", 1)     // Small blind
	p3 := game.getPlayer("Charlie", 2) // Big blind

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game
	game.StartGame()
	game.PreFlop()

	// Call the blinds
	p1.Call(game)      // Call the big blind
	p2.Call(game)      // Call the small blind
	p3.Raise(50, game) // Raise to $100
	p1.Call(game)      // Call the raise
	p2.Call(game)      // Call the raise
	p3.Check(game)     // Check the raise

	// Debug
	for i, player := range game.Players {
		t.Logf("Player %d (%s) bet: %d", i, player.Name, player.bet)
	}

	// Call AddBetsToPots
	game.AddBetsToPots()

	// Debug
	for i, player := range game.Players {
		t.Logf("Player %d (%s) bet: %d", i, player.Name, player.bet)
	}

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

func TestFlop(t *testing.T) {
	// Setup
	game := NewGame(1000, 50)
	game.Initialise()
	game.AddPlayer("Alice")
	game.AddPlayer("Bob")

	// Ensure players are ready
	for i := range game.Players {
		game.Players[i].IsReady = true
	}

	// Start the game and transition to PreFlop
	game.StartGame()
	p1 := game.getPlayer("Alice", 0) // Small blind
	p1.Call(game)                    // Call the big blind
	game.PreFlop()

	// Flop phase
	game.Flop()

	// Assertions
	if game.GameStatus != Flop {
		t.Errorf("Expected game status to be Flop, got %v", game.GameStatus)
	}

	if game.Community.Count() != 3 {
		t.Errorf("Expected 3 community cards, got %d", game.Community.Count())
	}

	// Ensure the community cards are not nil
	for _, card := range game.Community.cards {
		if card == (Card{}) {
			t.Error("Community card is uninitialized")
		}
	}
}
