package main

import (
	"fmt"
	"strings"
)

type Hand struct {
	Bet uint32
	CardSet
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
	for _, v := range p.Options {
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

func (p *Player) Hit(dealer *Dealer) {
	curHand := p.CurrentHand()
	curHand.CardSet = dealer.DealTo(curHand.CardSet, 1)
	if curHand.CardSet.GetPoint().Hard > 21 {
		p.HandIndex++
	}
}

func (p *Player) Double(dealer *Dealer) {
	curHand := p.CurrentHand()
	p.Credit -= curHand.Bet
	curHand.Bet *= 2
	curHand.CardSet = dealer.DealTo(curHand.CardSet, 1)
	p.HandIndex++
}

func (p *Player) Split() {
	p.Credit -= p.CurrentHand().Bet

	hand := Hand{}
	hand.Bet = p.CurrentHand().Bet
	hand.CardSet = []Card{p.CurrentHand().CardSet[1]}
	p.Hands = append(p.Hands, hand)

	p.CurrentHand().CardSet = []Card{p.CurrentHand().CardSet[0]}
}

func (p *Player) IsAllHandsFinished() bool {
	return p.HandIndex >= len(p.Hands)
}

func (p *Player) IsAllHandsBust() bool {
	for i := 0; i < len(p.Hands); i++ {
		if p.Hands[i].CardSet.GetPoint().Hard <= 21 {
			return false
		}
	}

	return true
}

func (p *Player) reset() {
	p.Hands = []Hand{}
	p.Hands = append(p.Hands, Hand{})
	p.HandIndex = 0
	p.Options = nil
}

func (p *Player) String() string {
	builder := strings.Builder{}
	builder.WriteString(">> PLAYER n")
	builder.WriteString(fmt.Sprintf("Credit: %d \n ", p.Credit))

	for i := 0; i < len(p.Hands); i++ {
		builder.WriteString(fmt.Sprintf("# %d \t Bet %d \t ", i, p.Hands[i].Bet))
		builder.WriteString(fmt.Sprintf("%v %v\n ", p.Hands[i].CardSet, p.Hands[i].CardSet.GetPoint()))
	}

	return builder.String()
}
