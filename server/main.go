package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Option string

const (
	BET   Option = "bet"
	HIT   Option = "hit"
	STAND Option = "stand"
	//SPLIT  Option = "split"  // Do it later.
	DOUBLE Option = "double" // Do it later.
	RESULT Option = "result"
)

type BetSet struct {
	Bet      uint32
	Cards    []Card
	isDouble bool
}

const DEBUG bool = true

type Player struct {
	Credit uint32

	Bets       []BetSet
	BetIndex   int
	BetOptions []Option
}

func (g *Player) getCurrentBet() *BetSet {
	if g.BetIndex >= len(g.Bets) {
		return nil
	}
	return &(g.Bets[g.BetIndex])
}

func (g *Player) finished() bool {
	return g.BetIndex >= len(g.Bets)
}

func (g *Player) reset() {
	g.Bets = []BetSet{{0, []Card{}, false}}
	g.BetIndex = 0
	g.BetOptions = nil
}

type PlayerActionRequest struct {
	Action Option
	Value  any
}

type Response struct {
	Player      Player
	BankerCards []Card

	BetSetIndex uint8
}

// var players = make(map[string]PlayerGame)

var game *Player
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

func getActionOption(game *Player) []Option {
	rtn := []Option{}

	if game.finished() {
		rtn = append(rtn, BET)
	} else {
		betSet := game.getCurrentBet()
		if betSet.Bet == 0 {
			rtn = append(rtn, BET)
		} else {
			rtn = append(rtn, HIT, STAND)
		}
		if len(betSet.Cards) == 2 {
			rtn = append(rtn, DOUBLE)
		}
	}

	return rtn
}

func getResponseBankerCards(finished bool) []Card {

	if !finished {
		if len(banker.Cards) < 2 {
			return banker.Cards
		} else {
			rtnCards := make([]Card, len(banker.Cards))

			copy(rtnCards, banker.Cards)
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

	game = new(Player)
	game.Credit = 10000
	game.Bets = append(game.Bets, BetSet{0, []Card{}, false})
	game.BetOptions = getActionOption(game)

	return getResponse(Response{
		Player:      *game,
		BankerCards: getResponseBankerCards(false),
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

	isActionVaild := false
	for _, v := range game.BetOptions {
		if v == action.Action {
			isActionVaild = true
		}
	}

	if !isActionVaild {
		return "", runtime.NewError("action not allow", 100)
	}

	currentBetSet := game.getCurrentBet()
	switch action.Action {
	case BET:
		banker.CheckReshuffleCard()
		banker.ClearCards()
		game.reset()
		currentBetSet = game.getCurrentBet()
		if currentBetSet.Bet == 0 {
			betAmount := uint32(action.Value.(float64))
			game.Credit -= betAmount

			currentBetSet.Bet = betAmount
			currentBetSet.Cards = append(currentBetSet.Cards, banker.Deal(), banker.Deal())

			banker.Cards = append(banker.Cards, banker.Deal(), banker.Deal())
		}
	case HIT:
		// game.BetIndex = 0
		currentBetSet.Cards = append(currentBetSet.Cards, banker.Deal())
		p := getCardsPoint(currentBetSet.Cards)
		if p[0] >= 21 { //&& p[1] >= 21 {
			game.BetIndex++
		}
	case STAND:
		game.BetIndex++
	// game.BetIndex = 0
	// currentBetSet.Bet = 0
	//currentBetSet.Cards = []Card{}

	// case SPLIT:
	// 	// game.BetIndex = 0
	// 	// currentBetSet.Bet = 0
	// 	// currentBetSet.Cards = []Card{}
	// 	game.Bets = append(game.Bets, BetSet{Bet: currentBetSet.Bet, Cards: []Card{currentBetSet.Cards[1]}})
	// 	currentBetSet.Cards = currentBetSet.Cards[:1]
	case DOUBLE:
		// 	// game.BetIndex = 0
		currentBetSet.isDouble = true
		game.Credit -= currentBetSet.Bet
		currentBetSet.Bet *= 2
		currentBetSet.Cards = append(currentBetSet.Cards, banker.Deal())
		game.BetIndex++
	}

	if game.finished() {
		banker.DrawCards()
	}

	game.BetOptions = getActionOption(game)

	return getResponse(Response{
		Player:      *game,
		BankerCards: getResponseBankerCards(game.finished()),
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
