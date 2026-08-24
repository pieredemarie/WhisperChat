package handler

import (
	"encoding/json"
	"net/http"
)

type createRoomRequest struct {
	Limit int `json:"limit"`
}

type createRoomResponse struct {
	ID string `json:"id"`
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body request", http.StatusBadRequest)
		return
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	rm := h.service.CreateRoom(req.Limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createRoomResponse{ID: rm.ID})
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("id")

	messages, err := h.service.GetHistory(roomID, 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
