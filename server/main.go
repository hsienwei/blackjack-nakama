package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

const DEBUG bool = true

type RequestPlayerAction struct {
	Action Option
	Value  any
}

type Response struct {
	Player      Player
	BankerCards []Card
	Result      Result
}

// var players = make(map[string]PlayerGame)

var player *Player
var banker *Banker

var dealer *Dealer

func getMarshalString(obj any) (string, error) {

	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", runtime.NewError("unable to marshal payload", 13)
	}

	return string(bytes), nil
}

func Join(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)

	banker = new(Banker)
	dealer = new(Dealer)
	player = new(Player)
	player.Credit = 10000
	player.Options = getActionOption(player)

	return getMarshalString(Response{
		Player:      *player,
		BankerCards: banker.displayHand(false),
		//Result:      *result,
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

	if player.actionIllegal(action.Action) {
		return "", runtime.NewError("Action Illegal", 100)
	}

	switch action.Action {
	case BET:
		dealer.CheckReshuffleCard()
		banker.reset()
		player.reset()

		betAmount := uint32(action.Value.(float64))
		player.Bet(betAmount)

		curHand := player.getCurrentHand()
		curHand.Cards = dealer.DealTo(curHand.Cards, 1)
		banker.Cards = dealer.DealTo(banker.Cards, 1)
		curHand.Cards = dealer.DealTo(curHand.Cards, 1)
	case HIT:
		player.Hit()
	case STAND:
		player.Stand()
	case SPLIT:
		player.Split()
	case DOUBLE:
		player.Split()
	}

	result := new(Result)

	if player.IsAllHandsFinished() {
		banker.DrawCards(dealer)
		result = banker.getResult(player)
	}

	player.Options = getActionOption(player)

	return getMarshalString(Response{
		Player:      *player,
		BankerCards: banker.displayHand(player.IsAllHandsFinished()),
		Result:      *result,
	})
}

var x int = 0

func Test(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("test %s %v", payload, x)
	x++
	logger.Info("%v", ctx)
	time.Sleep(10 * time.Second)
	logger.Info("sleep end")

	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if ok {
		logger.Info(userId)
	} else {
		logger.Info("no user id")
	}

	return "", nil
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

	if err := initializer.RegisterRpc("test", Test); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	return nil
}
