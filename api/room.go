package api

import (
	"encoding/json"
	"net/http"
)

type brieftRoom struct {
	Name         string `json:"name"`
	King         string `json:"king"`
	Participants uint   `json:"n_of_people"`
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	infos := service.RoomStore.GetInfos()
	json.NewEncoder(w).Encode(infos)
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
