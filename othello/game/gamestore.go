package game

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/sdbx/othello-server/othello/dbs"
	"github.com/sdbx/othello-server/othello/ws"
)

type (
	GameStore struct {
		*ws.WSStore
	}
	gameRoom struct {
		*ws.WSRoom
		game     *Game
		ticker1  *time.Ticker
		ticker10 *time.Ticker
	}
)

func NewGameStore() *GameStore {
	gs := &GameStore{
		WSStore: ws.NewWSStore(),
	}
	gs.WSStore.Handlers["enter"] = gs.enterHandler
	return gs
}

func (gs *GameStore) CreateGame(room string, black string, white string, gameType GameType) (*Game, error) {
	gs.Lock()
	defer gs.Unlock()
	if _, ok := gs.Rooms[room]; ok {
		return nil, errors.New("game already exist")
	}
	log.Println(room, "game created")
	gameroom := &gameRoom{
		WSRoom:   ws.NewWSRoom(room, gs.WSStore),
		ticker1:  time.NewTicker(time.Second),
		ticker10: time.NewTicker(time.Second * 10),
	}
	gam := newGame(gameroom, black, white, gameType)
	gameroom.game = gam
	gs.Rooms[room] = gameroom
	go gameroom.runGame()
	return gam, nil
}

func (g *gameRoom) runGame() {
	end := g.game.Emitter.On("end")
	turn := g.game.Emitter.On("turn")
	undo := g.game.Emitter.On("undo")
	undoAns := g.game.Emitter.On("undo_answer")
	for {
		select {
		case event := <-end:
			g.Emit("end", event.Args[0].(ws.H))
			g.Close()
			return
		case event := <-turn:
			g.Emit("turn", event.Args[0].(ws.H))
		case event := <-undo:
			g.Emit("undo", event.Args[0].(ws.H))
		case event := <-undoAns:
			g.Emit("undo_answer", event.Args[0].(ws.H))
		case <-g.ticker1.C:
			g.game.TimeCount()
		case <-g.ticker10.C:
			g.Emit("tick", ws.H{
				"black": g.game.BlackTime,
				"white": g.game.WhiteTime,
			})
		}
	}
}

func (gs *GameStore) GetGame(room string) *Game {
	gs.RLock()
	defer gs.RUnlock()

	if groom, ok := gs.Rooms[room]; ok {
		return groom.(*gameRoom).game
	}
	return nil
}

type enterRequest struct {
	_      string `json:"type"`
	Secret string `json:"secret"`
	Game   string `json:"game"`
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
	user, err := dbs.GetUserBySecret(req.Secret)
	if err != nil {
		cli.EmitError("user doesn't exist", "enter")
		return cli
	}
	cli, err = gs.Enter(cli, user, req.Game)
	if err != nil {
		cli.EmitError(err.Error(), "enter")
	}
	return cli
}
