package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello"
)

type actionFunc func(w http.ResponseWriter, r *http.Request, game *othello.Game, bytes []byte)

var actions = map[string]actionFunc{
	"put": actionsPut,
}

type BriefRoom struct {
	Name string `json:"name"`
}

func gameListHandler(w http.ResponseWriter, r *http.Request) {
	list := []BriefRoom{}
	for _, item := range service.GameStore.ListGames() {
		list = append(list, BriefRoom{
			Name: item,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h{
		"games": list,
	})
}

type gameCreateRequest struct {
	Blackname string `json:"black"`
	Whitename string `json:"white"`
}

func gameCreateHandler(w http.ResponseWriter, r *http.Request) {
	if !jsonTest(w, r) {
		return
	}

	req := gameCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorWrite(w, r, err.Error(), "gameCreateHandler")
		return
	}

	vars := mux.Vars(r)
	err = service.GameStore.CreateGame(vars["game"], req.Blackname, req.Whitename, othello.DefaultOthello{})
	if err != nil {
		errorWrite(w, r, err.Error(), "gameCreateHandler")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func gameGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	game := service.GameStore.GetGame(vars["game"])
	if game == nil {
		errorWrite(w, r, "game doesn't exist", "gameGetHandler")
		return
	}
	list := [8][8]uint{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			list[i][j] = uint(game.GetTile(othello.Coordinate{j, i}))
		}
	}
	resp := h{
		"black":   game.Black,
		"white":   game.White,
		"board":   game.Board,
		"history": game.History,
		"initial": game.GameType.Initial(),
		"list":    list,
	}
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
		errorWrite(w, r, "game doesn't exist", "gameActionHandler")
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorWrite(w, r, err.Error(), "gameActionHandler")
		return
	}
	typ, err := jsonparser.GetString(bytes, "type")
	if err != nil {
		errorWrite(w, r, err.Error(), "gameActionHandler")
		return
	}
	if action, ok := actions[typ]; ok {
		action(w, r, game, bytes)
	}
}
