package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (h *Hub) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data struct {
		username string
	}

	decoder.Decode(&data)

	id := "Room_" + fmt.Sprintf("%d", time.Now().Unix())

	room := NewRoom(id, data.username)

	log.Printf("Created room: %v+\n", room)

	h.rooms[room] = true

	go room.Run()

	w.WriteHeader(http.StatusCreated)
}

func (h *Hub) WebsocketHandler(w http.ResponseWriter, r *http.Request) {

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

	room, err := h.getRoomById(roomId[0])

	if err != nil {
		fmt.Println("Here")
		log.Printf("Room with id %s not found", roomId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

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
}
