package main

import (
	"encoding/json"
	"log"
)

const (
	InfoReqAction   = "INFO_REQ"
	InfoRespAction  = "INFO_RESP"
	SeekAction      = "SEEK"
	PauseAction     = "PAUSE"
	PlayAction      = "PLAY"
	AddTrackAction  = "ADD_TRACK"
	NextTrackAction = "NEXT_TRACK"
	PrevTrackAction = "PREV_TRACK"
	JumpTrackAction = "JUMP_TRACK"
	JoinRoomAction  = "JOIN_ROOM"
	LeaveRoomAction = "LEAVE_ROOM"
)

type Message struct {
	Action string                 `json:"action"`
	Args   map[string]interface{} `json:"args"`
	Target *Room                  `json:"target"`
	Sender *Client                `json:"sender"`
}

func NewMessage(action string) *Message {
	return &Message{
		Action: action,
		Args:   make(map[string]interface{}),
	}
}

func (m *Message) SetArgs(key string, value interface{}) {
	m.Args[key] = value
}

func (m *Message) Encode() []byte {
	encoded, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error in encoding message: %s\n", err)
	}
	return encoded
}
