package room

import (
	"encoding/json"

	"github.com/buger/jsonparser"
	"github.com/sdbx/othello-server/othello/dbs"
	"github.com/sdbx/othello-server/othello/ws"
)

const jsonErrorMsg = `{"type":"error","msg":"json error","from":"none"}`
const enterErrorMsg = `{"type":"error","msg":"%s","from":"enter"}`

type enterRequest struct {
	Secret string `json:"secret"`
	Room   string `json:"room"`
}

func (rs *RoomStore) enterHandler(cli ws.Client, message []byte) ws.Client {
	req := enterRequest{}

	err := json.Unmarshal(message, &req)
	if err != nil {
		cli.EmitError("json error", "enter")
		return cli
	}

	user, err := dbs.GetUserBySecret(req.Secret)
	if err != nil {
		cli.EmitError("user doesn't exist", "enter")
		return cli
	}

	if _, ok := rs.Rooms[req.Room]; !ok {
		room, _ := rs.CreateRoom(req.Room)
		room.info.King = user.Name
	}

	cli, err = rs.Enter(cli, user, req.Room)
	if err != nil {
		cli.EmitError(err.Error(), "enter")
	}

	return cli
}

func (rs *RoomStore) actionsHandler(cli ws.Client, message []byte) ws.Client {
	room := cli.Room.(*Room)

	if room.GetInfo().King != cli.User.Name {
		cli.EmitError("not enough permission", "action")
		return cli
	}

	typ, err := jsonparser.GetString(message, "action")
	if err != nil {
		cli.EmitError(err.Error(), "artion")
		return cli
	}

	switch typ {
	case "color":
		req := struct {
			Color    string `json:"color"`
			Username string `json:"to"`
		}{}

		err = json.Unmarshal(message, &req)
		if err != nil {
			cli.EmitError(err.Error(), "color")
			return cli
		}

		err = room.ChangeColor(req.Color, req.Username)
		if err != nil {
			cli.EmitError(err.Error(), "color")
		}

	case "kick":
		req := struct {
			Target string `json:"target"`
		}{}

		err = json.Unmarshal(message, &req)
		if err != nil {
			cli.EmitError(err.Error(), "kick")
			return cli
		}

		err = room.Kick(req.Target)
		if err != nil {
			cli.EmitError(err.Error(), "kick")
		}

	case "king":
		req := struct {
			Target string `json:"target"`
		}{}

		err = json.Unmarshal(message, &req)
		if err != nil {
			cli.EmitError(err.Error(), "king")
			return cli
		}

		room.ChangeKing(req.Target)
		if err != nil {
			cli.EmitError(err.Error(), "king")
		}

	case "gamestart":
		_, err := room.StartGame()
		if err != nil {
			cli.EmitError(err.Error(), "start")
		}

	default:
		cli.EmitError("no such action", "action")
	}

	return cli
}
