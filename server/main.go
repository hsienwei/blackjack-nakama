package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"

	"github.com/heroiclabs/nakama-common/runtime"
)

// const (
// 	BET uint8 = iota
// 	ACTION
// )

const (
	BET uint8 = iota
	HIT
	STAND
	SPLIT
	DOUBLE
)

type BetSet struct {
	Bet   uint32
	Cards []int
}

type Banker struct {
	Cards        []uint8
	ShuffleCards []uint8
	DealIndex    uint8
}

var defaultDeckOfCards = [52]uint8{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C,
}

func (banker *Banker) setDeckOfCount(count int) {

	if count <= 0 {
		banker.setDeckOfCount(1)
		return
	}

	banker.ShuffleCards = make([]uint8, count*52)

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

func (banker *Banker) Deal() uint8 {
	rtn := banker.ShuffleCards[banker.DealIndex]
	banker.DealIndex++

	return rtn
}

type PlayerGame struct {
	Credit uint32
	Status uint8

	Bets     []BetSet
	BetIndex uint8
}

type PlayerActionRequest struct {
	Action uint8
	Value  any
}

// func (s *BetSet) MarshalJSON() ([]byte, error) {
// 	var array string
// 	if s.Cards == nil {
// 		array = "[]"
// 	} else {
// 		array = strings.Join(strings.Fields(fmt.Sprintf("%d", s.Cards)), ",")
// 	}
// 	jsonResult := fmt.Sprintf(`{"Bet":%q,"Cards":%s}`, s.Bet, array)
// 	return []byte(jsonResult), nil
// }

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

func Join(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	logger.Info("userId %s %s", userId, ok)

	banker = Banker{
		Cards:        []uint8{},
		ShuffleCards: nil,
	}

	banker.ShuffleCard()

	game = PlayerGame{
		Credit:   10000,
		Status:   BET,
		BetIndex: 0,
		Bets: []BetSet{
			{0, []int{}},
		},
	}

	if err := json.Unmarshal([]byte(payload), &game); err != nil {
		return "", runtime.NewError("unable to unmarshal payload", 13)
	}

	response, err := json.Marshal(game)
	if err != nil {
		return "", runtime.NewError("unable to marshal payload", 13)
	}

	return string(response), nil
}

func getActionOption(game *PlayerGame) []uint8 {
	return []uint8{HIT, STAND}
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

	if action.Action == BET {
		if game.Bets[game.BetIndex].Bet == 0 {
			betAmount := action.Value.(uint32)
			game.Bets[game.BetIndex].Bet = betAmount
			game.Credit -= betAmount

			game.Bets[game.BetIndex].Cards = append(game.Bets[game.BetIndex].Cards, int(banker.Deal()), int(banker.Deal()))
		}
	} else if action.Action == HIT {
		game.Bets[game.BetIndex].Cards = append(game.Bets[game.BetIndex].Cards, int(banker.Deal()))
	} else if action.Action == STAND {
		game.BetIndex = 0
		game.Bets[game.BetIndex].Bet = 0
		game.Bets[game.BetIndex].Cards = []int{}
	}

	getActionOption(&game)

	if err := json.Unmarshal([]byte(payload), &game); err != nil {
		return "", runtime.NewError("unable to unmarshal payload", 13)
	}

	response, err := json.Marshal(game)
	if err != nil {
		return "", runtime.NewError("unable to marshal payload", 13)
	}

	return string(response), nil
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Hello World!")

	// players = map[string]PlayerGame{}

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
