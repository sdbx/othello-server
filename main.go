package main

import (
	"net/http"

	"github.com/sdbx/othello-server/api"
	"github.com/sdbx/othello-server/othello"
)

func main() {
	service := othello.NewService()
	app := api.New(service)
	r := app.GetRouter()
	http.ListenAndServe("127.0.0.1:8080", r)
}
