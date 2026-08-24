package handler

import "net/http"

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /rooms", h.CreateRoom)
	mux.HandleFunc("GET /rooms/{id}/history", h.GetHistory)
	mux.HandleFunc("GET /rooms/{id}/ws", h.ServeWS)
	return mux
}
