package poker

import (
	"testing"
)

func TestPlayerDeal(t *testing.T) {
	tests := []struct {
		name      string
		deckCards []Card
		expected1 Card
		expected2 Card
	}{
		{
			name: "Deal two cards",
			deckCards: []Card{
				{Clubs, King},
				{Clubs, Queen},
			},
			expected1: Card{Clubs, King},
			expected2: Card{Clubs, Queen},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deck{CardStack{tt.deckCards}}
			p := NewPlayer("Bob", 1000)
			p.Deal(d)
			p.Deal(d)

			if p.CardStack.cards[0] != tt.expected1 {
				t.Errorf("Expected first card %s, got %s", tt.expected1, p.CardStack.cards[0])
			}
			if p.CardStack.cards[1] != tt.expected2 {
				t.Errorf("Expected second card %s, got %s", tt.expected2, p.CardStack.cards[1])
			}
		})
	}
}

func TestRaise(t *testing.T) {
	tests := []struct {
		name          string
		initialMoney  int
		raiseAmount   int
		expectedBet   int
		expectedMoney int
		success       bool
	}{
		{
			name:          "Sufficient balance to raise",
			initialMoney:  1000,
			raiseAmount:   100,
			expectedBet:   150, // Includes small blind (50)
			expectedMoney: 850, // Deducts small blind and raise
			success:       true,
		},
		{
			name:          "Insufficient balance to raise",
			initialMoney:  1000,
			raiseAmount:   2000,
			expectedBet:   50,  // Only small blind is posted
			expectedMoney: 950, // Deducts only small blind
			success:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(tt.initialMoney, 100) // Big blind is 100, small blind is 50
			g.Initialise()
			g.AddPlayer("Bob")
			g.AddPlayer("Alice")
			p := g.getPlayer("Bob", 0)
			p.IsReady = true
			g.Players[1].IsReady = true
			g.StartGame()
			g.PreFlop()

			success := p.Raise(tt.raiseAmount, g)

			if success != tt.success {
				t.Errorf("Expected success to be %v, got %v", tt.success, success)
			}
			if p.bet != tt.expectedBet {
				t.Errorf("Expected bet to be %d, got %d", tt.expectedBet, p.bet)
			}
			if p.money != tt.expectedMoney {
				t.Errorf("Expected money to be %d, got %d", tt.expectedMoney, p.money)
			}
		})
	}
}

func TestPlayerCheck(t *testing.T) {
	tests := []struct {
		name                 string
		previousPlayerRaises bool
		expectedCheck        bool
	}{
		{
			name:                 "Player may check",
			previousPlayerRaises: false,
			expectedCheck:        true,
		},
		{
			name:                 "Player may not check",
			previousPlayerRaises: true,
			expectedCheck:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(1000, 50)
			g.Initialise()
			g.AddPlayer("Bob")
			g.AddPlayer("Alice")
			p1 := g.getPlayer("Bob", 0)   // Small blind
			p2 := g.getPlayer("Alice", 1) // Big blind
			p1.IsReady = true
			p2.IsReady = true
			g.StartGame()
			g.PreFlop()
			if tt.previousPlayerRaises {
				p1.Raise(100, g) // Simulate a raise from the previous player
			} else {
				p1.Call(g) // Call the big blind
			}
			result := p2.Check(g)

			if result != tt.expectedCheck {
				t.Errorf("Expected check result to be %v, got %v", tt.expectedCheck, result)
			}
		})
	}
}

func TestPlayerFold(t *testing.T) {
	tests := []struct {
		name           string
		initialMoney   int
		responseBet    int
		expectedMoney  int
		expectedStatus PlayerStatus
	}{
		{
			name:           "Player folds with bet",
			initialMoney:   1000,
			responseBet:    50,
			expectedMoney:  900,
			expectedStatus: Folded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(1000, 50)
			g.Initialise()
			g.AddPlayer("Bob")
			g.AddPlayer("Alice")
			p1 := g.getPlayer("Bob", 0)   // Small blind
			p2 := g.getPlayer("Alice", 1) // Big blind
			// Ensure players are ready
			for i := range g.Players {
				g.Players[i].IsReady = true
			}
			g.StartGame()
			g.PreFlop()
			p1.Call(g)                  // Call the big blind
			p2.Raise(tt.responseBet, g) // Raise to $100
			p1.AllIn(g)
			p2.Fold() // No longer capturing a return value

			if p2.money != tt.expectedMoney {
				t.Errorf("Expected money to be %d, got %d", tt.expectedMoney, p2.money)
			}
			if p2.PlayerStatus != tt.expectedStatus {
				t.Errorf("Expected status to be %v, got %v", tt.expectedStatus, p2.PlayerStatus)
			}
		})
	}
}

func TestPlayerCall(t *testing.T) {
	tests := []struct {
		name           string
		initialMoney   int
		responseBet    int
		highestBet     int
		expectedBet    int
		expectedMoney  int
		expectedStatus PlayerStatus
	}{
		{
			name:           "Call with sufficient money",
			initialMoney:   1000,
			responseBet:    0,
			highestBet:     100,
			expectedBet:    100,
			expectedMoney:  900,
			expectedStatus: Called,
		},
		{
			name:           "Call with existing bet",
			initialMoney:   900,
			responseBet:    100,
			highestBet:     200,
			expectedBet:    200,
			expectedMoney:  800,
			expectedStatus: Called,
		},
		{
			name:           "Call goes All In",
			initialMoney:   800,
			responseBet:    200,
			highestBet:     1000,
			expectedBet:    1000,
			expectedMoney:  0,
			expectedStatus: AllIn,
		},
		{
			name:           "Call with insufficient money (All In)",
			initialMoney:   1000,
			responseBet:    0,
			highestBet:     2000,
			expectedBet:    1000,
			expectedMoney:  0,
			expectedStatus: AllIn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGame(1000, 50)
			p := NewPlayer("TestPlayer", tt.initialMoney)
			p.bet = tt.responseBet
			g.highestBet = tt.highestBet

			p.Call(g)

			if p.bet != tt.expectedBet {
				t.Errorf("Expected bet to be %d, got %d", tt.expectedBet, p.bet)
			}
			if p.money != tt.expectedMoney {
				t.Errorf("Expected money to be %d, got %d", tt.expectedMoney, p.money)
			}
			if p.PlayerStatus != tt.expectedStatus {
				t.Errorf("Expected status to be %s, got %s", tt.expectedStatus, p.PlayerStatus)
			}
		})
	}
}
