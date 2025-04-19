package poker

import (
	"testing"
)

func TestCardString(t *testing.T) {
	card := Card{Suit: Spades, Value: Ace}
	if card.String() != "ACE of SPADES ♠" {
		t.Errorf("Expected ACE of SPADES ♠ but got %s", card.String())
	}
}

func TestDefaultDeck(t *testing.T) {
	deck := NewDeck()
	firstCard, success := deck.Pop()
	const expected = "KING of CLUBS ♣"
	if firstCard.String() != expected {
		t.Errorf("Expected first card to be %s but got %s", expected, firstCard.String())
	}
	if !success {
		t.Error("Expected deck to not be empty")
	}
}