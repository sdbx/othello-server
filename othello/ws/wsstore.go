package ws

import (
	"errors"
	"log"
	"sync"

	"github.com/buger/jsonparser"
	websocket "github.com/kataras/go-websocket"
	"github.com/sdbx/othello-server/othello/dbs"
)

type (
	H map[string]interface{}

	WSListenHandler func(Client, []byte) Client
	WSStore         struct {
		sync.RWMutex
		WS       websocket.Server
		Handlers map[string]WSListenHandler
		Rooms    map[string]Room
	}
)

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

func (rs *WSStore) GetRoom(name string) Room {
	rs.RLock()
	defer rs.RUnlock()
	if room, ok := rs.Rooms[name]; !ok {
		return nil
	} else {
		return room
	}
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
