package main

import "fmt"

type Room struct {
	id        string
	clients   map[*Client]bool
	join      chan *Client
	leave     chan *Client
	broadcast chan *Message
}

func NewRoom(id string) *Room {
	return &Room{
		id:        id,
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan *Message),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.joinClientInRoom(client)
		case client := <-r.leave:
			r.leaveClientFromRoom(client)
		case message := <-r.broadcast:
			fmt.Println("New message in room", r.id)
			r.broadcastMessageToClientsInRoom(message)
		}
	}
}

func (r *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action: LeaveAction,
		Target: r,
	}
	r.broadcastMessageToClientsInRoom(message)
}

func (r *Room) joinClientInRoom(client *Client) {
	r.notifyClientJoined(client)
	r.clients[client] = true
}

func (r *Room) leaveClientFromRoom(client *Client) {
	if _, ok := r.clients[client]; ok {
		delete(r.clients, client)
	}
}

func (r *Room) broadcastMessageToClientsInRoom(m *Message) {
	bytes := m.Encode()
	for client := range r.clients {
		client.send <- bytes
	}
}
