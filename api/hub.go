package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello"
	"github.com/sdbx/othello-server/util"
)

func hubWebsocketHandler(h *othello.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if hub, ok := h.Rooms[vars["room"]]; !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			hub.ServeWs(w, r)
		}
	}
}

type hubCreateRequest struct {
	Secret string `json:"secret"`
}

func hubCreateHandler(h *othello.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !util.JsonTest(w, r) {
			return
		}
		vars := mux.Vars(r)
		hub := vars["hub"]
		request := hubCreateRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			util.ErrorWrite(w, r, err.Error(), "hubCreateHandler")
			return
		}
		hub2, err := h.AddRoom(hub, request.Secret)
		if err != nil {
			util.ErrorWrite(w, r, err.Error(), "hubCreateHandler")
			return
		}
		go hub2.Run()
		w.WriteHeader(http.StatusOK)
	}
}
