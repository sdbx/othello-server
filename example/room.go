package main

import (
	"fmt"
	"net/http"

	websocket "github.com/kataras/go-websocket"
)

func main() {
	ws := websocket.New(websocket.Config{})
	http.Handle("/abc", ws.Handler())
	ws.OnConnection(handleWebsocketConnection)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

var myChatRoom = "room1"

func handleWebsocketConnection(c websocket.Connection) {

	c.Join(myChatRoom)

	c.OnMessage(func(message []byte) {
		fmt.Println(string(message))
	})

	c.OnDisconnect(func() {
		fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
	})
}
