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
