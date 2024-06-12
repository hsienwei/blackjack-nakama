package main

import (
	"fmt"
	"slices"
	"testing"
)

var bj *BlackJackGame

func initial() {
	bj = new(BlackJackGame)
	bj.player.Credit = 10000
	bj.player.Options = bj.GetActionOption()
	bj.dealer.ShuffleCard()
}

func _CheckPointEqual(t *testing.T, cards CardSet, targetPoint int, targetCount int) {
	points := cards.GetPoint()
	fmt.Printf("%v %v\n", cards, points)
	if points.Soft != targetPoint {
		t.Fail()
	}
	if targetCount == 2 && points.Soft == points.Hard {
		t.Fail()
	}
}

func TestGetPoint(t *testing.T) {

	_CheckPointEqual(t, []Card{0x16, 0x14, 0x1c}, 22, 1)
	_CheckPointEqual(t, []Card{0x10, 0x1a}, 21, 2)
	_CheckPointEqual(t, []Card{}, 0, 1)
	_CheckPointEqual(t, []Card{0x0a, 0x0b, 0x0c}, 30, 1)
	_CheckPointEqual(t, []Card{0x0e}, 0, 1)
	_CheckPointEqual(t, []Card{0x40}, 0, 1)
}

func TestSplit(t *testing.T) {
	initial()
	bj.player.Credit = 10000
	bj.player.Hands = append(bj.player.Hands, Hand{50, []Card{0x0a, 0x1a}})
	fmt.Println(bj.player)
	bj.player.Split()

	fmt.Println(bj.player)
}

func TestGetActionOption(t *testing.T) {
	initial()
	opts := bj.GetActionOption()
	if !slices.Contains(opts, BET) {
		t.Fail()
	}

	bj.player.Hands = append(bj.player.Hands, Hand{50, []Card{0x16, 0x14, 0x1c}})
	opts = bj.GetActionOption()
	fmt.Printf("%v", opts)
	if slices.Contains(opts, DOUBLE) {
		t.Fail()
	}

	curHand := bj.player.CurrentHand()
	curHand.CardSet = []Card{0x10, 0x1a}
	opts = bj.GetActionOption()
	fmt.Printf("%v", opts)
	if slices.Contains(opts, DOUBLE) {
		t.Fail()
	}

	curHand.CardSet = []Card{0x10, 0x20}
	opts = bj.GetActionOption()
	fmt.Printf("%v", opts)
	if !slices.Contains(opts, SPLIT) {
		t.Fail()
	}

	curHand.CardSet = []Card{0x1a, 0x2b}
	opts = bj.GetActionOption()
	fmt.Printf("%v", opts)
	if !slices.Contains(opts, SPLIT) {
		t.Fail()
	}
}

func TestSimulationGame(t *testing.T) {
	initial()

	bj.player.Options = bj.GetActionOption()

	for {
		// Bet at start.
		fmt.Println(bj.dealer)
		fmt.Println(bj.player)
		fmt.Println(bj.banker)
		opts := bj.GetActionOption()
		fmt.Printf("可選選項 %v \n", opts)

		if slices.Contains(opts, BET) {
			fmt.Println("選擇 BET")
			bj.ExecAction(BET, 50)
		} else if slices.Contains(opts, DOUBLE) {
			fmt.Println("選擇 DOUBLE")
			bj.ExecAction(DOUBLE, 0)
		} else if slices.Contains(opts, SPLIT) {
			fmt.Println("選擇 SPLIT")
			bj.ExecAction(SPLIT, 0)
		} else if slices.Contains(opts, STAND) &&
			bj.player.CurrentHand().CardSet.GetPoint().Soft > 17 {
			fmt.Println("選擇 STAND")
			bj.ExecAction(STAND, 0)
		} else if slices.Contains(opts, HIT) {
			fmt.Println("選擇 HIT")
			bj.ExecAction(HIT, 0)
		}

		if bj.player.IsAllHandsFinished() {
			break
		}

		bj.player.Options = bj.GetActionOption()
	}
	fmt.Println("玩家操作結束")
	fmt.Println(bj.player)
	fmt.Println(bj.banker)

	if bj.player.IsAllHandsFinished() {
		result := bj.ExecBankerAction()
		fmt.Printf("%#v", result)
	}

	// if !bj.player.IsAllHandsBust() {
	// 	fmt.Println("莊家抽牌")
	// 	banker.DrawCards(dealer)
	// }

	// fmt.Println(banker)
	// fmt.Println("比較與賠付")

	// var totalPay uint32 = 0
	// for i := 0; i < len(player.Hands); i++ {
	// 	totalPay += CompareAndPay(player.Hands[i].Cards, player.Hands[i].Bet, banker.Cards).WinAmount
	// }

	// fmt.Printf("拿到 %v\n", totalPay)
	// player.Credit += totalPay
	// fmt.Println(player)

}
