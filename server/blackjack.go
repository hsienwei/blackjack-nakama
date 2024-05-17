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

func IsAllowSplit(cards []Card) bool {
	if len(cards) != 2 {
		return false
	}

	return cards[0].Point() == cards[1].Point()
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

type HandResult struct {
	WinAmount uint32
	Comment   string
}

func CompareAndPay(playerCards []Card, bet uint32, bankerCards []Card) HandResult {
	fmt.Printf("PlayerCard:%v   bankerCards:%v  bet: %d \n", playerCards, bankerCards, bet)
	playerPoint := GetPoint(playerCards)
	bankerPoint := GetPoint(bankerCards)
	var rtn HandResult
	if playerPoint.Soft > 21 {
		fmt.Println("Player hands bust , Banker Win")
		rtn = HandResult{0, "Player burst, Banker Win"}
	} else if bankerPoint.Soft > 21 {
		fmt.Println("Banker hands bust , Player Win")
		rtn = HandResult{bet * 2, "Banker burst, Player Win"}
	} else if playerPoint.Soft > bankerPoint.Soft {
		fmt.Println("Player hands > Banker hands , Player Win")
		rtn = HandResult{bet * 2, "Player bigger, Player Win"}
	} else if playerPoint.Soft < bankerPoint.Soft {
		fmt.Println("Banker hands > Player hands , Banker Win")
		rtn = HandResult{0, "Banker bigger, Banker Win"}
	} else {
		fmt.Println("Banker hands = Player hands , no win")
		rtn = HandResult{bet, "Banker = Player"}
	}

	return rtn
}

func getActionOption(p *Player) []Option {
	rtn := []Option{}

	if p.IsAllHandsFinished() && p.Credit >= 50 {
		rtn = append(rtn, BET)
	} else {
		curHand := p.CurrentHand()
		point := GetPoint(curHand.Cards)

		if curHand.Bet == 0 {
			rtn = append(rtn, BET)
		} else {
			if point.Soft < 21 {
				rtn = append(rtn, HIT, STAND)
			} else {
				rtn = append(rtn, STAND)
			}

			if len(curHand.Cards) == 2 && p.Credit >= curHand.Bet {
				if !(point.Soft < 21) {
					rtn = append(rtn, DOUBLE)
				}
				if IsAllowSplit(curHand.Cards) {
					rtn = append(rtn, SPLIT)
				}
			}
		}

	}

	return rtn
}
