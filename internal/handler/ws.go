package handler

import (
	"log"
	"net/http"
	"time"
	"whisperchat/internal/domain"
	"whisperchat/internal/room"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("id")
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error %v", err)
		return
	}

	cl := &domain.Client{
		Conn:        conn,
		Send:        make(chan []byte, 255),
		DisplayName: name,
		RoomID:      roomID,
	}

	rm, err := h.service.JoinRoom(roomID, cl)
	if err != nil {
		cl.Conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
		cl.Conn.Close()
		return
	}

	go writePump(cl)
	readPump(cl, rm)

}

func writePump(client *domain.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func readPump(client *domain.Client, r *room.Room) {
	defer func() {
		client.Conn.Close()
		r.Unregister <- client
	}()

	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close for %q: %v", client.DisplayName, err)
				return
			}
		}
		r.Broadcast <- &room.IncomingMessage{
			Client:  client,
			Content: string(msg),
		}
	}
}
