package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	var h = NewHub()

	go h.Run()

	r := mux.NewRouter()

	r.HandleFunc("/", h.WebsocketHandler)
	r.HandleFunc("/create", h.CreateRoomHandler)

	fmt.Println("Serving on :8080")
	http.ListenAndServe(":8080", r)
}
