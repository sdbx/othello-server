package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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

func gameCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := service.GameStore.CreateGame(vars["game"])
	if err != nil {
		errorWrite(w, r, err.Error(), "gameCreateHandler")
		return
	}
	w.WriteHeader(http.StatusOK)
}
