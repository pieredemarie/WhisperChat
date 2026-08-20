package service

import "whisperchat/internal/domain"

type ChatRepository interface {
	Save(msg *domain.Message) error
	GetRecent(chatID string, limit int) (*[]domain.Message, error)
}

// TODO: Add methods here

type ChatService struct {
	repo ChatRepository
}

func NewChatService()
