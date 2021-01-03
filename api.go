package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *Hub) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	decoder := json.NewDecoder(r.Body)

	var data struct {
		Username string `json:"username"`
	}

	var resp struct {
		RoomId string `json:"room_id"`
	}

	decoder.Decode(&data)

	id := "Room_" + fmt.Sprintf("%d", time.Now().Unix())

	room := NewRoom(id, data.Username)

	log.Printf("Created room: %v+\n", room)

	h.rooms[room] = true

	go room.Run()

	w.WriteHeader(http.StatusCreated)

	resp.RoomId = room.Id

	respBody, err := json.Marshal(resp)

	if err != nil {
		log.Printf("Cannot write response: %s\n", err)
	}

	_, err = w.Write(respBody)
}

func (h *Hub) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

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
