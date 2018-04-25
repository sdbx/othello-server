package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello/dbs"
)

type registerRequest struct {
	Username string `json:"username"`
}
type userInfo struct {
	Name    string `json:"username"`
	Profile string `json:"profile"`
	WinLose string `json:"win_lose"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	user, err := dbs.GetUserByName(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := userInfo{
		Name:    user.Name,
		Profile: user.Profile,
		WinLose: user.GetWinLose(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func battleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	user, err := dbs.GetUserByName(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("p"))
	battles := user.GetBattles(page)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(battles)
}
