package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/api"
	"github.com/sdbx/othello-server/othello"
	"github.com/sdbx/othello-server/room"
)

func main() {
	service := othello.NewService()
	r := room.Start(service)
	r2 := api.Start(service)
	group := mux.NewRouter()
	group.Handle("/room", r)
	group.Handle("/api", r2)
	http.ListenAndServe("127.0.0.1:8080", group)
}
