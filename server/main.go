package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/heroiclabs/nakama-common/runtime"
)

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

const DEBUG bool = true

func getOptionString(c Option) string {
	o := "undefined"
	switch c {
	case BET:
		o = "BET"
	case HIT:
		o = "HIT"
	default:
		o = "UNKNOWN"
	}

	return o
}

func (c Option) MarshalJSON() ([]byte, error) {
	if DEBUG {
		return []byte(fmt.Sprintf("%q", getOptionString(c))), nil
	}
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

func (g *PlayerGame) getCurrentBet() *BetSet {
	return &(g.Bets[g.BetIndex])
}

type PlayerActionRequest struct {
	Action Option
	Value  any
}

type Response struct {
	Player      PlayerGame
	BankerCards []Card

	BetSetIndex uint8
	Option      []Option
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
			rtnCards[1] = HIDE_CARD
			return rtnCards
		}
	}

	return banker.Cards

}

func Join(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)

	// banker = Banker{
	// 	Cards:        []Card{},
	// 	ShuffleCards: nil,
	// }

	banker := new(Banker)

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

	currentBetSet := game.getCurrentBet()
	switch action.Action {
	case BET:
		if currentBetSet.Bet == 0 {
			betAmount := uint32(action.Value.(float64))
			game.Credit -= betAmount

			currentBetSet.Bet = betAmount
			currentBetSet.Cards = append(currentBetSet.Cards, banker.Deal(), banker.Deal())

			banker.Cards = append(banker.Cards, banker.Deal(), banker.Deal())
		}
	case HIT:
		game.BetIndex = 0
		currentBetSet.Cards = append(currentBetSet.Cards, banker.Deal())
	case STAND:
		// game.BetIndex = 0
		// currentBetSet.Bet = 0
		// currentBetSet.Cards = []Card{}
	case SPLIT:
		// game.BetIndex = 0
		// currentBetSet.Bet = 0
		// currentBetSet.Cards = []Card{}
		game.Bets = append(game.Bets, BetSet{Bet: currentBetSet.Bet, Cards: []Card{currentBetSet.Cards[1]}})
		currentBetSet.Cards = currentBetSet.Cards[:1]
	case DOUBLE:
		// game.BetIndex = 0
		// currentBetSet.Bet = 0
		game.Credit -= currentBetSet.Bet
		currentBetSet.Bet *= 2
		// currentBetSet.Cards = []Card{}
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
