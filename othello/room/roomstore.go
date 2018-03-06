package room

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/olebedev/emitter"
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/models"
	"github.com/sdbx/othello-server/othello/ws"
)

type (
	State     uint
	RoomStore struct {
		mu sync.Mutex
		*ws.WSStore
		gameStore *game.GameStore
	}
	Room struct {
		mu sync.Mutex
		*ws.WSRoom
		Participants  uint
		State         State
		Black         string
		White         string
		King          string
		Game          string
		gameStore     *game.GameStore
		lastConnected time.Time
	}
	keyStore struct {
		mu   sync.Mutex
		keys map[string]bool
	}
)

// TODO use db instead
var KeyStore = keyStore{
	keys: make(map[string]bool),
}

func genKey() string {
	return "abc"
}

func (s *keyStore) Gen() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for {
		key := genKey()
		if !s.keys[key] {
			s.keys[key] = true
			return key
		}
	}
}

const (
	StateReady State = iota
	StateGame
)

const NoneUser = "none"
const NoneGame = "none"

func (r *Room) StartGame() (string, error) {
	if r.State == StateGame {
		return "", errors.New("already ingame")
	}
	if r.Black == NoneUser || r.White == NoneUser {
		return "", errors.New("some color isn't selected")
	}
	key := KeyStore.Gen()
	gam, err := r.gameStore.CreateGame(key, r.Black, r.White, game.DefaultOthello{})
	if err != nil {
		return "", err
	}
	r.State = StateGame
	r.Game = key
	gam.Emitter.On("end", func(e *emitter.Event) {
		r.State = StateReady
		r.Game = NoneGame
		go func() {
			t := time.NewTimer(time.Second * 30)
			<-t.C
			diff := time.Now().Sub(r.lastConnected)
			if diff == time.Second*30 {
				r.Close()
			}
		}()
	})
	return key, nil
}

func (r *Room) GetClient(username string) *ws.Client {
	targets := r.GetClientsByName(username)
	if len(targets) == 0 {
		return nil
	}
	return targets[0]
}

func (r *Room) ChangeKing(target string) error {
	cli := r.GetClient(target)
	if cli == nil {
		return errors.New("user with target username doesn't exist")
	}
	r.King = target
	return nil
}

func (r *Room) ChangeColor(color string, target string) error {
	cli := r.GetClient(target)
	if cli == nil && target != NoneUser {
		return errors.New("user with target username doesn't exist")
	}
	if color == "black" {
		r.Black = target
		return nil
	}
	if color == "white" {
		r.White = target
		return nil
	}
	return errors.New("no such color")
}

func (r *Room) Kick(target string) error {
	cli := r.GetClient(target)
	if cli == nil {
		return errors.New("user with target username doesn't exist")
	}
	if cli.User.Name == r.King {
		return errors.New("what are you doing here master?")
	}
	cli.Connection.Disconnect()
	return nil
}

func (rs *RoomStore) CreateRoom(roomn string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if _, ok := rs.Rooms[roomn]; ok {
		return errors.New("room already exist")
	}
	log.Println(roomn, "room created")
	room := &Room{
		WSRoom:       ws.NewWSRoom(roomn, rs.WSStore),
		State:        StateReady,
		Participants: 0,
		Black:        NoneUser,
		White:        NoneUser,
		Game:         NoneGame,
		gameStore:    rs.gameStore,
	}
	rs.Rooms[roomn] = room
	go room.Run()
	return nil
}

func (r *Room) Register(cli *ws.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Participants++
	r.WSRoom.Register(cli)
	r.lastConnected = time.Now()
}

func (r *Room) Unregister(cli *ws.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Participants--
	r.WSRoom.Unregister(cli)
	if r.Participants == 0 && r.State == StateReady {
		r.Close()
	}
}

func NewRoomStore(userStore models.UserStore, gameStore *game.GameStore) *RoomStore {
	rs := &RoomStore{
		WSStore: ws.NewWSStore(userStore),
	}
	rs.gameStore = gameStore
	rs.Handlers = map[string]ws.WSListenHandler{
		"enter":   rs.enterHandler,
		"actions": rs.actionsHandler,
	}
	return rs
}
