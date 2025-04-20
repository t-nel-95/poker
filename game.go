package poker

import "fmt"

// Status of a game
type GameStatus int

// Game State enums
const (
	Init              GameStatus = iota // Create deck, shuffle
	WaitingForPlayers                   // Players join lobby, at least 2 req to proceed
	StartGame                           // Choose dealer position, deal cards, exec small and big blind
	PreFlop                             // First round of betting, Initia with player after big blind, ends when everyone is called/folded
	Flop                                // Deal three community cards, Initia with first active player after dealer position, ends when everyone is called/folded
	Turn                                // Deal one community card, Initia with first active player after dealer position, ends when everyone is called/folded
	River                               // Deal one community card, Initia with first active player after dealer position, ends when everyone is called/folded
	DetermineWinner                     // Determine who the winner(s) are
)

// Pot represents the accumulated money from a round of betting and the eligible players
type Pot struct {
	Amount   int
	Eligible []Player
}

// Game structure
type Game struct {
	Players []Player
	GameStatus
	StartingMoney int
	BigBlind      int
	DealerIndex   int
	Deck          *Deck
	Community     CardStack
	Pots          []Pot // Main pot and optional side pots
	highestBet    int   // Tracks the current highest bet during the game
}

// NewGame creates a new game instance with initial values
func NewGame(startingMoney, bigBlind int) *Game {
	fmt.Println("Starting a new game!")
	deck := NewDeck()
	deck.Shuffle()
	return &Game{
		Players:       []Player{},
		GameStatus:    Init,
		StartingMoney: startingMoney,
		BigBlind:      bigBlind,
		DealerIndex:   -1,
		Deck:          deck,
	}
}

// AddPlayer to the game instance
func (g *Game) AddPlayer(name string) {
	p := NewPlayer(name, g.StartingMoney)
	g.Players = append(g.Players, *p)
	fmt.Printf("Player %s joined the game\n", p.Name)
}

// Initialise transitions the game from Init to WaitingForPlayers state
func (g *Game) Initialise() {
	if g.GameStatus != Init {
		fmt.Println("Game cannot be initialised. Current state:", g.GameStatus)
		return
	}
	g.GameStatus = WaitingForPlayers
	fmt.Println("Game has been initialised. Waiting for Players to join.")
}

// StartGame transitions the game from WaitingForPlayers to StartGame
func (g *Game) StartGame() {
	if g.GameStatus != WaitingForPlayers {
		fmt.Println("Game cannot Start. Current state:", g.GameStatus)
		return
	}

	// Ensure all Players are ready
	for _, player := range g.Players {
		if !player.IsReady {
			fmt.Printf("Player %s is not ready. Cannot Start the game.\n", player.Name)
			return
		}
	}

	// Transition to StartGame state
	g.GameStatus = StartGame
	fmt.Println("Game has Started. Setting up the game...")

	// Set dealer position
	g.DealerIndex = 0
	fmt.Printf("Player %s is the dealer.\n", g.Players[g.DealerIndex].Name)

	// Deal two cards to each player
	deck := NewDeck()
	deck.Shuffle()
	for i := range g.Players {
		g.Players[i].Deal(deck)
		g.Players[i].Deal(deck)
	}

	// Handle blinds
	smallBlindIndex := -1
	bigBlindIndex := -1

	// Special case for two Players: small blind and big blind are the same player
	if len(g.Players) == 2 {
		smallBlindIndex = g.DealerIndex
		bigBlindIndex = (g.DealerIndex + 1) % len(g.Players)
	} else {
		smallBlindIndex = (g.DealerIndex + 1) % len(g.Players)
		bigBlindIndex = (g.DealerIndex + 2) % len(g.Players)
	}

	g.Players[smallBlindIndex].Raise(g.BigBlind/2, g)
	fmt.Printf("Player %s posts the small blind of $%d.\n", g.Players[smallBlindIndex].Name, g.BigBlind/2)

	g.Players[bigBlindIndex].Raise(g.BigBlind, g)
	fmt.Printf("Player %s posts the big blind of $%d.\n", g.Players[bigBlindIndex].Name, g.BigBlind)

	// Initialize the Community CardStack
	g.Community = CardStack{}

	// Initialize the main pot
	g.Pots = []Pot{{Amount: 0, Eligible: g.Players}}
}

// AddBetsToPots adds the current bets of players to the pots
func (g *Game) AddBetsToPots() {
	for i := range g.Players {
		player := &g.Players[i]

		// Skip players with no bets
		if player.bet == 0 {
			continue
		}

		switch player.PlayerStatus {
		case Called:
			// Add their bet to the latest side pot they are eligible for, if side pots exist
			addedToSidePot := false
			for j := len(g.Pots) - 1; j > 0; j-- { // Skip the main pot (index 0)
				if containsPlayer(g.Pots[j].Eligible, *player) {
					g.Pots[j].Amount += player.bet
					player.bet = 0
					addedToSidePot = true
					break
				}
			}
			// If no side pot exists or they are not eligible for any side pot, add to the main pot
			if !addedToSidePot {
				g.Pots[0].Amount += player.bet
				player.bet = 0
			}

		case Folded:
			// Add their bet to the latest side pot they were eligible for
			addedToSidePot := false
			for j := len(g.Pots) - 1; j >= 0; j-- {
				if containsPlayer(g.Pots[j].Eligible, *player) {
					g.Pots[j].Amount += player.bet
					player.bet = 0
					addedToSidePot = true
					break
				}
			}
			// If not eligible for any side pot, add to the main pot
			if !addedToSidePot {
				g.Pots[0].Amount += player.bet
				player.bet = 0
			}

		case AllIn:
			// Check if the player is already in a side pot
			alreadyInSidePot := false
			for _, pot := range g.Pots {
				if containsPlayer(pot.Eligible, *player) {
					alreadyInSidePot = true
					break
				}
			}

			if alreadyInSidePot && player.bet == 0 {
				// If already in a side pot and no new bets, skip
				continue
			}

			// Determine if a side pot is needed
			highestBet := 0
			for _, p := range g.Players {
				if p.bet > highestBet {
					highestBet = p.bet
				}
			}

			if player.bet < highestBet {
				// Create a side pot for the excess bets
				sidePot := Pot{Amount: 0, Eligible: []Player{}}
				for j := range g.Players {
					otherPlayer := &g.Players[j]
					if otherPlayer.bet > player.bet {
						excess := otherPlayer.bet - player.bet
						sidePot.Amount += excess
						otherPlayer.bet -= excess
						sidePot.Eligible = append(sidePot.Eligible, *otherPlayer)
					}
				}
				g.Pots = append(g.Pots, sidePot)
			}

			// Add the All In player's bet to the main pot
			g.Pots[0].Amount += player.bet
			player.bet = 0
		}
	}

	// Handle the case where only one player remains active
	activePlayers := 0
	for _, player := range g.Players {
		if player.PlayerStatus != Folded {
			activePlayers++
		}
	}
	if activePlayers == 1 {
		fmt.Println("Only one player remains active. They win the pots.")
		for _, pot := range g.Pots {
			for i := range g.Players {
				if g.Players[i].PlayerStatus != Folded {
					g.Players[i].money += pot.Amount
					pot.Amount = 0
					break
				}
			}
		}
	}
}

// containsPlayer checks if a player is in the eligible list for a pot
func containsPlayer(players []Player, player Player) bool {
	for _, p := range players {
		if p.Name == player.Name {
			return true
		}
	}
	return false
}

// PreFlop transitions the game to the PreFlop state and deals three community cards
func (g *Game) PreFlop() {
	if g.GameStatus != StartGame {
		fmt.Println("Game cannot transition to PreFlop. Current state:", g.GameStatus)
		return
	}

	// Transition to PreFlop state
	g.GameStatus = PreFlop
	fmt.Println("Transitioning to PreFlop phase...")

	// Deal three cards to the Community
	for i := 0; i < 3; i++ {
		card, success := g.Deck.Pop()
		if !success {
			panic("Deck is empty! Cannot deal community cards.")
		}
		g.Community.Push(card)
	}

	// Print the Community cards
	fmt.Println("Community cards dealt:")
	g.Community.ForEach(func(c Card) {
		fmt.Println(c)
	})

	// Add bets to pots if all players have called, folded, or gone all in
	g.AddBetsToPots()
}

// Turn transitions the game to the Turn state and deals one additional community card
func (g *Game) Turn() {
	if g.GameStatus != PreFlop {
		fmt.Println("Game cannot transition to Turn. Current state:", g.GameStatus)
		return
	}

	// Transition to Turn state
	g.GameStatus = Turn
	fmt.Println("Transitioning to Turn phase...")

	// Deal one card to the Community
	card, success := g.Deck.Pop()
	if !success {
		panic("Deck is empty! Cannot deal community card.")
	}
	g.Community.Push(card)

	// Print the Community cards
	fmt.Println("Community cards dealt:")
	g.Community.ForEach(func(c Card) {
		fmt.Println(c)
	})

	// Add bets to pots if all players have called, folded, or gone all in
	g.AddBetsToPots()
}

// River transitions the game to the River state and deals one final community card
func (g *Game) River() {
	if g.GameStatus != Turn {
		fmt.Println("Game cannot transition to River. Current state:", g.GameStatus)
		return
	}

	// Transition to River state
	g.GameStatus = River
	fmt.Println("Transitioning to River phase...")

	// Deal one card to the Community
	card, success := g.Deck.Pop()
	if !success {
		panic("Deck is empty! Cannot deal community card.")
	}
	g.Community.Push(card)

	// Print the Community cards
	fmt.Println("Community cards dealt:")
	g.Community.ForEach(func(c Card) {
		fmt.Println(c)
	})

	// Add bets to pots if all players have called, folded, or gone all in
	g.AddBetsToPots()
}

// DetermineWinner transitions the game to the DetermineWinner state, evaluates player hands, and announces the winner(s)
func (g *Game) DetermineWinner() {
	if g.GameStatus != River {
		fmt.Println("Game cannot transition to DetermineWinner. Current state:", g.GameStatus)
		return
	}

	// Transition to DetermineWinner state
	g.GameStatus = DetermineWinner
	fmt.Println("Determining the winner(s)...")

	// Evaluate hands and determine winners
	winners := EvaluateGame(g.Players, g.Community.cards)

	// Announce winners
	if len(winners) == 1 {
		fmt.Printf("The winner is %s!\n", winners[0].Name)
	} else {
		fmt.Print("The winners are: ")
		for i, winner := range winners {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(winner.Name)
		}
		fmt.Println("!")
	}

	// Distribute pot winnings to winners
	for _, pot := range g.Pots {
		winners := EvaluateGame(pot.Eligible, g.Community.cards)
		winnings := pot.Amount / len(winners)
		for _, winner := range winners {
			for i := range g.Players {
				if g.Players[i].Name == winner.Name {
					g.Players[i].money += winnings
					break
				}
			}
		}
	}

	// Print final player balances
	fmt.Println("Final player balances:")
	for _, player := range g.Players {
		fmt.Printf("%s: $%d\n", player.Name, player.money)
	}
}

func (g *Game) PlayerRaise(playerIndex, amount int) {
	if g.Players[playerIndex].Raise(amount, g) {
		// highestBet is updated within Raise
	}
}

// getPlayer retrieves a pointer to a player by name or index.
// If both name and index are provided, index takes precedence.
// Returns nil if the player is not found.
func (g *Game) getPlayer(name string, index int) *Player {
	if index >= 0 && index < len(g.Players) {
		return &g.Players[index]
	}
	for i := range g.Players {
		if g.Players[i].Name == name {
			return &g.Players[i]
		}
	}
	return nil
}
