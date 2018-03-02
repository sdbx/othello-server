package main

import (
	"fmt"

	"github.com/olebedev/emitter"
)

func main() {
	e := &emitter.Emitter{}
	go func() {
		<-e.Emit("change", 42) // wait for the event sent successfully
		<-e.Emit("change", 37)
	}()

	for event := range e.On("change") {
		// do something with event.Args
		fmt.Println(event.Int(0)) // cast the first argument to int
	}
}
