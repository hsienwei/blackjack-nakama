package main

import (
	"fmt"
	"strings"
)

type Banker struct {
	CardSet
}

func (banker *Banker) displayHand(finished bool) []Card {

	if !finished {
		if len(banker.CardSet) < 2 {
			return banker.CardSet
		} else {
			rtnCards := make([]Card, len(banker.CardSet))

			copy(rtnCards, banker.CardSet)
			rtnCards[1] = HIDE_CARD
			return rtnCards
		}
	}
	return banker.CardSet

}

func (banker *Banker) DrawCards(dealer *Dealer) {
	p := banker.CardSet.GetPoint()
	for p.Soft < 17 {
		banker.CardSet = dealer.DealTo(banker.CardSet, 1)
		p = banker.CardSet.GetPoint()
	}

}

type Result struct {
	BankerPoint Point
	HandResult  []HandResult
}

func (banker *Banker) getResult(player *Player) *Result {
	result := new(Result)
	result.BankerPoint = banker.CardSet.GetPoint()
	result.HandResult = make([]HandResult, len(player.Hands))
	for i := 0; i < len(player.Hands); i++ {
		result.HandResult[i] = CompareAndPay(player.Hands[i].CardSet, player.Hands[i].Bet, banker.CardSet)
	}

	return result
}

func (b *Banker) String() string {
	builder := strings.Builder{}
	builder.WriteString(">> BANKER \n")
	builder.WriteString(fmt.Sprintf("%v %v\n ", b.CardSet, b.CardSet.GetPoint()))

	return builder.String()
}
