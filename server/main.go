package main

import (
	"context"
	"database/sql"
	"encoding/json"

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

// var player *Player
// var banker *Banker

// var dealer *Dealer

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

	bj := new(BlackJackGame)

	bj.player.Credit = 10000
	bj.player.Options = bj.GetActionOption()

	return getMarshalString(Response{
		Player:      bj.player,
		BankerCards: bj.banker.displayHand(false),
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

	bj := new(BlackJackGame)
	bj.ExecAction(action.Action, action.Value)

	result := new(Result)
	result = bj.banker.getResult(&bj.player)

	return getMarshalString(Response{
		Player:      bj.player,
		BankerCards: bj.banker.displayHand(bj.player.IsAllHandsFinished()),
		Result:      *result,
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
