package ws

import (
	"encoding/json"
	"fmt"

	websocket "github.com/kataras/go-websocket"
)

type enterRequest struct {
	Type   string `json:"type"`
	Secret string `json:"secret"`
	Room   string `json:"game"`
}

func WSEnterHandler(c websocket.Connection, client *Client, message []byte) {
	if client.user != nil {
		c.EmitMessage([]byte(onceMsg))
		return
	}
	req := enterRequest{}
	err := json.Unmarshal(message, &req)
	if err != nil {
		c.EmitMessage([]byte(jsonErrorMsg))
		return
	}
	store := client.room.Store()
	user := store.UserStore.GetUserBySecret(req.Secret)
	if user == nil {
		c.EmitMessage([]byte(fmt.Sprintf(userNoMsg, "enter")))
		return
	}
	room, ok := store.Rooms[req.Room]
	if !ok {
		c.EmitMessage([]byte(fmt.Sprintf(gameNoMsg)))
		return
	}
	client.user = user
	client.room = room
	c.Join(room.Name())
	room.Emit("connect", H{
		"username": client.user.Name,
	})
	room.Register() <- client
}

func WSPingHandler(c websocket.Connection, client *Client, message []byte) {
	c.EmitMessage([]byte(pongMsg))
}
