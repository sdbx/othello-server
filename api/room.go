package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello/room"
)

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	infos := service.RoomStore.GetInfos()
	json.NewEncoder(w).Encode(infos)
}

func roomDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ro := service.RoomStore.GetRoom(vars["room"])
	if ro == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no such room")
		return
	}

	info := ro.(*room.Room).GetInfo()
	json.NewEncoder(w).Encode(info)
	w.WriteHeader(http.StatusOK)
}
