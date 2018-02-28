package room

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello"
	"github.com/sdbx/othello-server/util"
)

var service *othello.Service

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if room, ok := RoomHub[vars["room"]]; !ok {
		w.WriteHeader(http.StatusNotFound)
	} else {
		room.serveWs(w, r)
	}
}

type roomCreateRequest struct {
	Roomname string `json:"roomname`
	Secret   string `json:"secret"`
}

func roomCreateHandler(w http.ResponseWriter, r *http.Request) {
	if !util.JsonTest(w, r) {
		return
	}
	request := roomCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		util.ErrorWrite(w, r, err.Error(), "RoomCreateHandler")
	}
	RoomHub.AddRoom(request.Roomname, request.Secret)
}

func Start(serv *othello.Service) *mux.Router {
	service = serv
	r := mux.NewRouter()
	r.HandleFunc("/ws/{room}", handleWebsocket)

	r.HandleFunc("/{room}", roomCreateHandler).
		Methods("POST")

	r.HandleFunc("/", roomCreateHandler).
		Methods("GET")
	return r
}
