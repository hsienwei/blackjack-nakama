package main

import (
	"fmt"
	"strings"
)

type Hand struct {
	Bet   uint32
	Cards []Card
}

type Player struct {
	Credit uint32

	Hands     []Hand
	HandIndex int
	Options   []Option
}

func (p *Player) CurrentHand() *Hand {
	if p.IsAllHandsFinished() {
		return nil
	}
	return &(p.Hands[p.HandIndex])
}

func (p *Player) actionIllegal(action Option) bool {

	actionIllegal := true
	for _, v := range player.Options {
		if v == action {
			actionIllegal = false
			break
		}
	}

	return actionIllegal

}

func (p *Player) Bet(amount uint32) {
	curHand := p.CurrentHand()
	p.Credit -= amount
	curHand.Bet = amount
}

func (p *Player) Stand() {
	p.HandIndex++
}

func (p *Player) Hit() {
	curHand := p.CurrentHand()
	curHand.Cards = dealer.DealTo(curHand.Cards, 1)
	if GetPoint(curHand.Cards).Hard > 21 {
		p.HandIndex++
	}
}

func (p *Player) Double() {
	curHand := p.CurrentHand()
	p.Credit -= curHand.Bet
	curHand.Bet *= 2
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

func (p *Player) IsAllHandsBust() bool {
	for i := 0; i < len(p.Hands); i++ {
		if GetPoint(p.Hands[i].Cards).Hard <= 21 {
			return false
		}
	}

	return true
}

func (p *Player) reset() {
	p.Hands = []Hand{}
	p.Hands = append(player.Hands, Hand{})
	p.HandIndex = 0
	p.Options = nil
}

func (p *Player) String() string {
	builder := strings.Builder{}
	builder.WriteString(">> PLAYER n")
	builder.WriteString(fmt.Sprintf("Credit: %d \n ", p.Credit))

	for i := 0; i < len(p.Hands); i++ {
		builder.WriteString(fmt.Sprintf("# %d \t Bet %d \t ", i, p.Hands[i].Bet))
		builder.WriteString(fmt.Sprintf("%v %v\n ", p.Hands[i].Cards, GetPoint(p.Hands[i].Cards)))
	}

	return builder.String()
}
