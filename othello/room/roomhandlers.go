package room

import (
	"encoding/json"

	"github.com/sdbx/othello-server/othello/ws"
)

const jsonErrorMsg = `
{
	"type":"error",
	"msg":"json error",
	"from":"none"
}
`
const enterErrorMsg = `
{
	"type":"error",
	"msg":"%s",
	"from":"enter"
}
`

type enterRequest struct {
	_      string `json:"type"`
	Secret string `json:"secret"`
	Room   string `json:"room"`
}

func (rs *RoomStore) enterHandler(cli *ws.Client, message []byte) {
	req := enterRequest{}
	err := json.Unmarshal(message, &req)
	if err != nil {
		cli.EmitError("json error", "enter")
		return
	}
	user := rs.UserStore.GetUserBySecret(req.Secret)
	if user == nil {
		cli.EmitError("user doesn't exist", "enter")
		return
	}

	if room, ok := rs.Rooms[req.Room]; !ok {
		rs.CreateRoom(req.Room)
		rs.Rooms[req.Room].(*Room).King = user.Name
	} else if len(room.GetClientsByName(user.Name)) != 0 {
		cli.EmitError("allows only one connection per user", "enter")
		return
	}

	err = rs.Enter(cli, user, req.Room)
	if err != nil {
		cli.EmitError(err.Error(), "enter")
	}
}
