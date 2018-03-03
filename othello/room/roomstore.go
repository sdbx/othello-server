package room

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
	h         map[string]interface{}
	RoomStore struct {
		WS        websocket.Server
		userStore models.UserStore
		rooms     map[string]*Room
	}
	roomClient struct {
		user *models.User
		room *Room
	}
	Room struct {
		Name  string
		store *RoomStore
		ws    websocket.Server

		clients    map[*roomClient]bool
		register   chan *roomClient
		unregister chan *roomClient
	}
)

func NewRoomStore(userStore models.UserStore) *RoomStore {
	return &RoomStore{
		WS:        websocket.New(websocket.Config{}),
		userStore: userStore,
		rooms:     make(map[string]*Room),
	}
}

func (r *Room) emitMessage(message []byte) {
	log.Println("websocket sent from", r.name, ":", string(message))
	for _, con := range r.ws.GetConnectionsByRoom(r.name) {
		con.EmitMessage(message)
	}
}

func (r *Room) emit(typ string, ho h) {
	ho["type"] = typ
	content, _ := json.Marshal(ho)
	r.emitMessage(content)
}

func (rs *RoomStore) CreateGame(roomn string) (*Room, error) {
	log.Println(roomn, "room created")
	if _, ok := rs.rooms[roomn]; ok {
		return nil, errors.New("room already exist")
	}
	room := &Room{
		Name:       roomn,
		store:      rs,
		ws:         rs.WS,
		clients:    make(map[*roomClient]bool),
		register:   make(chan *roomClient),
		unregister: make(chan *roomClient),
	}
	go room.run()
	return room, nil
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			delete(r.clients, client)
		}
	}
}

type loginRequest struct {
	Type   string `json:"type"`
	Secret string `json:"secret"`
	Room   string `json:"room"`
}

// 중복 실화?
func (rs *RoomStore) handleConnection(c websocket.Connection) {
	client := &roomClient{}
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
			user := rs.userStore.GetUserBySecret(req.Secret)
			if user == nil {
				c.EmitMessage([]byte(fmt.Sprintf(userNoMsg, "enter")))
				return
			}
			room, ok := rs.rooms[req.Room]
			if !ok {
				c.EmitMessage([]byte(fmt.Sprintf(gameNoMsg)))
				return
			}
			client.user = user
			client.room = room
			c.Join(room.Name)
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
