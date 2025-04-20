package poker

import (
	"testing"
)

func TestPlayerDeal(t *testing.T) {
	d := NewDeck()
	p := NewPlayer("Bob", 1000)
	p.Deal(d)
	p.Deal(d)
	expected1 := Card{Clubs, King}
	expected2 := Card{Clubs, Queen}
	if p.CardStack.cards[0] != expected1 {
		t.Errorf("First expected card dealt to player from deck is %s but found %s", expected1, p.CardStack.cards[0])
		t.Errorf("Second expected card dealt to player from deck is %s but found %s", expected2, p.CardStack.cards[0])

	}
}

func TestRaise(t *testing.T) {
	g := NewGame(1000, 50)     // Create a new game instance with starting money and big blind
	g.AddPlayer("Bob")         // Add a player to the game
	p := g.getPlayer("Bob", 0) // Get the player instance from the game
	success := g.Players[0].Raise(100, g)
	if !success {
		t.Error("Expected player with a sufficent balance to be able to Raise")
	}
	if p.money != 900 {
		t.Errorf("Expected balance of $900 after raising by $100, but got $%d", p.money)
	}
	if p.bet != 100 {
		t.Errorf("Expected bet to be $100 after raising by $100 but got $%d", p.bet)
	}
	success = p.Raise(2000, g)
	if success {
		t.Error("Expected player with insufficient balance to be unable to Raise")
	}
}

func TestPlayerCheck(t *testing.T) {
	maxBet := 100
	startingMoney := 1000
	bigBlind := 50
	g := NewGame(startingMoney, bigBlind)
	g.AddPlayer("Bob")
	g.AddPlayer("Mike")
	g.Players[0].Raise(100, g)
	if g.Players[1].Check(maxBet) {
		t.Error("Expected a player with insufficent bet to be unable to check")
	}
	if !g.Players[0].Check(maxBet) {
		t.Error("Expected a player with sufficent bet to be able to check")
	}
}

func TestPlayerFold(t *testing.T) {
	g := NewGame(1000, 50)     // Create a new game instance with starting money and big blind
	g.AddPlayer("Bob")         // Add a player to the game
	p := g.getPlayer("Bob", 0) // Get the player instance from the game
	p.Raise(100, g)
	amountForfeit := p.Fold()
	if amountForfeit != 100 {
		t.Errorf("Expected a player with a bet of $100 to fold $100 but instead got $%d", amountForfeit)
	}
	if p.money != 900 {
		t.Errorf("Expected a player who folds a $100 bet from a $1000 starting balance to have $900 remaining but got $%d", p.money)
	}
}

func TestPlayerCall(t *testing.T) {
	g := NewGame(1000, 50) // Create a new game instance with starting money and big blind
	p := NewPlayer("Bob", 1000)
	g.highestBet = 100
	// Case 1
	p.Call(g)
	if p.bet != 100 {
		t.Errorf("Expected a player with no bet calling a $100 bet to have a resulting bet of $100 but got $%d", p.bet)
	}
	if p.money != 900 {
		t.Errorf("Expected a player with a call of $100 from a initial balance of $1000 to have a resulting balance of $900 but got %d", p.money)
	}
	//Case 2
	g.highestBet = 200
	p.Call(g)
	if p.bet != 200 {
		t.Errorf("Expected a player with an existing bet of $100 calling a $200 bet to have a resulting bet of $200 but got $%d", p.bet)
	}
	if p.money != 800 {
		t.Errorf("Expected a player with a call of $200 from a initial balance of $1000 to have a resulting balance of $800 but got %d", p.money)
	}
	// Case 3
	g.highestBet = 1000
	p.Call(g)
	if p.bet != 1000 {
		t.Errorf("Expected a player with an existing bet of $200 calling a $1000 bet to have a resulting bet of $1000 but got $%d", p.bet)
	}
	if p.money != 0 {
		t.Errorf("Expected a player with a call of $1000 from a initial balance of $1000 to have a resulting balance of $0 but got %d", p.money)
	}
	if p.PlayerStatus != AllIn {
		t.Errorf("Expected a player with a call of $1000 from a initial balance of $1000 to have status of All in but got %s", p.PlayerStatus)
	}
	// Case 4
	p = NewPlayer("Mike", 1000)
	g.highestBet = 2000
	p.Call(g)
	if p.bet != 1000 {
		t.Errorf("Expected a player with an existing bet of $0 calling a $2000 bet with a balance of $1000 to have a resulting bet of $1000 but got $%d", p.bet)
	}
	if p.money != 0 {
		t.Errorf("Expected a player with a call of $2000 from a initial balance of $1000 to have a resulting balance of $0 but got %d", p.money)
	}
	if p.PlayerStatus != AllIn {
		t.Errorf("Expected a player with a call of $2000 from a initial balance of $1000 to have status of All in but got %s", p.PlayerStatus)
	}
}
