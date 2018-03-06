package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/buger/jsonparser"
	websocket "github.com/kataras/go-websocket"
	"github.com/sdbx/othello-server/othello/models"
)

type (
	H      map[string]interface{}
	Client struct {
		User       *models.User
		Room       Room
		Connection websocket.Connection
	}

	Room interface {
		Name() string
		Register(*Client)
		Unregister(*Client)
		Emit(string, H)
		Store() *WSStore
		GetClientsByName(name string) []*Client
	}

	WSRoom struct {
		clients    map[string]map[*Client]bool
		name       string
		store      *WSStore
		register   chan *Client
		unregister chan *Client
		close      chan bool
	}
	WSListenHandler func(*Client, []byte)
	WSStore         struct {
		WS        websocket.Server
		Handlers  map[string]WSListenHandler
		UserStore models.UserStore
		Rooms     map[string]Room
	}
)

func (cli *Client) EmitError(msg string, from string) {
	cli.Connection.EmitMessage([]byte(fmt.Sprintf(`{"type":"error","msg":"%s","from":"%s"}`, msg, from)))
}

func NewWSStore(userStore models.UserStore) *WSStore {
	gs := &WSStore{
		WS:        websocket.New(websocket.Config{}),
		UserStore: userStore,
		Rooms:     make(map[string]Room),
		Handlers:  make(map[string]WSListenHandler),
	}
	gs.Handlers["ping"] = WSPingHandler
	gs.WS.OnConnection(gs.handleConnection)
	return gs
}

const jsonErrorMsg = `
{
	"type":"error",
	"msg":"json error",
	"from":"none"
}
`
const typeErrorMsg = `
{
	"type":"error",
	"msg":"no such type",
	"from":"none"
}
`

func (rs *WSStore) handleConnection(c websocket.Connection) {
	client := &Client{
		Connection: c,
	}
	c.OnMessage(func(message []byte) {
		log.Println("websocket recieved:", string(message))
		typ, err := jsonparser.GetString(message, "type")
		if err != nil {
			c.EmitMessage([]byte(jsonErrorMsg))
			return
		}
		fmt.Println(rs.Handlers)
		if fun, ok := rs.Handlers[typ]; !ok {
			c.EmitMessage([]byte(typeErrorMsg))
			return
		} else {
			fun(client, message)
		}
	})
	c.OnDisconnect(func() {
		if client.Room != nil {
			client.Room.Emit("disconnect", H{
				"username": client.User.Name,
			})
			client.Room.Unregister(client)
		}
	})
}

func (r *WSStore) Enter(cli *Client, user *models.User, roomn string) error {
	if cli.User != nil {
		return errors.New("enter should be occured once in a session")
	}
	room, ok := r.Rooms[roomn]
	if !ok {
		return errors.New("no such room")
	}
	cli.User = user
	cli.Room = room
	cli.Connection.Join(room.Name())
	room.Emit("connect", H{
		"username": user.Name,
	})
	room.Register(cli)
	return nil
}

func NewWSRoom(name string, store *WSStore) *WSRoom {
	return &WSRoom{
		clients:    make(map[string]map[*Client]bool),
		name:       name,
		store:      store,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		close:      make(chan bool),
	}
}

func (r *WSRoom) GetClientsByName(name string) []*Client {
	if clili, ok := r.clients[name]; ok {
		list := []*Client{}
		for cli := range clili {
			list = append(list, cli)
		}
		return list
	}
	return []*Client{}
}

func (r *WSRoom) Close() {
	go func() {
		r.close <- true
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
	r.EmitMsg(content)
	log.Println("websocket sent from", r.name, ":", string(content))
}

func (r *WSRoom) EmitMsg(content []byte) {
	for _, con := range r.Store().WS.GetConnectionsByRoom(r.Name()) {
		con.EmitMessage(content)
	}
	log.Println("websocket sent from", r.name, ":", string(content))
}

func (r *WSRoom) Run() {
	for {
		select {
		case client := <-r.register:
			name := client.User.Name
			_, ok := r.clients[name]
			if !ok {
				r.clients[name] = make(map[*Client]bool)
			}
			r.clients[name][client] = true
		case client := <-r.unregister:
			delete(r.clients[client.User.Name], client)
			fmt.Printf("%v", r.clients[client.User.Name])
		case <-r.close:
			return
		}
	}
}

func (r *WSRoom) Name() string {
	return r.name
}

func (r *WSRoom) Store() *WSStore {
	return r.store
}

func (r *WSRoom) Register(cli *Client) {
	r.register <- cli
}

func (r *WSRoom) Unregister(cli *Client) {
	r.unregister <- cli
}
