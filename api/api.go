package api

import (
	"encoding/json"
	"fmt"
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
	r.Path("/register").
		HandlerFunc(registerHandler).
		Methods("POST").
		Name("register")

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

func errorWrite(w http.ResponseWriter, r *http.Request, err string, from string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(Error{
		Msg:  err,
		From: from,
	})
	fmt.Println(err)
}
