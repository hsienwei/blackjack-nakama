package main

import (
	"fmt"
)

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

type Point struct {
	Soft int
	Hard int
}

type CardSet []Card

func (cards CardSet) IsAllowSplit() bool {
	if len(cards) != 2 {
		return false
	}

	return cards[0].Point() == cards[1].Point()
}

func (cards CardSet) GetPoint() Point {
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

type BlackJackGame struct {
	player Player
	banker Banker
	dealer Dealer
}

func (bj *BlackJackGame) ExecAction(action Option, value any) (bool, string) {
	if bj.player.actionIllegal(action) {
		return false, "Action Illegal"
	}

	switch action {
	case BET:
		bj.dealer.CheckReshuffleCard()
		// bj.banker.reset()
		bj.player.reset()

		betAmount, ok := value.(int)
		if !ok {
			return false, "Value Illegal"
		}

		bj.player.Bet(uint32(betAmount))

		curHand := bj.player.CurrentHand()
		curHand.CardSet = bj.dealer.DealTo(curHand.CardSet, 1)
		bj.banker.CardSet = bj.dealer.DealTo(bj.banker.CardSet, 1)
		curHand.CardSet = bj.dealer.DealTo(curHand.CardSet, 1)
		bj.banker.CardSet = bj.dealer.DealTo(bj.banker.CardSet, 1)

	case HIT:
		bj.player.Hit(&bj.dealer)
	case STAND:
		bj.player.Stand()
	case SPLIT:
		bj.player.Split()
	case DOUBLE:
		bj.player.Double(&bj.dealer)
	}

	// var result *Result

	if bj.player.IsAllHandsFinished() {
		bj.ExecBankerAction()
	}

	bj.player.Options = bj.GetActionOption()

	return true, ""
}

func (bj *BlackJackGame) ExecBankerAction() *Result {
	if !bj.player.IsAllHandsBust() {
		bj.banker.DrawCards(&bj.dealer)
	}

	return bj.banker.getResult(&bj.player)
	// bj.banker.getResult(&bj.player)

	// for i := 0; i < len(result.HandResult); i++ {
	// 	bj.player.Credit += result.HandResult[i].WinAmount
	// }
}

type HandResult struct {
	WinAmount uint32
	Comment   string
}

func CompareAndPay(playerCards CardSet, bet uint32, bankerCards CardSet) HandResult {
	fmt.Printf("PlayerCard:%v   bankerCards:%v  bet: %d \n", playerCards, bankerCards, bet)
	playerPoint := playerCards.GetPoint()
	bankerPoint := bankerCards.GetPoint()
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

func (bj *BlackJackGame) GetActionOption() []Option {
	rtn := []Option{}

	if bj.player.IsAllHandsFinished() && bj.player.Credit >= 50 {
		rtn = append(rtn, BET)
	} else {
		curHand := bj.player.CurrentHand()
		point := curHand.CardSet.GetPoint()

		if curHand.Bet == 0 {
			rtn = append(rtn, BET)
		} else {
			if point.Soft < 21 {
				rtn = append(rtn, HIT, STAND)
			} else {
				rtn = append(rtn, STAND)
			}

			if len(curHand.CardSet) == 2 && bj.player.Credit >= curHand.Bet {
				if point.Soft < 21 {
					rtn = append(rtn, DOUBLE)
				}
				if curHand.CardSet.IsAllowSplit() {
					rtn = append(rtn, SPLIT)
				}
			}
		}

	}

	return rtn
}
