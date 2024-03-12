package main

type Hand struct {
	Bet      uint32
	Cards    []Card
	isDouble bool
}

type Player struct {
	Credit uint32

	Hands     []Hand
	HandIndex int
	Options   []Option
}

func (g *Player) getCurrentHand() *Hand {
	if g.finished() {
		return nil
	}
	return &(g.Hands[g.HandIndex])
}

func (g *Player) finished() bool {
	return g.HandIndex >= len(g.Hands)
}

func (g *Player) reset() {
	g.Hands = []Hand{{0, []Card{}, false}}
	g.HandIndex = 0
	g.Options = nil
}
