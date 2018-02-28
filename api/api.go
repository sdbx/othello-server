package api

import (
	"github.com/gorilla/mux"
	"github.com/sdbx/othello-server/othello"
)

var service *othello.Service

type h map[string]interface{}

func Start(serv *othello.Service) *mux.Router {
	service = serv
	r := mux.NewRouter()
	r.Path("/register").
		HandlerFunc(registerHandler).
		Methods("POST").
		Name("register")
	return r
}
