package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second
	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

type Client struct {
	Username string `json:"username"`
	conn     *websocket.Conn
	hub      *Hub
	room     *Room
	send     chan []byte
}

func NewClient(username string, conn *websocket.Conn, hub *Hub, room *Room) *Client {
	return &Client{
		Username: username,
		conn:     conn,
		hub:      hub,
		room:     room,
		send:     make(chan []byte),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.disconnect()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.handleMessage(jsonMessage)
	}
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			fmt.Println("here")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) disconnect() {
	c.hub.unregister <- c
	close(c.send)
	c.conn.Close()
}

func (c *Client) handleMessage(jsonMessage []byte) {
	var message Message

	log.Printf("Message recieved from %s\n", c.Username)
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Invalid message %s", err)
		return
	}

	message.Sender = c

	log.Printf("%#v\n", message)

	c.room.broadcast <- &message
}
