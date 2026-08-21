package service

import (
	"whisperchat/internal/domain"
	"whisperchat/internal/room"
)

type ChatService struct {
	manager *room.Manager
}

func NewChatService(m *room.Manager) *ChatService {
	return &ChatService{
		manager: m,
	}
}

func (s *ChatService) JoinRoom(roomID string, client *domain.Client) (*room.Room, error) {
	r, err := s.manager.GetRoom(roomID)
	if err != nil {
		return nil, err
	}
	if err := r.Join(client); err != nil {
		return nil, err
	}
	return r, nil
}
