package game

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/sdbx/othello-server/othello/dbs"
	"github.com/sdbx/othello-server/othello/ws"
)

type GameStore struct {
	*ws.WSStore
}

func NewGameStore() *GameStore {
	gs := &GameStore{
		WSStore: ws.NewWSStore(),
	}
	gs.WSStore.Handlers["enter"] = gs.enterHandler
	return gs
}

func (gs *GameStore) CreateGame(room string, black string, white string, gameType GameType) (*Game, error) {
	gs.WSStore.Lock()
	defer gs.WSStore.Unlock()

	if _, ok := gs.Rooms[room]; ok {
		return nil, errors.New("game already exist")
	}
	log.Println(room, "game created")

	gameroom := &gameRoom{
		WSRoom: ws.NewWSRoom(room, gs.WSStore),
		timer:  time.NewTimer(gameType.Time()),
	}
	gam := newGame(gameroom, black, white, gameType)
	gameroom.game = gam
	gs.Rooms[room] = gameroom
	go gameroom.runGame()

	return gam, nil
}

func (gs *GameStore) GetGame(room string) *Game {
	gs.WSStore.RLock()
	defer gs.WSStore.RUnlock()

	if groom, ok := gs.Rooms[room]; ok {
		return groom.(*gameRoom).game
	}
	return nil
}

type enterRequest struct {
	_    string `json:"type"`
	Game string `json:"game"`
}

const jsonErrorMsg = `{"type":"error","msg":"json error","from":"none"}`
const enterErrorMsg = `{"type":"error","msg":"%s","from":"enter"}`

func (gs *GameStore) enterHandler(cli ws.Client, message []byte) ws.Client {
	req := enterRequest{}

	err := json.Unmarshal(message, &req)
	if err != nil {
		cli.EmitError("json error", "enter")
		return cli
	}

	cli, err = gs.Enter(cli, dbs.User{}, req.Game)
	if err != nil {
		cli.EmitError(err.Error(), "enter")
	}

	return cli
}
