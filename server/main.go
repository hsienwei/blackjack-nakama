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

const DEBUG bool = true

type RequestPlayerAction struct {
	Action Option
	Value  any
}

type Response struct {
	Player      Player
	BankerCards []Card

	BetSetIndex uint8
}

// var players = make(map[string]PlayerGame)

var player *Player
var banker *Banker

var dealer *Dealer

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

func getMarshalString(obj any) (string, error) {

	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", runtime.NewError("unable to marshal payload", 13)
	}

	return string(bytes), nil
}

func getActionOption(game *Player) []Option {
	rtn := []Option{}

	if game.finished() {
		rtn = append(rtn, BET)
	} else {
		betSet := game.getCurrentHand()
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

func Join(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)

	banker = new(Banker)
	dealer = new(Dealer)
	player = new(Player)
	player.Credit = 10000
	player.Hands = append(player.Hands, Hand{0, []Card{}, false})
	player.Options = getActionOption(player)

	return getMarshalString(Response{
		Player:      *player,
		BankerCards: banker.displayHand(false),
	})
}

func Action(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)
	logger.Info("payload %s", payload)

	action := new(RequestPlayerAction)

	if err := json.Unmarshal([]byte(payload), action); err != nil {
		return "", runtime.NewError("unable to unmarshal payload", 13)
	}
	logger.Info("payload %v", action)

	isActionVaild := false
	for _, v := range player.Options {
		if v == action.Action {
			isActionVaild = true
		}
	}

	if !isActionVaild {
		return "", runtime.NewError("action not allow", 100)
	}

	currentBetSet := player.getCurrentHand()
	switch action.Action {
	case BET:
		dealer.CheckReshuffleCard()
		banker.reset()
		player.reset()
		currentBetSet = player.getCurrentHand()
		if currentBetSet.Bet == 0 {
			betAmount := uint32(action.Value.(float64))
			player.Credit -= betAmount

			currentBetSet.Bet = betAmount
			currentBetSet.Cards = append(currentBetSet.Cards, dealer.Deal(), dealer.Deal())

			banker.Cards = append(banker.Cards, dealer.Deal(), dealer.Deal())
		}
	case HIT:
		// game.BetIndex = 0
		currentBetSet.Cards = append(currentBetSet.Cards, dealer.Deal())
		p := getCardsPoint(currentBetSet.Cards)
		if p[0] >= 21 { //&& p[1] >= 21 {
			player.HandIndex++
		}
	case STAND:
		player.HandIndex++
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
		currentBetSet.isDouble = true
		player.Credit -= currentBetSet.Bet
		currentBetSet.Bet *= 2
		currentBetSet.Cards = append(currentBetSet.Cards, dealer.Deal())
		player.HandIndex++
	}

	if player.finished() {
		banker.DrawCards(dealer)
	}

	player.Options = getActionOption(player)

	return getMarshalString(Response{
		Player:      *player,
		BankerCards: banker.displayHand(player.finished()),
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
