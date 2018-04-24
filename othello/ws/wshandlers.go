package ws

const pongMsg = `{"type":"pong"}`

func WSPingHandler(client Client, message []byte) Client {
	client.Connection.EmitMessage([]byte(pongMsg))
	return client
}
