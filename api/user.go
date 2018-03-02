package api

import (
	"encoding/json"
	"net/http"
)

type registerRequest struct {
	Username string `json:"username"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if !jsonTest(w, r) {
		return
	}

	request := registerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorWrite(w, r, err.Error(), "RegisterHandler")
		return
	}

	name := request.Username
	if service.UserStore.GetUserByID(name) != nil {
		errorWrite(w, r, "username already exist", "RegisterHandler")
		return
	}

	key := service.UserStore.Register(name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h{
		"secret": key,
	})
}
