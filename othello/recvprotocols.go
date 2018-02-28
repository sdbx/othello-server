package othello

import (
	"encoding/json"

	"github.com/buger/jsonparser"
)

func protoErr(cli *Client, from string, msg string) {
	content, _ := json.Marshal(h{
		"type": "error",
		"from": from,
		"msg":  msg,
	})
	cli.Send <- content
}

func ping(room *Room, cli *Client, message []byte) {
	cli.Send <- []byte(`{"type":"pong"}`)
}

func login(room *Room, cli *Client, message []byte) {
	secret, err := jsonparser.GetString(message, "secret")
	if err != nil {
		protoErr(cli, "login", "json error")
	}
	user := room.hub.userStore.GetUserBySecret(secret)
	if user == nil {
		protoErr(cli, "login", "user doesn't exist")
	} else if cli.User != nil {
		protoErr(cli, "login", "login should be done once in a session")
	} else {
		cli.User = user
		content, _ := json.Marshal(h{
			"type":     "success",
			"from":     "login",
			"username": user.Name,
		})
		cli.Send <- content

		content, _ = json.Marshal(h{
			"type":     "connect",
			"username": user.Name,
		})
		room.Broadcast <- content
	}
}
