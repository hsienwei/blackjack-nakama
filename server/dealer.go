package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Dealer struct {
	ShuffleCards []Card
	DealIndex    int
}

func (dealer *Dealer) setDeckOfCount(count int) {

	if count <= 0 {
		dealer.setDeckOfCount(1)
		return
	}

	dealer.ShuffleCards = make([]Card, count*52)

	for i, v := range defaultDeckOfCards {

		for c := 0; c < count; c++ {
			dealer.ShuffleCards[i*count+c] = v
		}
	}
}

func (dealer *Dealer) ShuffleCard() {

	if dealer.ShuffleCards == nil {
		dealer.setDeckOfCount(1)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(
		len(dealer.ShuffleCards),
		func(i, j int) {
			(dealer.ShuffleCards)[i], (dealer.ShuffleCards)[j] = (dealer.ShuffleCards)[j], (dealer.ShuffleCards)[i]
		})
	dealer.DealIndex = 0
}

func (dealer *Dealer) CheckReshuffleCard() bool {
	// if len(banker.ShuffleCards)-int(banker.DealIndex) < 15 {
	// 	banker.ShuffleCard()
	// 	return true
	// }

	// return false

	dealer.ShuffleCard()
	return true
}

func (dealer *Dealer) Deal() Card {
	card := dealer.ShuffleCards[dealer.DealIndex]
	dealer.DealIndex++

	return card
}

func (dealer *Dealer) DealTo(target []Card, count int) []Card {
	c := make([]Card, count)
	for i := 0; i < count; i += 1 {
		c[i] = dealer.Deal()
	}

	return append(target, c...)
}

func (dealer *Dealer) String() string {
	return fmt.Sprintf("Less %d Cards.", len(dealer.ShuffleCards)-dealer.DealIndex)
}
