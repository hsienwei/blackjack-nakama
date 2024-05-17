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

func (g *Player) getCurrentHand() *Hand {
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
	curHand := g.getCurrentHand()
	g.Credit -= amount
	curHand.Bet = amount
}

func (g *Player) Stand() {
	g.HandIndex++
}

func (g *Player) Hit() {
	curHand := g.getCurrentHand()
	curHand.Cards = dealer.DealTo(curHand.Cards, 1)
	if GetPoint(curHand.Cards).Hard > 21 {
		g.HandIndex++
	}
}

func (g *Player) Double() {
	curHand := g.getCurrentHand()
	g.Credit -= curHand.Bet
	curHand.Bet *= 2
	curHand.Doubled = true
	curHand.Cards = dealer.DealTo(curHand.Cards, 1)
	g.HandIndex++
}

func (g *Player) Split() {
	curHand := g.getCurrentHand()
	g.Credit -= curHand.Bet

	hand := Hand{}
	hand.Bet = curHand.Bet
	hand.Cards = []Card{curHand.Cards[1]}
	g.Hands = append(g.Hands, hand)
	curHand.Cards = []Card{curHand.Cards[0]}
}

func (g *Player) IsAllHandsFinished() bool {
	return g.HandIndex >= len(g.Hands)
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
