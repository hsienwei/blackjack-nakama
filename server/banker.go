package main

import "math/rand"

type Banker struct {
	Cards        []Card
	ShuffleCards []Card
	DealIndex    uint8
}

func (banker *Banker) setDeckOfCount(count int) {

	if count <= 0 {
		banker.setDeckOfCount(1)
		return
	}

	banker.ShuffleCards = make([]Card, count*52)

	for i, v := range defaultDeckOfCards {

		for c := 0; c < count; c++ {
			banker.ShuffleCards[i*count+c] = v
		}
	}
}

func (banker *Banker) ShuffleCard() {

	if banker.ShuffleCards == nil {
		banker.setDeckOfCount(1)
	}

	r := rand.New(rand.NewSource(0))
	r.Shuffle(
		len(banker.ShuffleCards),
		func(i, j int) {
			(banker.ShuffleCards)[i], (banker.ShuffleCards)[j] = (banker.ShuffleCards)[j], (banker.ShuffleCards)[i]
		})
	banker.DealIndex = 0
}

func (banker *Banker) CheckReshuffleCard() bool {
	isShuffle := false
	if len(banker.ShuffleCards)-int(banker.DealIndex) < 15 {
		banker.ShuffleCard()
		isShuffle = true
	}

	return isShuffle
}

func (banker *Banker) Deal() Card {
	rtn := banker.ShuffleCards[banker.DealIndex]
	banker.DealIndex++

	return rtn
}
