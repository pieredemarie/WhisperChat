package model

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID         string
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Limit      int
	mu         sync.RWMutex
}

type Client struct {
	Conn        *websocket.Conn
	Send        chan []byte
	DisplayName string
	RoomID      string
}

type Message struct {
	ID          int       `json:"id"`
	RoomID      string    `json:"room_id"`
	DisplayName string    `json:"display_name"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

type WMessage struct {
	DisplayName string `json:"display_name"`
	Content     string `json:"content"`
	Timestamp   string `json:"timestamp,omitempty"`
}

type SystemMessage struct {
	// Type = Join || Leave
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func NewJoinMessage(name string) []byte {
	msg := SystemMessage{
		Type:    "join",
		Message: name + "has joined the chat",
		Time:    time.Now().Format(time.RFC3339),
	}

	data, _ := json.Marshal(msg)

	return data
}

func NewLeaveMessage(name string) []byte {
	msg := SystemMessage{
		Type:    "leave",
		Message: name + "has left the chat",
		Time:    time.Now().Format(time.RFC3339),
	}

	data, _ := json.Marshal(msg)

	return data
}
