package room

import (
	"errors"
	"log"
	"sync"

	websocket "github.com/kataras/go-websocket"
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/ws"
)

type RoomStore struct {
	sync.RWMutex
	*ws.WSStore
	gameStore *game.GameStore
	userConns map[string]websocket.Connection
}

func NewRoomStore(gameStore *game.GameStore) *RoomStore {
	rs := &RoomStore{
		WSStore:   ws.NewWSStore(),
		gameStore: gameStore,
		userConns: make(map[string]websocket.Connection),
	}

	rs.Handlers = map[string]ws.WSListenHandler{
		"enter":  rs.enterHandler,
		"action": rs.actionsHandler,
	}

	return rs
}

func (rs *RoomStore) CreateRoom(roomn string) (*Room, error) {
	rs.WSStore.Lock()
	defer rs.WSStore.Unlock()

	if _, ok := rs.Rooms[roomn]; ok {
		return nil, errors.New("room already exist")
	}

	log.Println(roomn, "room created")

	room := &Room{
		WSRoom:    ws.NewWSRoom(roomn, rs.WSStore),
		roomStore: rs,
		gameStore: rs.gameStore,
		info: RoomInfo{
			State:        StatePerparing,
			Participants: 0,
			Black:        NoneUser,
			White:        NoneUser,
			Game:         NoneGame,
		},
	}
	rs.Rooms[roomn] = room

	return room, nil
}

func (rs *RoomStore) GetInfos() []RoomInfo {
	rs.WSStore.RLock()
	defer rs.WSStore.RUnlock()

	list := []RoomInfo{}
	for name, room := range rs.Rooms {
		info := room.(*Room).GetInfo()
		info.Name = name
		list = append(list, info)
	}

	return list
}
