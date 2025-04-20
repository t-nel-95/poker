package poker

import (
	"fmt"
	"sort"
)

// HandRank defines hand strength
type HandRank int

// HandRank enums
const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

var rankNames = [...]string{
	"High Card", "One Pair", "Two Pair", "Three of a Kind",
	"Straight", "Flush", "Full House", "Four of a Kind", "Straight Flush",
}

func (hr HandRank) String() string {
	return rankNames[hr]
}

// Hand represents a ranked poker hand including tie-breaking info
type Hand struct {
	Rank    HandRank
	Values  []int // Primary values for ranking
	Kickers []int // Remaining cards for tie-breaking
}

func (h Hand) String() string {
	return fmt.Sprintf("%s with %v, kickers %v", h.Rank, h.Values, h.Kickers)
}

// BestHand evaluates the best possible 5-card hand from 7 cards
func BestHand(cards []Card) Hand {
	combinations := generate5CardCombos(cards)
	bestChan := make(chan Hand, len(combinations)) // Buffered channel to collect results

	// Worker function to evaluate a combination
	evaluateCombo := func(combo []Card, ch chan<- Hand) {
		ch <- evaluateFiveCardHand(combo)
	}

	// Launch goroutines to evaluate combinations concurrently
	for _, combo := range combinations {
		go evaluateCombo(combo, bestChan)
	}

	// Collect results and determine the best hand
	var best Hand
	for range combinations {
		current := <-bestChan
		if compare := compareRankedHands(current, best); compare > 0 {
			best = current
		}
	}

	return best
}

// evaluateFiveCardHand ranks a 5-card hand
func evaluateFiveCardHand(cards []Card) Hand {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value > cards[j].Value
	})

	isFlush := true
	firstSuit := cards[0].Suit
	for _, c := range cards {
		if c.Suit != firstSuit {
			isFlush = false
			break
		}
	}

	values := make([]int, len(cards))
	for i, c := range cards {
		values[i] = int(c.Value)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	isStraight := true
	for i := 1; i < 5; i++ {
		if values[i] != values[i-1]-1 {
			isStraight = false
			break
		}
	}
	// Ace-low straight
	if !isStraight && values[0] == 14 && values[1] == 5 &&
		values[2] == 4 && values[3] == 3 && values[4] == 2 {
		isStraight = true
		values = []int{5, 4, 3, 2, 1}
	}

	counts := map[int]int{}
	for _, v := range values {
		counts[v]++
	}

	groups := make([]int, 0, 5)
	for v := range counts {
		groups = append(groups, v)
	}
	sort.Slice(groups, func(i, j int) bool {
		if counts[groups[i]] == counts[groups[j]] {
			return groups[i] > groups[j]
		}
		return counts[groups[i]] > counts[groups[j]]
	})

	var hand Hand

	switch {
	case isFlush && isStraight:
		hand = Hand{Rank: StraightFlush, Values: []int{values[0]}}
	case counts[groups[0]] == 4:
		hand = Hand{Rank: FourOfAKind, Values: []int{groups[0]}, Kickers: []int{groups[1]}}
	case counts[groups[0]] == 3 && counts[groups[1]] == 2:
		hand = Hand{Rank: FullHouse, Values: []int{groups[0], groups[1]}}
	case isFlush:
		hand = Hand{Rank: Flush, Values: values}
	case isStraight:
		hand = Hand{Rank: Straight, Values: []int{values[0]}}
	case counts[groups[0]] == 3:
		hand = Hand{Rank: ThreeOfAKind, Values: []int{groups[0]}, Kickers: groups[1:3]}
	case counts[groups[0]] == 2 && counts[groups[1]] == 2:
		hand = Hand{Rank: TwoPair, Values: []int{groups[0], groups[1]}, Kickers: []int{groups[2]}}
	case counts[groups[0]] == 2:
		hand = Hand{Rank: OnePair, Values: []int{groups[0]}, Kickers: groups[1:4]}
	default:
		hand = Hand{Rank: HighCard, Values: values[:1], Kickers: values[1:]}
	}

	return hand
}

// CompareHands compares two poker hands and returns:
// 1 if hand1 wins, -1 if hand2 wins, 0 if tie
func CompareHands(hand1Cards, hand2Cards []Card) int {
	h1 := BestHand(hand1Cards)
	h2 := BestHand(hand2Cards)
	return compareRankedHands(h1, h2)
}

// compareRankedHands compares two Hand structs
func compareRankedHands(h1, h2 Hand) int {
	if h1.Rank > h2.Rank {
		return 1
	}
	if h2.Rank > h1.Rank {
		return -1
	}
	for i := 0; i < len(h1.Values) && i < len(h2.Values); i++ {
		if h1.Values[i] > h2.Values[i] {
			return 1
		}
		if h1.Values[i] < h2.Values[i] {
			return -1
		}
	}
	for i := 0; i < len(h1.Kickers) && i < len(h2.Kickers); i++ {
		if h1.Kickers[i] > h2.Kickers[i] {
			return 1
		}
		if h1.Kickers[i] < h2.Kickers[i] {
			return -1
		}
	}
	return 0
}

// EvaluateGame determines the winner(s) among players based on the best hand
func EvaluateGame(players []Player, community []Card) []Player {
	var winners []Player

	for _, player := range players {
		if player.PlayerStatus == Folded {
			continue
		}

		hand := append([]Card{}, community...)
		player.CardStack.ForEach(func(c Card) {
			hand = append(hand, c)
		})

		if len(winners) == 0 {
			winners = []Player{player}
			continue
		}

		compare := CompareHands(hand, getCombinedHand(winners[0], community))
		if compare > 0 {
			winners = []Player{player}
		} else if compare == 0 {
			winners = append(winners, player)
		}
	}

	return winners
}

// getCombinedHand builds the 7-card hand from player and community
func getCombinedHand(p Player, community []Card) []Card {
	hand := append([]Card{}, community...)
	p.CardStack.ForEach(func(c Card) {
		hand = append(hand, c)
	})
	return hand
}

// generate5CardCombos generates all 21 possible 5-card hands from 7 cards
func generate5CardCombos(cards []Card) [][]Card {
	var combos [][]Card
	n := len(cards)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			var combo []Card
			for k := 0; k < n; k++ {
				if k != i && k != j {
					combo = append(combo, cards[k])
				}
			}
			combos = append(combos, combo)
		}
	}
	return combos
}
