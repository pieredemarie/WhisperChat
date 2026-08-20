package handler

import "whisperchat/internal/service"

type Handler struct {
	service *service.ChatService
}

func NewHandler(s *service.ChatService) *Handler {
	return &Handler{
		service: s,
	}
}


