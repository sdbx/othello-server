package room

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/olebedev/emitter"
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/utils"
	"github.com/sdbx/othello-server/othello/ws"
)

type (
	State uint

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

func (r *Room) StartGame() (string, error) {
	r.Lock()
	defer r.Unlock()

	if r.info.State == StateGame {
		return "", errors.New("already ingame")
	}

	if r.info.Black == NoneUser || r.info.White == NoneUser {
		return "", errors.New("some color isn't selected")
	}

	name := r.runGame()

	return name, nil
}

func (r *Room) runGame() string {
	name := utils.GenKey()
	gam, _ := r.gameStore.CreateGame(name, r.info.Black, r.info.White, game.DefaultOthello{})

	r.info.State = StateGame
	r.info.Game = name

	gam.Emitter.On("end", func(e *emitter.Event) {
		r.info.State = StatePerparing
		r.info.Game = NoneGame
		r.Emit("gameend", ws.H{})
		go r.timeout()
	})

	r.Emit("gamestart", ws.H{
		"game": name,
	})

	return name
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
		return errors.New("no such user")
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
		return errors.New("no such user")
	}
	if cli.User.Name == r.info.King {
		return errors.New("what are you doing here sir?")
	}

	cli.Connection.Disconnect()
	r.Emit("action", ws.H{
		"action": "kick",
		"target": target,
	})

	return nil
}

func (r *Room) Register(cli ws.Client) {
	// deal with room
	r.Lock()
	r.info.Participants++
	r.lastConnected = time.Now()
	r.Unlock()

	r.WSRoom.Register(cli)

	// duplicate user
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

func (r *Room) pickNext() string {
	list := r.GetClientNames()
	return list[rand.Intn(len(list))]
}

func (r *Room) timeout() {
	t := time.NewTimer(time.Second * 30)
	<-t.C

	diff := time.Now().Sub(r.lastConnected)
	if diff >= time.Second*30 {
		r.Close()
	}
}
