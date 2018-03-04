package ws

const pongMsg = `
{
	"type":"pong"
}
`

func WSPingHandler(client *Client, message []byte) {
	client.Connection.EmitMessage([]byte(pongMsg))
}
