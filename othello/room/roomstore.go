package room

import (
	"errors"
	"log"
	"math/rand"
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

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genKey(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s *keyStore) Gen() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for {
		key := genKey(10)
		if !s.keys[key] {
			s.keys[key] = true
			return key
		}
	}
}

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
		r.State = StatePerparing
		r.Game = NoneGame
		r.Emit("gameend", ws.H{})
		go func() {
			t := time.NewTimer(time.Second * 30)
			<-t.C
			diff := time.Now().Sub(r.lastConnected)
			if diff == time.Second*30 {
				r.Close()
			}
		}()
	})
	r.Emit("gamestart", ws.H{
		"game": key,
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
	r.Emit("actions", ws.H{
		"action": "king",
		"target": target,
	})
	return nil
}

func (r *Room) ChangeColor(color string, target string) error {
	if r.State == StateGame {
		return errors.New("changin color during the game is not allowed")
	}
	if color != "black" && color != "white" {
		return errors.New("no such color")
	}
	cli := r.GetClient(target)
	if cli == nil && target != NoneUser {
		return errors.New("user with target username doesn't exist")
	}
	if color == "black" {
		r.Black = target
	}
	if color == "white" {
		r.White = target
	}
	r.Emit("actions", ws.H{
		"action": "color",
		"to":     target,
		"color":  color,
	})
	return nil
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
	r.Emit("actions", ws.H{
		"action": "kick",
		"target": target,
	})
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
		State:        StatePerparing,
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

func (r *Room) pickNext(current string) string {
	list := r.GetClientNames()
	next := ""
	for {
		next = list[rand.Intn(len(list))]
		if next != current {
			return next
		}
	}
}

func (r *Room) Unregister(cli *ws.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	name := cli.User.Name
	r.Participants--
	r.WSRoom.Unregister(cli)
	if r.Participants == 0 && r.State == StatePerparing {
		r.Close()
		return
	}
	if r.State == StatePerparing {
		if r.Black == name {
			r.ChangeColor("black", NoneUser)
		}
		if r.White == name {
			r.ChangeColor("white", NoneUser)
		}
	}
	if r.King == name {
		r.ChangeKing(r.pickNext(r.King))
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
