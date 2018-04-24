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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	user, err := dbs.GetUserBySecret(secret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "no such user")
		return
	}

	cord, err := game.CordFromMove(game.Move(req.Move))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	if gam.Black == user.Name {
		err = gam.Put(cord, game.GameTileBlack)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	if gam.White == user.Name {
		err = gam.Put(cord, game.GameTileWhite)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "you are not a player")
}
