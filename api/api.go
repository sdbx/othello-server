package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello"
)

type API struct {
	service *othello.Service
}

type Error struct {
	Msg  string `json:"msg"`
	From string `json:"from"`
}

type h map[string]interface{}

func New(service *othello.Service) *API {
	return &API{
		service: service,
	}
}

func jsonTest(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{
			Msg:  "Content-Type is not application/json",
			From: "RegisterHandler",
		})
		return false
	}
	return true
}

func errorWrite(w http.ResponseWriter, r *http.Request, err string, from string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Error{
		Msg:  err,
		From: from,
	})
}

func (ap *API) GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path("/register").
		HandlerFunc(ap.RegisterHandler).
		Methods("POST").
		Name("register")

	r.Path("/rooms").
		HandlerFunc(ap.RoomsHandler).
		Methods("GET").
		Name("rooms")

	r.Path("/rooms/{room}").
		HandlerFunc(ap.RoomsCreateHandler).
		Methods("POST").
		Name("room create")

	return r
}
