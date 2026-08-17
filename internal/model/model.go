package model

import (
	"sync"

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
