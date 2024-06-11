package main

import (
	"fmt"
	"slices"
	"testing"
)

func initial() {
	player = new(Player)
	player.Credit = 10000
	player.Hands = append(player.Hands, Hand{})
	player.Options = getActionOption(player)

	banker = new(Banker)
	banker.reset()

	dealer = new(Dealer)
	dealer.ShuffleCard()
}

func _CheckPointEqual(t *testing.T, cards []Card, targetPoint int, targetCount int) {
	points := GetPoint(cards)
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
	player = new(Player)
	player.Credit = 10000
	player.Hands = append(player.Hands, Hand{50, []Card{0x0a, 0x1a}})
	fmt.Println(player)
	player.Split()

	fmt.Println(player)
}

func TestGetActionOption(t *testing.T) {
	initial()
	opts := getActionOption(player)
	if !slices.Contains(opts, BET) {
		t.Fail()
	}

	curHand := player.CurrentHand()
	curHand.Bet = 50
	curHand.Cards = []Card{0x16, 0x14, 0x1c}

	opts = getActionOption(player)
	fmt.Printf("%v", opts)
	if slices.Contains(opts, DOUBLE) {
		t.Fail()
	}

	curHand.Cards = []Card{0x10, 0x1a}
	opts = getActionOption(player)
	fmt.Printf("%v", opts)
	if slices.Contains(opts, DOUBLE) {
		t.Fail()
	}

	curHand.Cards = []Card{0x10, 0x20}
	opts = getActionOption(player)
	fmt.Printf("%v", opts)
	if !slices.Contains(opts, SPLIT) {
		t.Fail()
	}

	curHand.Cards = []Card{0x1a, 0x2b}
	opts = getActionOption(player)
	fmt.Printf("%v", opts)
	if !slices.Contains(opts, SPLIT) {
		t.Fail()
	}
}

func TestSimulationGame(t *testing.T) {
	initial()

	for {
		// Bet at start.
		fmt.Println(dealer)
		fmt.Println(player)
		fmt.Println(banker)
		curHand := player.CurrentHand()
		opts := getActionOption(player)
		fmt.Printf("可選選項 %v \n", opts)

		if slices.Contains(opts, BET) {
			fmt.Println("選擇 BET")
			player.Bet(50)
			curHand.Cards = dealer.DealTo(curHand.Cards, 1)
			banker.Cards = dealer.DealTo(banker.Cards, 1)
			curHand.Cards = dealer.DealTo(curHand.Cards, 1)
			banker.Cards = dealer.DealTo(banker.Cards, 1)
		} else if slices.Contains(opts, DOUBLE) {
			fmt.Println("選擇 DOUBLE")
			player.Double()
		} else if slices.Contains(opts, SPLIT) {
			fmt.Println("選擇 SPLIT")
			player.Split()
		} else if slices.Contains(opts, STAND) && GetPoint(curHand.Cards).Soft > 17 {
			fmt.Println("選擇 STAND")
			player.Stand()
		} else if slices.Contains(opts, HIT) {
			fmt.Println("選擇 HIT")
			player.Hit()
		}

		if player.IsAllHandsFinished() {
			break
		}

	}
	fmt.Println("玩家操作結束")
	fmt.Println(player)
	fmt.Println(banker)

	if !player.IsAllHandsBust() {
		fmt.Println("莊家抽牌")
		banker.DrawCards(dealer)
	}

	fmt.Println(banker)
	fmt.Println("比較與賠付")

	var totalPay uint32 = 0
	for i := 0; i < len(player.Hands); i++ {
		totalPay += CompareAndPay(player.Hands[i].Cards, player.Hands[i].Bet, banker.Cards).WinAmount
	}

	fmt.Printf("拿到 %v\n", totalPay)
	player.Credit += totalPay
	fmt.Println(player)

}
