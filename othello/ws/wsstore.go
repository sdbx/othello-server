package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/buger/jsonparser"
	websocket "github.com/kataras/go-websocket"
	"github.com/sdbx/othello-server/othello/models"
)

type (
	H      map[string]interface{}
	Client struct {
		user *models.User
		room Room
	}

	Room interface {
		Name() string
		Register() chan *Client
		Unregister() chan *Client
		Emit(string, H)
		Store() *WSStore
	}

	WSRoom struct {
		Clients     map[*Client]bool
		_name       string
		_store      *WSStore
		_register   chan *Client
		_unregister chan *Client
	}
	WSListenHandler func(websocket.Connection, *Client, []byte)
	WSStore         struct {
		WS        websocket.Server
		Handlers  map[string]WSListenHandler
		UserStore models.UserStore
		Rooms     map[string]Room
	}
)

func NewWSStore(userStore models.UserStore) *WSStore {
	gs := &WSStore{
		WS:        websocket.New(websocket.Config{}),
		UserStore: userStore,
		Rooms:     make(map[string]Room),
		Handlers:  make(map[string]WSListenHandler),
	}
	gs.Handlers["ping"] = WSPingHandler
	gs.Handlers["enter"] = WSEnterHandler
	gs.WS.OnConnection(gs.handleConnection)
	return gs
}

func (rs *WSStore) handleConnection(c websocket.Connection) {
	client := &Client{}
	c.OnMessage(func(message []byte) {
		log.Println("websocket recieved:", string(message))
		typ, err := jsonparser.GetString(message, "type")
		if err != nil {
			c.EmitMessage([]byte(jsonErrorMsg))
			return
		}
		if fun, ok := rs.Handlers[typ]; !ok {
			c.EmitMessage([]byte(typeErrorMsg))
			return
		} else {
			fun(c, client, message)
		}
	})
	c.OnDisconnect(func() {
		if client.room != nil {
			client.room.Emit("disconnect", H{
				"username": client.user.Name,
			})
			client.room.Unregister() <- client
		}
	})
}

func NewWSRoom(name string, store *WSStore) *WSRoom {
	return &WSRoom{
		Clients:     make(map[*Client]bool),
		_name:       name,
		_store:      store,
		_register:   make(chan *Client),
		_unregister: make(chan *Client),
	}
}

func (r *WSRoom) Close() {
	go func() {
		time.Sleep(time.Second)
		for _, conn := range r.Store().WS.GetConnectionsByRoom(r.Name()) {
			conn.Disconnect()
		}
		delete(r.Store().Rooms, r.Name())
	}()
}

func (r *WSRoom) Emit(typ string, ho H) {
	ho["type"] = typ
	content, _ := json.Marshal(ho)
	for _, con := range r.Store().WS.GetConnectionsByRoom(r.Name()) {
		con.EmitMessage(content)
	}
	log.Println("websocket sent from", r._name, ":", string(content))
}

func (r *WSRoom) Run() {
	for {
		select {
		case client := <-r._register:
			r.Clients[client] = true
		case client := <-r._unregister:
			delete(r.Clients, client)
		}
	}
}

func (r *WSRoom) Name() string {
	return r._name
}

func (r *WSRoom) Register() chan *Client {
	return r._register
}

func (r *WSRoom) Unregister() chan *Client {
	return r._unregister
}

func (r *WSRoom) Store() *WSStore {
	return r._store
}
