package main

import (
	"fmt"
	"strings"
)

type Hand struct {
	Bet     uint32
	Cards   []Card
	Doubled bool
}

type Player struct {
	Credit uint32

	Hands     []Hand
	HandIndex int
	Options   []Option
}

func (g *Player) CurrentHand() *Hand {
	if g.IsAllHandsFinished() {
		return nil
	}
	return &(g.Hands[g.HandIndex])
}

func (g *Player) actionIllegal(action Option) bool {

	actionIllegal := true
	for _, v := range player.Options {
		if v == action {
			actionIllegal = false
			break
		}
	}

	return actionIllegal

}

func (g *Player) Bet(amount uint32) {
	curHand := g.CurrentHand()
	g.Credit -= amount
	curHand.Bet = amount
}

func (g *Player) Stand() {
	g.HandIndex++
}

func (g *Player) Hit() {
	curHand := g.CurrentHand()
	curHand.Cards = dealer.DealTo(curHand.Cards, 1)
	if GetPoint(curHand.Cards).Hard > 21 {
		g.HandIndex++
	}
}

func (p *Player) Double() {
	curHand := p.CurrentHand()
	p.Credit -= curHand.Bet
	curHand.Bet *= 2
	curHand.Doubled = true
	curHand.Cards = dealer.DealTo(curHand.Cards, 1)
	p.HandIndex++
}

func (p *Player) Split() {
	p.Credit -= p.CurrentHand().Bet

	hand := Hand{}
	hand.Bet = p.CurrentHand().Bet
	hand.Cards = []Card{p.CurrentHand().Cards[1]}
	p.Hands = append(p.Hands, hand)

	p.CurrentHand().Cards = []Card{p.CurrentHand().Cards[0]}
}

func (p *Player) IsAllHandsFinished() bool {
	return p.HandIndex >= len(p.Hands)
}

func (g *Player) IsAllHandsBust() bool {
	for i := 0; i < len(g.Hands); i++ {
		if GetPoint(g.Hands[i].Cards).Hard <= 21 {
			return false
		}
	}

	return true
}

func (g *Player) reset() {
	g.Hands = []Hand{}
	g.Hands = append(player.Hands, Hand{})
	g.HandIndex = 0
	g.Options = nil
}

func (g *Player) String() string {
	builder := strings.Builder{}
	builder.WriteString(">> PLAYER n")
	builder.WriteString(fmt.Sprintf("Credit: %d \n ", g.Credit))

	for i := 0; i < len(g.Hands); i++ {
		builder.WriteString(fmt.Sprintf("# %d \t Bet %d \t ", i, g.Hands[i].Bet))
		builder.WriteString(fmt.Sprintf("%v %v\n ", g.Hands[i].Cards, GetPoint(g.Hands[i].Cards)))
	}

	return builder.String()
}
