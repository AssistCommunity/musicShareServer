package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Room struct {
	Id        string           `json:"room_id"`
	Host      string           `json:"room_host"`
	Clients   map[*Client]bool `json:"-"`
	Queue     *Queue           `json:"room_queue"`
	join      chan *Client
	leave     chan *Client
	broadcast chan *Message
}

func NewRoom(id string, host string) *Room {
	return &Room{
		Id:        id,
		Host:      host,
		Clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan *Message),
		Queue:     NewQueue(),
	}
}

func (r *Room) EncodeInfo() []byte {
	type EncodingRoom struct {
		Room
		Clients []string
	}
	encoded, err := json.Marshal(r)

	if err != nil {
		log.Printf("Could not encode room info: %q\n", err)
		return []byte("{\"error\": \"Cannot encode room info\"}")
	}

	return encoded
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.joinClientInRoom(client)
		case client := <-r.leave:
			r.leaveClientFromRoom(client)
		case message := <-r.broadcast:
			fmt.Println("New message in room", r.Id)
			r.handleMessage(message)
		}
	}
}

func (r *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action: JoinRoomAction,
		Target: r,
		Sender: client,
	}
	go r.handleMessage(message)
}

func (r *Room) notifyClientLeft(client *Client) {
	message := &Message{
		Action: LeaveRoomAction,
		Target: r,
		Sender: client,
	}
	go r.handleMessage(message)
}

func (r *Room) joinClientInRoom(client *Client) {
	r.Clients[client] = true
	r.notifyClientJoined(client)
}

func (r *Room) leaveClientFromRoom(client *Client) {
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
	}
}

func (r *Room) handleMessage(m *Message) {
	log.Printf("New message: %v\n", m)

	// Side effects
	switch m.Action {
	case InfoReqAction:
		// Send room info
		resp := NewMessage(InfoRespAction)
		resp.SetArgs("info", r)
		m.Sender.send <- resp.Encode()

	case AddTrackAction:
		if val, ok := m.Args["track_name"]; ok {
			r.Queue.AddToQueue(val.(string))
		} else {
			log.Printf("Invalid add track message recieved\n")
		}

	case NextTrackAction:
		r.Queue.NextTrack()

	case PrevTrackAction:
		r.Queue.PrevTrack()

	case JumpTrackAction:
		if val, ok := m.Args["track_index"]; ok {
			r.Queue.JumpToIndex(val.(int))
		}
	}

	fmt.Println("Messages are being handled. ")

	bytes := m.Encode()
	for client := range r.Clients {
		client.send <- bytes
	}
}
