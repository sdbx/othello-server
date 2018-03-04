package room

import (
	"errors"
	"log"
	"sync"

	"github.com/sdbx/othello-server/othello/models"
	"github.com/sdbx/othello-server/othello/ws"
)

type (
	State     uint
	RoomStore struct {
		sync.Mutex
		*ws.WSStore
	}
	Room struct {
		*ws.WSRoom
		Participants uint
		State        State
		Black        string
		White        string
		King         string
	}
)

const (
	StateReady State = iota
	StateGame
)

func NewRoomStore(userStore models.UserStore) *RoomStore {
	rs := &RoomStore{
		WSStore: ws.NewWSStore(userStore),
	}
	rs.Handlers["enter"] = rs.enterHandler
	return rs
}

func (rs *RoomStore) CreateRoom(roomn string) error {
	log.Println(roomn, "room created")
	rs.Lock()
	if _, ok := rs.Rooms[roomn]; ok {
		return errors.New("room already exist")
	}
	room := &Room{
		WSRoom:       ws.NewWSRoom(roomn, rs.WSStore),
		State:        StateReady,
		Participants: 0,
	}
	rs.Rooms[roomn] = room
	rs.Unlock()
	go room.Run()
	return nil
}

func (r *Room) Register(cli *ws.Client) {
	r.Participants++
	r.WSRoom.Register(cli)
}

func (r *Room) Unregister(cli *ws.Client) {
	r.Participants--
	r.WSRoom.Unregister(cli)
}
