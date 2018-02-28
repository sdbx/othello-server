package util

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Msg  string `json:"msg"`
	From string `json:"from"`
}

func JsonTest(w http.ResponseWriter, r *http.Request) bool {
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

func ErrorWrite(w http.ResponseWriter, r *http.Request, err string, from string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(Error{
		Msg:  err,
		From: from,
	})
}
