package othello

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Hub struct {
	Rooms        map[string]*Room
	ReadHandlers map[string]func(*Room, *Client, []byte)
	userStore    *UserStore
}

func NewHub(userStore *UserStore) *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
		ReadHandlers: map[string]func(*Room, *Client, []byte){
			"ping":  ping,
			"login": login,
		},
		userStore: userStore,
	}
}

func (h *Hub) AddRoom(roomname string, usersecret string) (*Room, error) {
	user := h.userStore.GetUserBySecret(usersecret)
	if user == nil {
		return nil, errors.New("user doesn't exist")
	}
	if user.Status != None {
		return nil, errors.New("user is already in room")
	}
	if _, ok := h.Rooms[roomname]; ok {
		return nil, errors.New("room already exist")
	}
	room := &Room{
		Name:       roomname,
		Clients:    make(map[*Client]bool),
		Ready:      false,
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		hub:        h,
	}
	h.Rooms[roomname] = room
	return room, nil
}

type Room struct {
	Name      string
	Clients   map[*Client]bool
	Ready     bool
	Broadcast chan []byte

	hub        *Hub
	register   chan *Client
	unregister chan *Client
}

func (h *Room) GetClient(username string) *Client {
	for item := range h.Clients {
		if item.User != nil {
			if item.User.Name == username {
				return item
			}
		}
	}
	return nil
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.Ready = true
			r.Clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.Clients[client]; ok {
				if client.User != nil {
					client.User.Status = None
				}
				close(client.Send)
				delete(r.Clients, client)
			}
			if len(r.Clients) == 0 {
				delete(r.hub.Rooms, r.Name)
				return
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}

func (r *Room) ServeWs(w http.ResponseWriter, re *http.Request) {
	conn, err := upgrader.Upgrade(w, re, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		User: nil,
		Send: make(chan []byte, 256),
		room: r,
		conn: conn,
	}
	r.register <- client

	go client.write()
	go client.read()
}
