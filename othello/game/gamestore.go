package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/buger/jsonparser"
	websocket "github.com/kataras/go-websocket"
	"github.com/olebedev/emitter"
	"github.com/sdbx/othello-server/othello/models"
)

type (
	GameStore struct {
		WS        websocket.Server
		userStore models.UserStore
		games     map[string]*gameRoom
	}
	gameClient struct {
		user *models.User
		room *gameRoom
	}
	gameRoom struct {
		name  string
		game  *Game
		store *GameStore
		ws    websocket.Server

		clients    map[*gameClient]bool
		register   chan *gameClient
		unregister chan *gameClient
		ticker1    *time.Ticker
		ticker10   *time.Ticker
	}
)

func (g *gameRoom) emitMessage(message []byte) {
	log.Println("websocket sent from", g.name, ":", string(message))
	for _, con := range g.ws.GetConnectionsByRoom(g.name) {
		con.EmitMessage(message)
	}
}

func (g *gameRoom) emit(typ string, ho h) {
	ho["type"] = typ
	content, _ := json.Marshal(ho)
	g.emitMessage(content)
}

func (g *gameRoom) run() {
	for {
		select {
		case client := <-g.register:
			g.clients[client] = true
		case client := <-g.unregister:
			delete(g.clients, client)
		case <-g.ticker1.C:
			if g.game.Turn() == GameTurnBlack {
				g.game.BlackTime--
				if g.game.BlackTime == 0 {
					g.emit("end", h{
						"winner": "white",
						"cause":  "timeout",
					})
					<-g.game.Emitter.Emit("end")
				}
			} else {
				g.game.WhiteTime--
				if g.game.WhiteTime == 0 {
					g.emit("end", h{
						"winner": "black",
						"cause":  "timeout",
					})
					<-g.game.Emitter.Emit("end")
				}
			}
		case <-g.ticker10.C:
			g.emit("tick", h{
				"black": g.game.BlackTime,
				"white": g.game.WhiteTime,
			})
		}
	}
}

func NewGameStore(userStore models.UserStore) *GameStore {
	gs := &GameStore{
		WS:        websocket.New(websocket.Config{}),
		userStore: userStore,
		games:     make(map[string]*gameRoom),
	}
	gs.WS.OnConnection(gs.handleConnection)
	return gs
}

func (gs *GameStore) CreateGame(room string, black string, white string, gameType GameType) error {
	log.Println(room, "game created")
	if _, ok := gs.games[room]; ok {
		return errors.New("game already exist")
	}
	gameroom := &gameRoom{
		name:       room,
		store:      gs,
		ws:         gs.WS,
		clients:    make(map[*gameClient]bool),
		register:   make(chan *gameClient),
		unregister: make(chan *gameClient),
		ticker1:    time.NewTicker(time.Second),
		ticker10:   time.NewTicker(time.Second * 10),
	}
	gameroom.game = newGame(gameroom, black, white, gameType)
	gameroom.game.Emitter.On("end", func(e *emitter.Event) {
		go func() {
			time.Sleep(time.Second)
			for _, conn := range gs.WS.GetConnectionsByRoom(room) {
				conn.Disconnect()
			}
			delete(gs.games, room)
		}()
	})
	gs.games[room] = gameroom
	go gameroom.run()
	return nil
}

func (gs *GameStore) GetGame(room string) *Game {
	if groom, ok := gs.games[room]; ok {
		return groom.game
	}
	return nil
}

type loginRequest struct {
	Type   string `json:"type"`
	Secret string `json:"secret"`
	Room   string `json:"game"`
}

func (gs *GameStore) handleConnection(c websocket.Connection) {
	client := &gameClient{}
	c.OnMessage(func(message []byte) {
		log.Println("websocket recieved:", string(message))
		typ, err := jsonparser.GetString(message, "type")
		if err != nil {
			c.EmitMessage([]byte(jsonErrorMsg))
			return
		}
		switch typ {
		case "ping":
			c.EmitMessage([]byte(pongMsg))
		case "enter":
			if client.user != nil {
				c.EmitMessage([]byte(onceMsg))
				return
			}
			req := loginRequest{}
			err = json.Unmarshal(message, &req)
			if err != nil {
				c.EmitMessage([]byte(jsonErrorMsg))
				return
			}
			user := gs.userStore.GetUserBySecret(req.Secret)
			if user == nil {
				c.EmitMessage([]byte(fmt.Sprintf(userNoMsg, "login")))
				return
			}
			room, ok := gs.games[req.Room]
			if !ok {
				c.EmitMessage([]byte(fmt.Sprintf(roomNoMsg, "login")))
				return
			}
			client.user = user
			client.room = room
			c.Join(room.name)
			room.emitMessage([]byte(fmt.Sprintf(connectMsg, client.user.Name)))
			room.register <- client
		default:
			c.EmitMessage([]byte(typeErrorMsg))
		}
	})

	c.OnDisconnect(func() {
		if client.room != nil {
			client.room.emitMessage([]byte(fmt.Sprintf(disconnectMsg, client.user.Name)))
			client.room.unregister <- client
		}
	})
}
