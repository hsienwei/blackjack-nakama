package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Option string

const (
	BET    Option = "bet"
	HIT    Option = "hit"
	STAND  Option = "stand"
	SPLIT  Option = "split"
	DOUBLE Option = "double"
	RESULT Option = "result"
)

type BetSet struct {
	Bet   uint32
	Cards []Card
}

const DEBUG bool = true

type PlayerGame struct {
	Credit uint32
	// Status Option

	Bets       []BetSet
	BetIndex   uint8
	BetOptions []Option
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
}

// var players = make(map[string]PlayerGame)

var game *PlayerGame
var banker *Banker

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

	banker = new(Banker)

	banker.ShuffleCard()

	game = new(PlayerGame)
	game.Credit = 10000
	game.Bets = append(game.Bets, BetSet{0, []Card{}})

	game.BetOptions = getActionOption(game)

	return getResponse(Response{
		Player:      *game,
		BankerCards: getRespinseBankerCards(BET),
	})
}

func Action(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)
	logger.Info("payload %s", payload)

	action := new(PlayerActionRequest)

	if err := json.Unmarshal([]byte(payload), action); err != nil {
		return "", runtime.NewError("unable to unmarshal payload", 13)
	}
	logger.Info("payload %v", action)

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

	game.BetOptions = getActionOption(game)

	return getResponse(Response{
		Player:      *game,
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
