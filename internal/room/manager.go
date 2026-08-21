package room

import (
	"errors"
	"sync"
	"whisperchat/internal/domain"
)

var ErrRoomNotFound = errors.New("room not found")

type Manager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

func (m *Manager) CreateRoom(id string, limit int, saveChan chan<- *domain.Message) *Room {
	r := NewRoom(id, limit, saveChan)
	m.mu.Lock()
	m.rooms[id] = r
	m.mu.Unlock()
	go r.Run()
	return r
}

func (m *Manager) GetRoom(id string) (*Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rooms[id]
	if !ok {
		return nil, ErrRoomNotFound
	}

	return r, nil
}
