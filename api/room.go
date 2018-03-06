package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello/room"
)

type brieftRoom struct {
	Name         string `json:"name"`
	King         string `json:"king"`
	Participants uint   `json:"n_of_people"`
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	list := []brieftRoom{}
	for _, t := range service.RoomStore.Rooms {
		rom := t.(*room.Room)
		list = append(list, brieftRoom{
			Name:         rom.Name(),
			King:         rom.King,
			Participants: rom.Participants,
		})
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h{
		"rooms": list,
	})
}

type detailRoom struct {
	Name         string   `json:"name"`
	King         string   `json:"king"`
	Participants []string `json:"participants"`
	Black        string   `json:"black"`
	White        string   `json:"white"`
	State        string   `json:"state"`
	Game         string   `json:"game"`
}

func roomsDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, ok := service.RoomStore.Rooms[vars["room"]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	rom := t.(*room.Room)
	resp := detailRoom{
		Name:         rom.Name(),
		King:         rom.King,
		Participants: rom.GetClientNames(),
		Black:        rom.Black,
		White:        rom.White,
		State:        rom.State.String(),
		Game:         rom.Game,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
