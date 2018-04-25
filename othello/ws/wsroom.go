package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	websocket "github.com/kataras/go-websocket"
	"github.com/sdbx/othello-server/othello/dbs"
)

type (
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
)

func (cli *Client) EmitError(msg string, from string) {
	cli.Connection.EmitMessage([]byte(fmt.Sprintf(`{"type":"error","msg":"%s","from":"%s"}`, msg, from)))
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
		con.EmitMessage(content)
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
