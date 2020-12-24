package main

import (
	"encoding/json"
	"log"
)

const (
	SeekAction      = "SEEK"
	PauseAction     = "PAUSE"
	PlayAction      = "PLAY"
	LeaveAction     = "LEAVE"
	AddSongAction   = "ADD_SONG"
	JoinRoomAction  = "JOIN_ROOM"
	LeaveRoomAction = "LEAVE_ROOM"
)

type Message struct {
	Action string            `json:"action"`
	Args   map[string]string `json:"args"`
	Target *Room             `json:"target"`
	Sender *Client           `json:"sender"`
}

func (m *Message) Encode() []byte {
	encoded, err := json.Marshal(m)
	if err != nil {
		log.Println("Error in encoding message")
	}
	return encoded
}
