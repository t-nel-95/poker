package poker

import (
	"testing"
)

func card(value Value, suit Suit) Card {
	return Card{Value: value, Suit: suit}
}

func TestBestHandRankings(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected HandRank
	}{
		{
			"Straight Flush",
			[]Card{
				card(Nine, Hearts), card(Ten, Hearts), card(Jack, Hearts),
				card(Queen, Hearts), card(King, Hearts), card(Ace, Spades), card(Two, Clubs),
			},
			StraightFlush,
		},
		{
			"Four of a Kind",
			[]Card{
				card(Ten, Spades), card(Ten, Hearts), card(Ten, Clubs),
				card(Ten, Diamonds), card(Ace, Hearts), card(Five, Spades), card(Three, Clubs),
			},
			FourOfAKind,
		},
		{
			"Full House",
			[]Card{
				card(Queen, Spades), card(Queen, Hearts), card(Queen, Clubs),
				card(Jack, Diamonds), card(Jack, Hearts), card(Two, Spades), card(Nine, Clubs),
			},
			FullHouse,
		},
		{
			"Flush",
			[]Card{
				card(Two, Spades), card(Five, Spades), card(Eight, Spades),
				card(Jack, Spades), card(King, Spades), card(Three, Clubs), card(Four, Hearts),
			},
			Flush,
		},
		{
			"Straight (Ace low)",
			[]Card{
				card(Ace, Diamonds), card(Two, Spades), card(Three, Clubs),
				card(Four, Hearts), card(Five, Diamonds), card(Ten, Spades), card(King, Clubs),
			},
			Straight,
		},
		{
			"Three of a Kind",
			[]Card{
				card(Nine, Spades), card(Nine, Hearts), card(Nine, Diamonds),
				card(Two, Clubs), card(Five, Hearts), card(Six, Spades), card(Jack, Clubs),
			},
			ThreeOfAKind,
		},
		{
			"Two Pair",
			[]Card{
				card(Ace, Spades), card(Ace, Clubs), card(King, Hearts),
				card(King, Diamonds), card(Five, Clubs), card(Nine, Spades), card(Three, Hearts),
			},
			TwoPair,
		},
		{
			"One Pair",
			[]Card{
				card(Queen, Spades), card(Queen, Diamonds), card(Ten, Clubs),
				card(Nine, Hearts), card(Four, Spades), card(Two, Diamonds), card(Six, Clubs),
			},
			OnePair,
		},
		{
			"High Card",
			[]Card{
				card(Two, Clubs), card(Four, Hearts), card(Six, Diamonds),
				card(Nine, Spades), card(Jack, Clubs), card(Queen, Hearts), card(Ace, Diamonds),
			},
			HighCard,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := BestHand(test.cards)
			if result.Rank != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result.Rank)
			}
		})
	}
}

func TestCompareHands(t *testing.T) {
	hand1 := []Card{
		card(Ten, Spades), card(Ten, Hearts), card(Ten, Diamonds),
		card(Three, Clubs), card(Five, Hearts), card(Nine, Spades), card(King, Clubs),
	}

	hand2 := []Card{
		card(Ace, Spades), card(Ace, Clubs), card(Eight, Hearts),
		card(Four, Diamonds), card(Six, Spades), card(Jack, Clubs), card(Two, Hearts),
	}

	result := CompareHands(hand1, hand2)
	if result != 1 {
		t.Errorf("expected hand1 to win, got result %d", result)
	}
}

func TestEvaluateGame(t *testing.T) {
	community := []Card{
		card(Ten, Spades), card(Jack, Hearts), card(Queen, Spades),
		card(King, Spades), card(Nine, Spades),
	}

	p1 := *NewPlayer("Alice", 1000)
	p1.CardStack.Push(card(Eight, Spades)) // makes straight flush
	p1.CardStack.Push(card(Seven, Spades))

	p2 := *NewPlayer("Bob", 1000)
	p2.CardStack.Push(card(Ace, Diamonds))
	p2.CardStack.Push(card(Two, Clubs))

	players := []Player{p1, p2}
	winners := EvaluateGame(players, community)

	if len(winners) != 1 || winners[0].Name != "Alice" {
		t.Errorf("expected Alice to win, got %v", winners)
	}
}
