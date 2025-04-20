package poker

import (
	"fmt"
	"math/rand"
	"time"
)

// Suit of a card
type Suit int

// Suit enums
const (
	Spades   Suit = iota // 0
	Hearts               // 1
	Diamonds             // 2
	Clubs                // 3
)

// String representation of a card Suit
func (s Suit) String() string {
	return [...]string{"SPADES", "HEARTS", "DIAMONDS", "CLUBS"}[s]
}

// Unicode representation of a card Suit
func suitToUnicode(s Suit) string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("invalid card suit")
	}
}

// Value of a card
type Value int

// Value enums
const (
	Ace   Value = iota + 1 // 1
	Two                    // 2
	Three                  // 3
	Four                   // 4
	Five                   // 5
	Six                    // 6
	Seven                  // 7
	Eight                  // 8
	Nine                   // 9
	Ten                    // 10
	Jack                   // 11
	Queen                  // 12
	King                   // 13
)

// String representation of a card value
func (v Value) String() string {
	switch v {
	case Ace:
		return "ACE"
	case Two:
		return "TWO"
	case Three:
		return "THREE"
	case Four:
		return "FOUR"
	case Five:
		return "FIVE"
	case Six:
		return "SIX"
	case Seven:
		return "SEVEN"
	case Eight:
		return "EIGHT"
	case Nine:
		return "NINE"
	case Ten:
		return "TEN"
	case Jack:
		return "JACK"
	case Queen:
		return "QUEEN"
	case King:
		return "KING"
	default:
		panic("invalid card value")
	}
}

// Card data structure
type Card struct {
	Suit  Suit
	Value Value
}

// String representation of a card
func (c Card) String() string {
	return fmt.Sprintf("%s of %s %s", c.Value, c.Suit, suitToUnicode(c.Suit))
}

// CardStack is a generic stack of cards with Push and Pop functionality
type CardStack struct {
	cards []Card
}

// Push adds a card to the top of the stack
func (cs *CardStack) Push(card Card) {
	// if len(d.Cards) == 52 {
	// 	panic("The deck is already full!")
	// }
	cs.cards = append(cs.cards, card)
}

// Pop removes and returns the card from the top of the stack
func (cs *CardStack) Pop() (Card, bool) {
	if len(cs.cards) == 0 {
		return Card{}, false // empty stack
	}
	card := cs.cards[0] // Removes the first card
	cs.cards = cs.cards[1:]
	return card, true
}

// Count the number of cards in a CardStack
func (cs *CardStack) Count() int {
	return len(cs.cards)
}

// ForEach iterates over all cards in the stack and calls the provided function on each one.
// It processes cards from bottom to top (index 0 to len-1)
func (cs *CardStack) ForEach(callback func(Card)) {
	for _, card := range cs.cards {
		callback(card)
	}
}

// Deck of cards structure
type Deck struct {
	CardStack
}

// NewDeck creates a full 52-card deck, unshuffled
func NewDeck() *Deck {
	//var cards []Card
	d := &Deck{}
	for suit := Spades; suit <= Clubs; suit++ {
		for value := Ace; value <= King; value++ {
			d.Push(Card{Suit: suit, Value: value})
			//cards = append(cards, Card{Suit: suit, Value: value})
		}
	}
	return d
}

// Shuffle randomizes the order of cards in the deck
func (d *Deck) Shuffle() {
	if len(d.CardStack.cards) != 52 {
		panic(fmt.Sprintf("cannot shuffle: deck has %d cards, expected 52", len(d.CardStack.cards)))
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.CardStack.cards), func(i, j int) {
		d.CardStack.cards[i], d.CardStack.cards[j] = d.CardStack.cards[j], d.CardStack.cards[i]
	})
}
