package room

import (
	"encoding/json"

	"github.com/buger/jsonparser"
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

func (rs *RoomStore) actionsHandler(cli *ws.Client, message []byte) {
	room := cli.Room.(*Room)
	if room.King != cli.User.Name {
		cli.EmitError("not enough permission", "action")
		return
	}
	typ, err := jsonparser.GetString(message, "action")
	if err != nil {
		cli.EmitError(err.Error(), "actions")
		return
	}
	switch typ {
	case "color":
		req := struct {
			Color    string `json:"color"`
			Username string `json:"username"`
		}{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			cli.EmitError(err.Error(), "color")
			return
		}

		err = room.ChangeColor(req.Color, req.Username)
		if err != nil {
			cli.EmitError(err.Error(), "color")
			return
		}
		room.EmitMsg(message)
	case "kick":
		req := struct {
			Target string `json:"target"`
		}{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			cli.EmitError(err.Error(), "kick")
			return
		}
		err = room.Kick(req.Target)
		if err != nil {
			cli.EmitError(err.Error(), "kick")
		}
		room.EmitMsg(message)
	case "king":
		req := struct {
			Target string `json:"target"`
		}{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			cli.EmitError(err.Error(), "king")
			return
		}
		room.ChangeKing(req.Target)
		if err != nil {
			cli.EmitError(err.Error(), "king")
			return
		}
		room.EmitMsg(message)
	case "start":
		name, err := room.StartGame()
		if err != nil {
			cli.EmitError(err.Error(), "start")
			return
		}
		room.Emit("start", ws.H{
			"game": name,
		})
	default:
		cli.EmitError("no such action", "actions")
	}
}
