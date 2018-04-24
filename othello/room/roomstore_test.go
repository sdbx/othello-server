package room_test

import (
	"testing"
	"time"

	"github.com/posener/wstest"
	"github.com/sdbx/othello-server/othello/dbs"
	"github.com/sdbx/othello-server/othello/game"
	"github.com/sdbx/othello-server/othello/room"
	"github.com/sdbx/othello-server/othello/ws"
	"github.com/stretchr/testify/assert"
)

func prepareDB() (string, string) {
	dbs.Clear()
	user := dbs.User{
		Name:   "hello",
		UserID: "asdf1",
	}
	dbs.AddUser(&user)
	sec1 := user.Secret
	user = dbs.User{
		Name:   "world",
		UserID: "asdf2",
	}
	sec2 := user.Secret
	return sec1, sec2
}

func TestConnect(t *testing.T) {
	sec1, _ := prepareDB()
	a := assert.New(t)
	g := game.NewGameStore()
	s := room.NewRoomStore(g)
	d := wstest.NewDialer(s.WS.Handler())

	c, _, _ := d.Dial("ws://asdf", nil)
	c.WriteJSON(ws.H{
		"type":   "enter",
		"secret": sec1,
		"room":   "asdf",
	})
	resp := ws.H{}
	c.ReadJSON(&resp)
	a.Equal(resp["type"].(string), "connect")

	// duplicate same room
	closed1 := false
	go func() {
		for {
			if _, _, err := c.NextReader(); err != nil {
				closed1 = true
				break
			}
		}
	}()

	d2 := wstest.NewDialer(s.WS.Handler())
	c2, _, _ := d2.Dial("ws://asdf", nil)
	c2.WriteJSON(ws.H{
		"type":   "enter",
		"secret": sec1,
		"room":   "asdf",
	})

	// duplicate different room
	closed2 := false
	go func() {
		for {
			if _, _, err := c2.NextReader(); err != nil {
				closed2 = true
				break
			}
		}
	}()

	d3 := wstest.NewDialer(s.WS.Handler())
	c3, _, _ := d3.Dial("ws://asdf", nil)
	c3.WriteJSON(ws.H{
		"type":   "enter",
		"secret": sec1,
		"room":   "asdf2",
	})

	time.Sleep(time.Second)
	a.Equal(closed1, true)
	a.Equal(closed2, true)
}
