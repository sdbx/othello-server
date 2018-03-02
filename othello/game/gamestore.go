package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/buger/jsonparser"
	websocket "github.com/kataras/go-websocket"
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
		name string
		game *Game
		ws   websocket.Server

		clients    map[*gameClient]bool
		register   chan *gameClient
		unregister chan *gameClient
		close      chan bool
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
		case <-g.close:
			return
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
	log.Println(room, " game created")
	if _, ok := gs.games[room]; ok {
		return errors.New("game already exist")
	}
	gameroom := &gameRoom{
		name:       room,
		ws:         gs.WS,
		clients:    make(map[*gameClient]bool),
		register:   make(chan *gameClient),
		unregister: make(chan *gameClient),
		close:      make(chan bool),
	}
	gameroom.game = newGame(gameroom, black, white, gameType)
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
		case "login":
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
