package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sdbx/othello-server/othello/dbs"
	"github.com/sdbx/othello-server/othello/game"
)

type actionsPutRequest struct {
	_    string `json:"type"`
	Move string `json:"move"`
}

func actionsPut(w http.ResponseWriter, r *http.Request, gam *game.Game, bytes []byte) {
	secret := r.Header.Get("X-User-Secret")
	req := actionsPutRequest{}
	err := json.Unmarshal(bytes, &req)
	if err != nil {
		errorWrite(w, r, err.Error(), "actionsPut")
		return
	}
	user, err := dbs.GetUserBySecret(secret)
	if err != nil {
		errorWrite(w, r, "user doesn't exist", "actionsPut")
		return
	}
	cord, err := game.CordFromMove(game.Move(req.Move))
	fmt.Println(cord)
	if err != nil {
		errorWrite(w, r, err.Error(), "actionsPut")
		return
	}
	if gam.Black == user.Name {
		err = gam.Put(cord, game.GameTileBlack)
	} else if gam.White == user.Name {
		err = gam.Put(cord, game.GameTileWhite)
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
