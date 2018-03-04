package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/sdbx/othello-server/api"
	"github.com/sdbx/othello-server/othello"
)

func main() {
	service := othello.NewService()
	r := api.Start(service)
	http.ListenAndServe("127.0.0.1:8080", r)
}
