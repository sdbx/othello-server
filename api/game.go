package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello/game"
)

type actionFunc func(w http.ResponseWriter, r *http.Request, gam *game.Game, bytes []byte)

var actions = map[string]actionFunc{
	"put": actionsPut,
}

type gameCreateRequest struct {
	Blackname string `json:"black"`
	Whitename string `json:"white"`
}

func gameGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	game := service.GameStore.GetGame(vars["game"])
	if game == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	game.RLock()
	resp := h{
		"board":   game.Board,
		"history": game.History,
		"initial": game.GameType.Initial(),
		"usernames": h{
			"black": game.Black,
			"white": game.White,
		},
		"times": h{
			"black": game.GetBlackTime(),
			"white": game.GetWhiteTime(),
		},
	}
	game.RUnlock()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func gameActionsHandler(w http.ResponseWriter, r *http.Request) {
	if !jsonTest(w, r) {
		return
	}

	vars := mux.Vars(r)
	game := service.GameStore.GetGame(vars["game"])
	if game == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	typ, err := jsonparser.GetString(bytes, "type")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	if action, ok := actions[typ]; ok {
		action(w, r, game, bytes)
	}
}
