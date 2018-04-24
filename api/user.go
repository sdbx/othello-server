package api

import (
	"encoding/json"
	"net/http"

	"github.com/sdbx/othello-server/othello/dbs"
)

type registerRequest struct {
	Username string `json:"username"`
}
type userInfo struct {
	Name   string `json:"username"`
	Secret string `json:"secret"`
}

func usersMeHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-User-Secret")

	user, err := dbs.GetUserBySecret(secret)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp := userInfo{
		Name:   user.Name,
		Secret: user.Secret,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
