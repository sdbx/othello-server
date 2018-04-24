package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/buger/jsonparser"
	websocket "github.com/kataras/go-websocket"
	"github.com/sdbx/othello-server/othello/dbs"
)

type (
	H      map[string]interface{}
	Client struct {
		ID         int
		Authed     bool
		User       dbs.User
		Room       Room
		Connection websocket.Connection
	}

	Room interface {
		Name() string
		Register(Client)
		Unregister(Client)
		Emit(string, H)
		Store() *WSStore
		GetClientsByName(name string) []Client
		GetClientNames() []string
	}

	WSRoom struct {
		sync.RWMutex
		maxid   int
		clients map[string]map[int]Client
		name    string
		store   *WSStore
	}

	WSListenHandler func(Client, []byte) Client
	WSStore         struct {
		sync.RWMutex
		WS       websocket.Server
		Handlers map[string]WSListenHandler
		Rooms    map[string]Room
	}
)

func (cli *Client) EmitError(msg string, from string) {
	cli.Connection.EmitMessage([]byte(fmt.Sprintf(`{"type":"error","msg":"%s","from":"%s"}`, msg, from)))
}

func NewWSStore() *WSStore {
	gs := &WSStore{
		WS:       websocket.New(websocket.Config{}),
		Rooms:    make(map[string]Room),
		Handlers: make(map[string]WSListenHandler),
	}
	gs.Handlers["ping"] = WSPingHandler
	gs.WS.OnConnection(gs.handleConnection)
	return gs
}

const jsonErrorMsg = `{"type":"error","msg":"json error","from":"none"}`
const typeErrorMsg = `{"type":"error","msg":"no such type","from":"none"}`

func (rs *WSStore) handleConnection(c websocket.Connection) {
	client := Client{
		Connection: c,
	}
	c.OnMessage(func(message []byte) {
		log.Println("websocket recieved:", string(message))
		typ, err := jsonparser.GetString(message, "type")
		if err != nil {
			c.EmitMessage([]byte(jsonErrorMsg))
			return
		}

		if fun, ok := rs.Handlers[typ]; !ok {
			c.EmitMessage([]byte(typeErrorMsg))
		} else {
			client = fun(client, message)
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

func (rs *WSStore) Enter(cli Client, user dbs.User, roomn string) (Client, error) {
	if cli.Authed {
		return cli, errors.New("enter should be occured once in a session")
	}

	rs.RLock()
	room, ok := rs.Rooms[roomn]
	rs.RUnlock()

	if !ok {
		return cli, errors.New("no such room")
	}
	cli.Authed = true
	cli.User = user
	cli.Room = room
	cli.Connection.Join(room.Name())
	room.Register(cli)
	room.Emit("connect", H{
		"username": user.Name,
	})
	return cli, nil
}

func NewWSRoom(name string, store *WSStore) *WSRoom {
	return &WSRoom{
		clients: make(map[string]map[int]Client),
		name:    name,
		store:   store,
	}
}

func (r *WSRoom) GetClientNames() []string {
	r.RLock()
	defer r.RUnlock()

	list := []string{}
	for name := range r.clients {
		list = append(list, name)
	}
	return list
}

func (r *WSRoom) GetClientsByName(name string) []Client {
	r.RLock()
	defer r.RUnlock()

	if clili, ok := r.clients[name]; ok {
		list := []Client{}
		for _, cli := range clili {
			list = append(list, cli)
		}
		return list
	}
	return []Client{}
}

func (r *WSRoom) Close() {
	r.Store().Lock()
	delete(r.Store().Rooms, r.Name())
	r.Store().Unlock()
	for _, conn := range r.Store().WS.GetConnectionsByRoom(r.Name()) {
		// will trigger unregisters
		conn.Disconnect()
	}
}

func (r *WSRoom) Emit(typ string, ho H) {
	ho["type"] = typ
	content, _ := json.Marshal(ho)
	r.EmitMsg(content)
}

func (r *WSRoom) EmitMsg(content []byte) {
	for _, con := range r.Store().WS.GetConnectionsByRoom(r.Name()) {
		go con.EmitMessage(content)
	}
	log.Println("websocket sent from", r.name, ":", string(content))
}

func (r *WSRoom) Name() string {
	return r.name
}

func (r *WSRoom) Store() *WSStore {
	return r.store
}

func (r *WSRoom) Register(cli Client) {
	r.Lock()
	defer r.Unlock()

	name := cli.User.Name
	_, ok := r.clients[name]
	if !ok {
		r.clients[name] = make(map[int]Client)
	}

	cli.ID = r.maxid
	r.clients[name][r.maxid] = cli
	r.maxid++
}

func (r *WSRoom) Unregister(cli Client) {
	r.Lock()
	defer r.Unlock()

	name := cli.User.Name
	delete(r.clients[name], cli.ID)
	if len(r.clients[name]) == 0 {
		delete(r.clients, name)
	}
}
