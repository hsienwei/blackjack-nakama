package main

import "fmt"

var cardPoint = [CARD_RANK_LEN]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}

type Option string

const (
	BET    Option = "bet"
	HIT    Option = "hit"
	STAND  Option = "stand"
	SPLIT  Option = "split"
	DOUBLE Option = "double"
	RESULT Option = "result"
)

func IsAllowSplit(card1, card2 Card) bool {
	return card1.Point() == card2.Point()
}

type Point struct {
	Soft int
	Hard int
}

func GetPoint(cards []Card) Point {
	rankAceCount := 0
	point := 0
	for _, v := range cards {
		cardRank := v.Rank()
		if int(cardRank) >= CARD_RANK_LEN {
			continue
		}
		if int(v.Suit()) >= CARD_SUIT_LEN {
			continue
		}

		if cardRank == CARD_RANK_A {
			rankAceCount++
		}
		point += cardPoint[cardRank]
	}

	if rankAceCount > 0 && point <= 11 {
		return Point{point + 10, point}
	} else {
		return Point{point, point}
	}
}

func CompareAndPay(playerCards []Card, bet uint32, bankerCards []Card) uint32 {
	fmt.Printf("PlayerCard:%v   bankerCards:%v  bet: %d \n", playerCards, bankerCards, bet)
	playerPoint := GetPoint(playerCards)
	bankerPoint := GetPoint(bankerCards)
	var rtn uint32 = 0
	if playerPoint.Soft > 21 {
		fmt.Println("Player hands bust , Banker Win")
		//rtn = -int32(bet)
	} else if bankerPoint.Soft > 21 {
		fmt.Println("Banker hands bust , Player Win")
		rtn = bet * 2
	} else if playerPoint.Soft > bankerPoint.Soft {
		fmt.Println("Player hands > Banker hands , Player Win")
		rtn = bet * 2
	} else if playerPoint.Soft < bankerPoint.Soft {
		fmt.Println("Banker hands > Player hands , Banker Win")
		//rtn = -int32(bet)
	} else {
		fmt.Println("Banker hands = Player hands , no win")
	}

	return rtn
}

func getActionOption(game *Player) []Option {
	rtn := []Option{}

	if game.IsAllHandsFinished() && game.Credit >= 50 {
		rtn = append(rtn, BET)
	} else {
		curHand := game.getCurrentHand()
		if curHand.Bet == 0 {
			rtn = append(rtn, BET)
		} else {
			rtn = append(rtn, HIT, STAND)
		}
		point := GetPoint(curHand.Cards)
		if len(curHand.Cards) == 2 && game.Credit >= curHand.Bet && !curHand.Doubled && point.Soft < 21 {
			rtn = append(rtn, DOUBLE)
		}
		if len(curHand.Cards) == 2 && game.Credit >= curHand.Bet && IsAllowSplit(curHand.Cards[0], curHand.Cards[1]) {
			rtn = append(rtn, SPLIT)
		}
	}

	return rtn
}
