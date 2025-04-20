package poker

import "fmt"

// PlayerStatus of the player during a round
type PlayerStatus int

// PlayerStatus enums
const (
	Waiting PlayerStatus = iota
	Folded
	Called
	Checked
	Raised
	AllIn
	Thinking
)

func (ps PlayerStatus) String() string {
	switch ps {
	case Waiting:
		return "Waiting"
	case Folded:
		return "Folded"
	case Called:
		return "Called"
	case Checked:
		return "Checked"
	case Raised:
		return "Raised"
	case AllIn:
		return "All In"
	case Thinking:
		return "Thinking"
	default:
		panic("invalid player status value")
	}
}

// Player data structure
type Player struct {
	Name  string
	money int
	bet   int
	CardStack
	IsReady  bool
	IsDealer bool
	PlayerStatus
}

// NewPlayer initialises a new player who has joined the game
func NewPlayer(name string, money int) *Player {
	return &Player{name, money, 0, CardStack{}, false, false, Waiting}
}

// Deal a card to the player's hand from the deck
func (p *Player) Deal(d *Deck) {
	if p.CardStack.Count() == 2 {
		panic("Player's hand cannot hold more than 2 cards")
	}
	dealtCard, success := d.Pop()
	if !success {
		panic("Unable to deal to player! The deck is empty")
	}
	p.CardStack.Push(dealtCard)
}

// StartTurn sets the status of the player to reflect that it is their turn
func (p *Player) StartTurn() {
	p.PlayerStatus = Thinking
}

// Fold the player's hand, their folded bet will be collected in the pot when the round ends
func (p *Player) Fold() {
	p.PlayerStatus = Folded
	fmt.Printf("Player %s folds.\n", p.Name)
}

// Check if their current bet suffices, return whether they were allowed to check
func (p *Player) Check(g *Game) bool {
	success := false
	if p.bet == g.highestBet {
		success = true
		p.PlayerStatus = Checked
		fmt.Printf("Player %s checks.\n", p.Name)
	} else {
		fmt.Printf("Player %s cannot check.\n", p.Name)
	}
	return success
}

// AllIn sets their entire remaining money balance as their bet,
// and adds them to a split pot if they do not have enough money for the maximum bet
func (p *Player) AllIn(g *Game) {
	if p.money > 0 {
		fmt.Printf("Player %s goes All In for $%d!\n", p.Name, p.bet+p.money)
		p.bet = p.bet + p.money
		p.money = 0
		p.PlayerStatus = AllIn
		if p.bet > g.highestBet {
			g.highestBet = p.bet
		}
	}
}

// Call the bet if they can afford it, otherwise go All In
func (p *Player) Call(g *Game) {
	if p.PlayerStatus != AllIn {
		amountToCall := g.highestBet - p.bet
		if p.money <= amountToCall {
			p.AllIn(g)
		} else {
			fmt.Printf("Player %s calls $%d\n", p.Name, g.highestBet)
			p.bet = p.bet + amountToCall
			p.money = p.money - amountToCall
			p.PlayerStatus = Called
			if p.bet > g.highestBet {
				g.highestBet = p.bet
			}
		}
	} else {
		fmt.Printf("Player %s is already All In", p.Name)
	}
}

// Raise by a specified amount if the player has suffient money.
// If it's the same as their amount of money, go All In
// Return whether the bet was successful
func (p *Player) Raise(amount int, g *Game) bool {
	success := false
	if amount < p.money {
		p.bet = p.bet + amount
		p.money = p.money - amount
		success = true
		fmt.Printf("Player %s raised by $%d\n", p.Name, amount)
	}
	if amount == p.money {
		p.AllIn(g)
		success = true
	}
	if amount > p.money {
		fmt.Printf("Player %s tried to raise by %d but their balance is insufficient. You can go All In instead and create a split pot.\n", p.Name, amount)
	}
	if success && p.bet > g.highestBet {
		g.highestBet = p.bet
	}
	fmt.Printf("Highest bet is %d\n", g.highestBet)
	return success
}
