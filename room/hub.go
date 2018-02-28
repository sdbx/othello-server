package room

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sdbx/othello-server/othello"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type roomHub map[string]*Room

func (rs *roomHub) AddRoom(roomname string, usersecret string) (*Room, error) {
	user := service.UserStore.GetUserBySecret(usersecret)
	if user == nil {
		return nil, errors.New("user doesn't exist")
	}
	if user.Status != othello.None {
		return nil, errors.New("user is already in room")
	}
	if _, ok := (*rs)[roomname]; ok {
		return nil, errors.New("room already exist")
	}
	room := &Room{
		Name:       roomname,
		Clients:    make(map[*Client]bool),
		Ready:      false,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	(*rs)[roomname] = room
	return room, nil
}

var RoomHub roomHub = make(map[string]*Room)

type Room struct {
	Name    string
	Clients map[*Client]bool
	Ready   bool

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (h *Room) getClient(username string) *Client {
	for item := range h.Clients {
		if item.User != nil {
			if item.User.Name == username {
				return item
			}
		}
	}
	return nil
}

func (h *Room) run() {
	for {
		select {
		case client := <-h.register:
			h.Ready = true
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				if client.User != nil {
					client.User.Status = othello.None
				}
				delete(h.Clients, client)
				close(client.send)
			}
			if len(h.Clients) == 0 {
				delete(RoomHub, h.Name)
				return
			}
		case message := <-h.broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Room) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		User: nil,
		room: h,
		conn: conn,
		send: make(chan []byte, 256),
	}
	h.register <- client

	go client.write()
	go client.read()
}
