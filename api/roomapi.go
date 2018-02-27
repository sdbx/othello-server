package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type RegisterRequest struct {
	Username string `json:"username"`
}

func (ap *API) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !jsonTest(w, r) {
		return
	}

	request := RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	encoder := json.NewEncoder(w)
	if err != nil {
		errorWrite(w, r, err.Error(), "RegisterHandler")
		return
	}

	name := request.Username
	if len(ap.service.UserStore.GetKey(name)) != 0 {
		errorWrite(w, r, "username already exist", "RegisterHandler")
		return
	}

	key := ap.service.UserStore.GenKey(name)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(h{
		"key": key,
	})
}

func (ap *API) RoomsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	list := []h{}
	for _, item := range ap.service.RoomStore {
		list = append(list, h{
			"name":  item.Name,
			"users": item.Users,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h{
		"rooms": list,
	})
}

type RoomsCreateRequest struct {
	Secret string `json:"secret"`
}

func (ap *API) RoomsCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !jsonTest(w, r) {
		return
	}
	vars := mux.Vars(r)
	room := vars["room"]

	if len(room) == 0 {
		errorWrite(w, r, "room name is empty", "RoomsCreateHandler")
		return
	}

	request := RoomsCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errorWrite(w, r, err.Error(), "RoomsCreateHandler")
		return
	}

	if user, ok := ap.service.UserStore[request.Secret]; !ok {
		errorWrite(w, r, "user with the secret doesn't exist", "RoomsCreateHandler")
	} else {
		w.WriteHeader(http.StatusOK)
		ap.service.RoomStore.AddRoom(user, room)
	}

}
