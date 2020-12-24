package main

import (
	"fmt"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[*Room]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			fmt.Printf("Message: %s\n", message)
			h.broadcastToClients(message)
		}
	}
}

func (h *Hub) registerClient(c *Client) {
	h.clients[c] = true
}

func (h *Hub) unregisterClient(c *Client) {
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
	}
}

func (h *Hub) broadcastToClients(message []byte) {
	for client := range h.clients {
		client.send <- message
	}
}

func (h *Hub) getRoomById(id string) *Room {
	var foundRoom *Room
	for room := range h.rooms {
		if id == room.id {
			foundRoom = room
			break
		}
	}

	if foundRoom == nil {
		foundRoom = NewRoom(id)
		h.rooms[foundRoom] = true
		go foundRoom.Run()
	}

	return foundRoom
}
