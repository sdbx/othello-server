package room

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	websocket "github.com/kataras/go-websocket"
	"github.com/olebedev/emitter"
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/utils"
	"github.com/sdbx/othello-server/othello/ws"
)

type (
	State     uint
	RoomStore struct {
		sync.RWMutex
		*ws.WSStore
		gameStore *game.GameStore
		userConns map[string]websocket.Connection
	}

	RoomInfo struct {
		Name         string `json:"name"`
		Participants uint   `json:"participants"`
		State        State  `json:"state"`
		Black        string `json:"black"`
		White        string `json:"white"`
		King         string `json:"king"`
		Game         string `json:"game"`
	}
	Room struct {
		sync.RWMutex
		*ws.WSRoom
		roomStore     *RoomStore
		gameStore     *game.GameStore
		lastConnected time.Time
		info          RoomInfo
	}
)

const (
	StatePerparing State = iota
	StateGame
)

func (s State) String() string {
	if s == StatePerparing {
		return "preparing"
	}
	return "ingame"
}

const NoneUser = "none"
const NoneGame = "none"

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

func (r *Room) timeout() {
	t := time.NewTimer(time.Second * 30)
	<-t.C
	diff := time.Now().Sub(r.lastConnected)
	if diff >= time.Second*30 {
		r.Close()
	}
}

func (r *Room) StartGame() (string, error) {
	r.Lock()
	defer r.Unlock()

	if r.info.State == StateGame {
		return "", errors.New("already ingame")
	}
	if r.info.Black == NoneUser || r.info.White == NoneUser {
		return "", errors.New("some color isn't selected")
	}
	key := utils.GenKey()
	gam, err := r.gameStore.CreateGame(key, r.info.Black, r.info.White, game.DefaultOthello{})
	if err != nil {
		return "", err
	}
	r.info.State = StateGame
	r.info.Game = key
	gam.Emitter.On("end", func(e *emitter.Event) {
		r.info.State = StatePerparing
		r.info.Game = NoneGame
		r.Emit("gameend", ws.H{})
		go r.timeout()
	})
	r.Emit("gamestart", ws.H{
		"game": key,
	})
	return key, nil
}

func (r *Room) GetClient(username string) (ws.Client, error) {
	targets := r.GetClientsByName(username)
	if len(targets) == 0 {
		return ws.Client{}, errors.New("no such user")
	}
	return targets[0], nil
}

func (r *Room) ChangeKing(target string) error {
	r.Lock()
	defer r.Unlock()

	_, err := r.GetClient(target)
	if err != nil {
		return errors.New("user with target username doesn't exist")
	}
	r.info.King = target
	r.Emit("action", ws.H{
		"action": "king",
		"target": target,
	})
	return nil
}

func (r *Room) ChangeColor(color string, target string) error {
	r.Lock()
	defer r.Unlock()

	if r.info.State == StateGame {
		return errors.New("changin color during the game is not allowed")
	}
	if color != "black" && color != "white" {
		return errors.New("no such color")
	}
	_, err := r.GetClient(target)
	if err == nil && target != NoneUser {
		return errors.New("user with target username doesn't exist")
	}
	if color == "black" {
		r.info.Black = target
	}
	if color == "white" {
		r.info.White = target
	}
	r.Emit("action", ws.H{
		"action": "color",
		"to":     target,
		"color":  color,
	})
	return nil
}

func (r *Room) Kick(target string) error {
	r.Lock()
	defer r.Unlock()

	cli, err := r.GetClient(target)
	if err != nil {
		return errors.New("user with target username doesn't exist")
	}
	if cli.User.Name == r.info.King {
		return errors.New("what are you doing here master?")
	}
	cli.Connection.Disconnect()
	r.Emit("action", ws.H{
		"action": "kick",
		"target": target,
	})
	return nil
}

func (rs *RoomStore) CreateRoom(roomn string) error {
	rs.WSStore.Lock()
	defer rs.WSStore.Unlock()
	if _, ok := rs.Rooms[roomn]; ok {
		return errors.New("room already exist")
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
	return nil
}

func (r *Room) Register(cli ws.Client) {
	// deal with room
	r.Lock()
	r.info.Participants++
	r.lastConnected = time.Now()
	r.Unlock()

	r.WSRoom.Register(cli)

	// disconnect duplicate user
	store := r.roomStore
	store.RLock()
	conn, ok := store.userConns[cli.User.Secret]
	store.RUnlock()
	if ok {
		conn.Disconnect()
	}
	store.Lock()
	store.userConns[cli.User.Secret] = cli.Connection
	store.Unlock()

}

func (r *Room) pickNext() string {
	list := r.GetClientNames()
	return list[rand.Intn(len(list))]
}

func (r *Room) Unregister(cli ws.Client) {
	// duplicate user
	store := r.roomStore
	store.Lock()
	delete(store.userConns, cli.User.Secret)
	store.Unlock()

	r.WSRoom.Unregister(cli)

	r.Lock()
	info := r.info
	r.info.Participants--
	if r.info.Participants == 0 && r.info.State == StatePerparing {
		r.Unlock()
		r.Close()
		return
	}
	r.Unlock()
	name := cli.User.Name
	if info.State == StatePerparing {
		if info.Black == name {
			r.ChangeColor("black", NoneUser)
		}
		if info.White == name {
			r.ChangeColor("white", NoneUser)
		}
	}
	if info.King == name {
		r.ChangeKing(r.pickNext())
	}
}

func (r *Room) GetInfo() RoomInfo {
	r.RLock()
	defer r.RUnlock()
	return r.info
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
