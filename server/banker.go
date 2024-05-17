package main

import (
	"fmt"
	"strings"
)

type Banker struct {
	Cards []Card
}

func (banker *Banker) reset() {
	banker.Cards = []Card{}
}

func (banker *Banker) displayHand(finished bool) []Card {

	if !finished {
		if len(banker.Cards) < 2 {
			return banker.Cards
		} else {
			rtnCards := make([]Card, len(banker.Cards))

			copy(rtnCards, banker.Cards)
			rtnCards[1] = HIDE_CARD
			return rtnCards
		}
	}
	return banker.Cards

}

func (banker *Banker) DrawCards(dealer *Dealer) {
	p := GetPoint(banker.Cards)
	for p.Soft < 17 {
		banker.Cards = dealer.DealTo(banker.Cards, 1)
		p = GetPoint(banker.Cards)
	}

}

type Result struct {
	BankerPoint Point
	HandResult  []HandResult
}

func (banker *Banker) getResult(player *Player) *Result {
	result := new(Result)
	result.BankerPoint = GetPoint(banker.Cards)
	result.HandResult = make([]HandResult, len(player.Hands))
	for i := 0; i < len(player.Hands); i++ {
		result.HandResult[i] = CompareAndPay(player.Hands[i].Cards, player.Hands[i].Bet, banker.Cards)
	}

	return result
}

func (b *Banker) String() string {
	builder := strings.Builder{}
	builder.WriteString(">> BANKER \n")
	builder.WriteString(fmt.Sprintf("%v %v\n ", b.Cards, GetPoint(b.Cards)))

	return builder.String()
}
