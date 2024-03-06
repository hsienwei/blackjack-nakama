package main

// var cardPoint1 = [13]int{11, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10}
var cardPoint = [13]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10}

func getCardsPoint(cards []Card) []int {
	rankAceCount := 0
	point := 0
	for _, v := range cards {
		cardRank := int(v % SUIT_MOD)
		if cardRank == CARD_RANK_A {
			rankAceCount++
		}
		point += cardPoint[cardRank]

	}

	if rankAceCount > 0 && point <= 11 {
		return []int{point, point + 10}
	} else {
		return []int{point}
	}
}
