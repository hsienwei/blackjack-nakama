package main

import (
	"fmt"
)

type Card uint8

const CARD_SUIT_LEN = 4
const CARD_RANK_LEN = 13
const CARD_LEN = CARD_SUIT_LEN * CARD_RANK_LEN

var CARD_SUITS = [CARD_SUIT_LEN]string{"♠", "♥", "♦", "♣"}
var CARD_RANK = [CARD_RANK_LEN]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

const CARD_RANK_A uint8 = 0

const SUIT_MOD = 16
const HIDE_CARD Card = 99

var defaultDeckOfCards = [CARD_LEN]Card{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C,
}

func (c *Card) MarshalJSON() ([]byte, error) {
	if DEBUG {
		return []byte(fmt.Sprintf("%q", c.toString())), nil
	}
	return []byte(fmt.Sprintf("%d", c)), nil
}

func (c *Card) Suit() uint8 {
	return uint8(*c / SUIT_MOD)
}
func (c *Card) SuitStr() string {
	return CARD_SUITS[c.Suit()]
}

func (c *Card) Rank() uint8 {
	return uint8(*c % SUIT_MOD)
}
func (c *Card) RankStr() string {
	return CARD_RANK[c.Rank()]
}

func (c *Card) Point() int {
	return cardPoint[c.Rank()]
}

func (c *Card) toString() string {
	if *c == HIDE_CARD {
		return "X"
	}
	if c.Rank() >= CARD_RANK_LEN {
		return "?"
	}
	if c.Suit() >= CARD_SUIT_LEN {
		return "?"
	}
	return fmt.Sprintf("%s%s", c.SuitStr(), c.RankStr())
}

func (c Card) String() string {
	return c.toString()
}
