package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	var h = NewHub()

	go h.Run()

	log.Fatal(http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query()

		username, ok := q["username"]
		if !ok {
			log.Println("Param username missing")
			return
		}

		roomId, ok := q["room_id"]

		if !ok {
			log.Println("Param username missing")
			return
		}

		room := h.getRoomById(roomId[0])

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := NewClient(username[0], conn, h, room)
		room.joinClientInRoom(client)

		fmt.Println("Got new client!")
		fmt.Printf("%#v\n", client)

		go client.readPump()
		go client.writePump()
	})))
}
