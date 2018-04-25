package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/naver"
	"github.com/sdbx/othello-server/api"
	"github.com/sdbx/othello-server/othello"
)

func main() {
	goth.UseProviders(
		naver.New(os.Getenv("NAVER_KEY"), os.Getenv("NAVER_SECRET"), os.Getenv("NAVER_CALLBACK")),
	)
	service := othello.NewService()
	r := api.Start(service)
	fmt.Println(os.Getenv("API_ADDR"))
	log.Println(http.ListenAndServe(os.Getenv("API_ADDR"), r))
}
