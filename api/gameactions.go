package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sdbx/othello-server/othello"
)

type actionsPutRequest struct {
	_      string `json:"type"`
	Secret string `json:"secret`
	Move   string `json:"move"`
}

func actionsPut(w http.ResponseWriter, r *http.Request, game *othello.Game, bytes []byte) {
	req := actionsPutRequest{}
	err := json.Unmarshal(bytes, &req)
	if err != nil {
		errorWrite(w, r, err.Error(), "actionsPut")
		return
	}
	user := service.UserStore.GetUserBySecret(req.Secret)
	if user == nil {
		errorWrite(w, r, "user doesn't exist", "actionsPut")
		return
	}
	cord, err := othello.CordFromMove(othello.Move(req.Move))
	fmt.Println(cord)
	if err != nil {
		errorWrite(w, r, err.Error(), "actionsPut")
		return
	}
	if game.Black == user.Name {
		err = game.Put(cord, othello.GameTileBlack)
	} else if game.White == user.Name {
		err = game.Put(cord, othello.GameTileWhite)
	} else {
		errorWrite(w, r, "you are not a player", "actionsPut")
		return
	}
	if err != nil {
		errorWrite(w, r, err.Error(), "actionsPut")
		return
	}
	w.WriteHeader(http.StatusOK)
}
