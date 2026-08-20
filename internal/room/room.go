package room

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	"whisperchat/internal/domain"
)

type joinRequest struct {
	client *domain.Client
	result chan error
}

var ErrRoomFull error = errors.New("room is full")

type IncomingMessage struct {
	Client  *domain.Client
	Content string
}

type Room struct {
	ID         string
	Clients    map[*domain.Client]bool
	Broadcast  chan *IncomingMessage
	Register   chan *joinRequest
	Unregister chan *domain.Client
	Limit      int
	SaveChan   chan<- *domain.Message
	Done       chan struct{}
}

func NewRoom(id string, limit int, saveChan chan<- *domain.Message) *Room {
	return &Room{
		ID:         id,
		Clients:    make(map[*domain.Client]bool),
		Broadcast:  make(chan *IncomingMessage),
		Register:   make(chan *joinRequest),
		Unregister: make(chan *domain.Client),
		Limit:      limit,
		SaveChan:   saveChan,
		Done:       make(chan struct{}),
	}
}

func (r *Room) Join(client *domain.Client) error {
	req := &joinRequest{
		client: client,
		result: make(chan error, 1),
	}

	r.Register <- req
	return <-req.result
}

func (r *Room) Run() {
	for {
		select {
		case req := <-r.Register:
			if len(r.Clients) >= r.Limit {
				req.result <- ErrRoomFull
				continue
			}
			r.Clients[req.client] = true
			req.result <- nil
			r.broadcastRaw(domain.NewJoinMessage(req.client.DisplayName))

		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
				r.broadcastRaw(domain.NewLeaveMessage(client.DisplayName))
			}

		case in := <-r.Broadcast:
			r.broadcastMessage(in)
		case <-r.Done:
			for client := range r.Clients {
				close(client.Send)
				delete(r.Clients, client)
			}
			return
		}
	}
}

func (r *Room) broadcastMessage(in *IncomingMessage) {
	now := time.Now()
	out := domain.WMessage{
		DisplayName: in.Client.DisplayName,
		Content:     in.Content,
		Timestamp:   now.Format(time.RFC3339),
	}

	payload, err := json.Marshal(out)
	if err != nil {
		log.Println("room %s: marshall outgoing message: %v", r.ID, err)
		return
	}
	r.broadcastRaw(payload)

	if r.SaveChan != nil {
		msg := &domain.Message{
			RoomID:      r.ID,
			DisplayName: in.Client.DisplayName,
			Content:     in.Content,
			CreatedAt:   now,
		}

		select {
		case r.SaveChan <- msg:
		default:
			log.Printf("room %s: save channel full, dropping message from history", r.ID)
		}
	}
}

func (r *Room) broadcastRaw(payload []byte) {
	for client := range r.Clients {
		select {
		case client.Send <- payload:
		default:
			close(client.Send)
			delete(r.Clients, client)
		}
	}
}
