package poker

import (
	"testing"
)

func TestCardString(t *testing.T) {
	tests := []struct {
		name     string
		card     Card
		expected string
	}{
		{
			name:     "Ace of Spades",
			card:     Card{Suit: Spades, Value: Ace},
			expected: "ACE of SPADES ♠",
		},
		{
			name:     "King of Hearts",
			card:     Card{Suit: Hearts, Value: King},
			expected: "KING of HEARTS ♥",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.card.String() != tt.expected {
				t.Errorf("Expected %s but got %s", tt.expected, tt.card.String())
			}
		})
	}
}

func TestDefaultDeck(t *testing.T) {
	deck := NewDeck()
	numCards := deck.CardStack.Count()
	if numCards != 52 {
		t.Errorf("Number of cards should be 52 but is %d", numCards)
	}

	firstCard, success := deck.Pop()
	expected := Card{Suit: Spades, Value: Ace} // First card based on Pop logic
	if firstCard != expected {
		t.Errorf("Expected first card to be %+v but got %+v", expected, firstCard)
	}
	if !success {
		t.Error("Expected deck to not be empty")
	}
}
