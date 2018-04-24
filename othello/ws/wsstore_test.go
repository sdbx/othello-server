package ws

import (
	"testing"
	"time"

	"github.com/posener/wstest"
	"github.com/sdbx/othello-server/othello/dbs"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	a := assert.New(t)
	s := NewWSStore()
	d := wstest.NewDialer(s.WS.Handler())
	c, _, err := d.Dial("ws://asdf", nil)
	if err != nil {
		panic(err)
	}
	c.WriteJSON(H{
		"type": "ping",
	})
	resp := H{}
	c.ReadJSON(&resp)
	a.Equal(resp["type"], "pong")
}
func TestEnter(t *testing.T) {
	a := assert.New(t)
	s := NewWSStore()
	s.Rooms["asdf"] = NewWSRoom("asdf", s)
	s.Handlers["enter"] = func(client Client, message []byte) Client {
		client, _ = s.Enter(client, dbs.User{Name: "asdf"}, "asdf")
		return client
	}

	d := wstest.NewDialer(s.WS.Handler())
	c, _, _ := d.Dial("ws://asdf", nil)
	c.WriteJSON(H{
		"type": "enter",
	})

	// emit
	resp := H{}
	c.ReadJSON(&resp)
	a.Equal(resp["type"], "connect")
	time.Sleep(time.Second / 2)

	// getclientsnames
	names := s.Rooms["asdf"].GetClientNames()
	a.Equal(len(names), 1)

	// getclientsbyname
	clients := s.Rooms["asdf"].GetClientsByName("asdf")
	a.Equal(len(clients), 1)

	// close
	c.Close()
	time.Sleep(time.Second / 2)
	names = s.Rooms["asdf"].GetClientNames()
	a.Equal(len(names), 0)
}
