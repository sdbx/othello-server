package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello"
)

var service *othello.Service

var Logger io.Writer

type h map[string]interface{}

func Start(serv *othello.Service) http.Handler {
	service = serv
	r := mux.NewRouter()

	r.Path("/games/{game}").
		HandlerFunc(gameGetHandler).
		Methods("GET").
		Name("game get")

	r.Path("/games/{game}").
		HandlerFunc(gameCreateHandler).
		Methods("POST").
		Name("game create")

	r.Path("/games/{game}/actions").
		HandlerFunc(gameActionsHandler).
		Methods("POST").
		Name("game actions")

	r.Path("/ws/games").
		Handler(service.GameStore.WS.Handler()).
		Name("game websocket")

	r.Path("/ws/rooms").
		Handler(service.RoomStore.WS.Handler()).
		Name("room websocket")

	r.Path("/rooms").
		HandlerFunc(roomsHandler).
		Methods("GET").
		Name("room list")

	r.Path("/rooms/{room}").
		HandlerFunc(roomDetailHandler).
		Methods("GET").
		Name("room detail")

	r.Path("/users/{name}").
		HandlerFunc(userHandler).
		Methods("GET").
		Name("personal info")

	r.Path("/users/{name}/battles").
		HandlerFunc(battleHandler).
		Methods("GET").
		Name("battles")

	r.Path("/auth/{provider}/callback").
		HandlerFunc(authCallbackHandler).
		Methods("GET").
		Name("auth callback")

	r.Path("/auth/{provider}").
		HandlerFunc(authHandler).
		Methods("GET").
		Name("auth")

	cors3 := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "X-User-Secret"})
	cors2 := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	cors := handlers.AllowedOrigins([]string{"*"})
	corsed := handlers.CORS(cors, cors2, cors3)(r)
	return handlers.LoggingHandler(os.Stdout, corsed)
}

type Error struct {
	Msg  string `json:"msg"`
	From string `json:"from"`
}

func jsonTest(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		json.NewEncoder(w).Encode(Error{
			Msg:  "Content-Type is not application/json",
			From: "RegisterHandler",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}
