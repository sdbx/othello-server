package othello

import (
	websocket "github.com/kataras/go-websocket"
)

type gameClientType uint

const (
	gameClientSpectator gameClientType = iota + 1
	gameClientBlack
	gameClientWhite
)

type tile uint

const (
	gameTileNone tile = iota
	gameTileBlack
	gameTileWhite
)

type (
	coordinate struct {
		x uint
		y uint
	}

	move       string
	history    []move
	board      [][]tile
	gameClient struct {
		typ  gameClientType
		user *User
		room *game
	}

	game struct {
		name    string
		black   string
		white   string
		initial board
		board   board
		history history
		ws      websocket.Server

		clients    map[*gameClient]bool
		register   chan *gameClient
		unregister chan *gameClient
		close      chan bool
	}
)

func (g *game) emitMessage(message []byte) {
	for _, con := range g.ws.GetConnectionsByRoom(g.name) {
		con.EmitMessage(message)
	}
}

func (g *game) run() {
	for {
		select {
		case client := <-g.register:
			g.clients[client] = true
		case client := <-g.unregister:
			delete(g.clients, client)
		case <-g.close:
			return
		}
	}
}
