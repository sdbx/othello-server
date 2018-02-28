package api

import (
	"encoding/json"
	"net/http"

	"github.com/sdbx/othello-server/util"
)

type registerRequest struct {
	Username string `json:"username"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !util.JsonTest(w, r) {
		return
	}

	request := registerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		util.ErrorWrite(w, r, err.Error(), "RegisterHandler")
		return
	}

	name := request.Username
	if service.UserStore.GetUserByName(name) != nil {
		util.ErrorWrite(w, r, "username already exist", "RegisterHandler")
		return
	}

	key := service.UserStore.Register(name)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h{
		"key": key,
	})
}
