package othello

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/buger/jsonparser"
	websocket "github.com/kataras/go-websocket"
)

type GameStore struct {
	WS        websocket.Server
	userStore *UserStore
	games     map[string]*game
}

func NewGameStore(userStore *UserStore) *GameStore {
	gs := &GameStore{
		WS:        websocket.New(websocket.Config{}),
		userStore: userStore,
		games:     make(map[string]*game),
	}
	gs.WS.OnConnection(gs.handleConnection)
	return gs
}

// for test
func (gs *GameStore) ListGames() []string {
	list := []string{}
	for key := range gs.games {
		list = append(list, key)
	}
	return list
}

func (gs *GameStore) CreateGame(room string) error {
	if _, ok := gs.games[room]; ok {
		return errors.New("game already exist")
	}
	gs.games[room] = &game{
		ws:         gs.WS,
		clients:    make(map[*gameClient]bool),
		register:   make(chan *gameClient),
		unregister: make(chan *gameClient),
		close:      make(chan bool),
	}
	go gs.games[room].run()
	return nil
}

type loginRequest struct {
	Type   string `json:"type"`
	Secret string `json:"secret"`
	Room   string `json:"room"`
}

func (gs *GameStore) handleConnection(c websocket.Connection) {
	client := &gameClient{}
	c.OnMessage(func(message []byte) {
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
