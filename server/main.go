package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Card uint8
type Option uint8

const (
	BET Option = iota
	HIT
	STAND
	SPLIT
	DOUBLE
	RESULT
)

type BetSet struct {
	Bet   uint32
	Cards []Card
}

type Banker struct {
	Cards        []Card
	ShuffleCards []Card
	DealIndex    uint8
}

var defaultDeckOfCards = [52]Card{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C,
}

const HideCard Card = 99

func (c Card) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", c)), nil
}
func (c Option) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", c)), nil
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

type PlayerGame struct {
	Credit uint32
	// Status Option

	Bets     []BetSet
	BetIndex uint8
}

type PlayerActionRequest struct {
	Action Option
	Value  any
}

type Response struct {
	Player      PlayerGame
	Option      []Option
	BankerCards []Card
}

// var players = make(map[string]PlayerGame)

var game PlayerGame
var banker Banker

// func getPlayerGame(logger runtime.Logger, userID string) *PlayerGame {

// 	value, isExist := players[userID]

// 	logger.Info("%s %s", value, isExist)

// 	if !isExist {
// 		players[userID] = PlayerGame{Status: 0}
// 		value = players[userID]
// 	}
// 	logger.Info("%s", players)
// 	return &value
// }

func getResponse(obj Response) (string, error) {

	response, err := json.Marshal(obj)
	if err != nil {
		return "", runtime.NewError("unable to marshal payload", 13)
	}

	return string(response), nil
}

func getActionOption(game *PlayerGame) []Option {
	return []Option{HIT, STAND, DOUBLE, SPLIT}
}

func getRespinseBankerCards(action Option) []Card {

	if action != RESULT {
		if len(banker.Cards) < 2 {
			return banker.Cards
		} else {
			rtnCards := banker.Cards
			rtnCards[1] = HideCard
			return rtnCards
		}
	}

	return banker.Cards

}

func Join(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)

	banker = Banker{
		Cards:        []Card{},
		ShuffleCards: nil,
	}

	banker.ShuffleCard()

	game = PlayerGame{
		Credit: 10000,
		// Status:   BET,
		BetIndex: 0,
		Bets: []BetSet{
			{0, []Card{}},
		},
	}

	return getResponse(Response{
		Player:      game,
		Option:      getActionOption(&game),
		BankerCards: getRespinseBankerCards(BET),
	})
}

func Action(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)
	logger.Info("payload %s", payload)

	action := PlayerActionRequest{}

	if err := json.Unmarshal([]byte(payload), &action); err != nil {
		return "", runtime.NewError("unable to unmarshal payload", 13)
	}
	logger.Info("payload %s", action)

	banker.CheckReshuffleCard()

	switch action.Action {
	case BET:
		if game.Bets[game.BetIndex].Bet == 0 {
			betAmount := uint32(action.Value.(float64))
			game.Bets[game.BetIndex].Bet = betAmount
			game.Credit -= betAmount

			game.Bets[game.BetIndex].Cards = append(game.Bets[game.BetIndex].Cards, banker.Deal(), banker.Deal())

			banker.Cards = append(banker.Cards, banker.Deal(), banker.Deal())
		}
	case HIT:
		game.BetIndex = 0
		game.Bets[game.BetIndex].Cards = append(game.Bets[game.BetIndex].Cards, banker.Deal())
	case STAND:
		// game.BetIndex = 0
		// game.Bets[game.BetIndex].Bet = 0
		// game.Bets[game.BetIndex].Cards = []Card{}
	case SPLIT:
		// game.BetIndex = 0
		// game.Bets[game.BetIndex].Bet = 0
		// game.Bets[game.BetIndex].Cards = []Card{}
		game.Bets = append(game.Bets, BetSet{Bet: game.Bets[game.BetIndex].Bet, Cards: []Card{game.Bets[game.BetIndex].Cards[1]}})
		game.Bets[game.BetIndex].Cards = game.Bets[game.BetIndex].Cards[:1]
	case DOUBLE:
		// game.BetIndex = 0
		// game.Bets[game.BetIndex].Bet = 0
		game.Credit -= game.Bets[game.BetIndex].Bet
		game.Bets[game.BetIndex].Bet *= 2
		// game.Bets[game.BetIndex].Cards = []Card{}
	}

	return getResponse(Response{
		Player:      game,
		Option:      getActionOption(&game),
		BankerCards: getRespinseBankerCards(action.Action),
	})
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Hello World!")

	if err := initializer.RegisterRpc("join", Join); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("action", Action); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	return nil
}

// if err := json.Unmarshal([]byte(payload), &game); err != nil {
// 	return "", runtime.NewError("unable to unmarshal payload", 13)
// }

// response, err := json.Marshal(game)
// if err != nil {
// 	return "", runtime.NewError("unable to marshal payload", 13)
// }

// return string(response), nil
