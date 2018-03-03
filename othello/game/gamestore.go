package game

import (
	"errors"
	"log"
	"time"

	"github.com/sdbx/othello-server/othello/models"
	"github.com/sdbx/othello-server/othello/ws"
)

type (
	h         map[string]interface{}
	GameStore struct {
		*ws.WSStore
	}
	gameClient struct {
		user *models.User
		room *gameRoom
	}
	gameRoom struct {
		*ws.WSRoom
		game     *Game
		ticker1  *time.Ticker
		ticker10 *time.Ticker
	}
)

func NewGameStore(userStore models.UserStore) *GameStore {
	gs := &GameStore{
		WSStore: ws.NewWSStore(userStore),
	}
	return gs
}

func (gs *GameStore) CreateGame(room string, black string, white string, gameType GameType) (*Game, error) {
	log.Println(room, "game created")
	if _, ok := gs.Rooms[room]; ok {
		return nil, errors.New("game already exist")
	}
	gameroom := &gameRoom{
		WSRoom:   ws.NewWSRoom(room, gs.WSStore),
		ticker1:  time.NewTicker(time.Second),
		ticker10: time.NewTicker(time.Second * 10),
	}
	gam := newGame(gameroom, black, white, gameType)
	gameroom.game = gam
	gs.Rooms[room] = gameroom
	go gameroom.Run()
	go gameroom.runGame()
	return gam, nil
}

func (g *gameRoom) runGame() {
	end := g.game.Emitter.On("end")
	turn := g.game.Emitter.On("turn")
	for {
		select {
		case event := <-end:
			g.Emit("end", event.Args[0].(ws.H))
			g.Close()
		case event := <-turn:
			g.Emit("turn", event.Args[0].(ws.H))
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
	if groom, ok := gs.Rooms[room]; ok {
		return groom.(*gameRoom).game
	}
	return nil
}
