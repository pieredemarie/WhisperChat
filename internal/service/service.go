package service

import (
	"errors"
	"log"
	"whisperchat/internal/domain"
	"whisperchat/internal/room"

	"github.com/google/uuid"
)

var ErrInvalidLimit = errors.New("limit must be positive")

type ChatRepository interface {
	Save(msg *domain.Message) error
	GetRecent(roomID string, limit int) ([]*domain.Message, error)
}

type ChatService struct {
	manager  *room.Manager
	repo     ChatRepository
	saveChan chan *domain.Message
}

func NewChatService(m *room.Manager, repo ChatRepository) *ChatService {
	s := &ChatService{
		manager:  m,
		repo:     repo,
		saveChan: make(chan *domain.Message, 1000),
	}
	go s.persistWorker()
	return s
}

func (s *ChatService) persistWorker() {
	for msg := range s.saveChan {
		if err := s.repo.Save(msg); err != nil {
			log.Printf("message failed (room %s): %v", msg.RoomID, err)
		}
	}
}

func (s *ChatService) CreateRoom(limit int) *room.Room {
	id := uuid.NewString()
	return s.manager.CreateRoom(id, limit, s.saveChan)
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

func (s *ChatService) GetHistory(roomID string, limit int) ([]*domain.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetRecent(roomID, limit)
}
