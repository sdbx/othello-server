package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/util"
)

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if room, ok := RoomHub[vars["room"]]; !ok {
		w.WriteHeader(http.StatusNotFound)
	} else {
		room.serveWs(w, r)
	}
}

type roomCreateRequest struct {
	Secret string `json:"secret"`
}

func roomCreateHandler(w http.ResponseWriter, r *http.Request) {
	if !util.JsonTest(w, r) {
		return
	}
	vars := mux.Vars(r)
	room := vars["room"]
	request := roomCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		util.ErrorWrite(w, r, err.Error(), "RoomCreateHandler")
		return
	}
	room2, err := RoomHub.AddRoom(room, request.Secret)
	if err != nil {
		util.ErrorWrite(w, r, err.Error(), "RoomCreateHandler")
		return
	}
	go room2.run()
	w.WriteHeader(http.StatusOK)
}
