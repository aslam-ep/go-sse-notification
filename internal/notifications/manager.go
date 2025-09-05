package notifications

import (
	"sync"

	"github.com/aslam-ep/go-sse-notification/internal/metrics"
)

type Notification struct {
	ID      int64  `json:"id"`
	UserID  string `json:"userID"`
	Message string `json:"message"`
}

type Manager struct {
	Mu      sync.RWMutex
	clients map[string][]chan Notification
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string][]chan Notification),
	}
}

func (m *Manager) AddClient(userID string, ch chan Notification) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.clients[userID] = append(m.clients[userID], ch)
	metrics.ActiveClients.Inc()
}

func (m *Manager) RemoveClient(userID string, ch chan Notification) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if channels, ok := m.clients[userID]; ok {
		for i, c := range channels {
			if c == ch {
				m.clients[userID] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
	}

	// Cleanup if no client left for the user
	if len(m.clients[userID]) == 0 {
		delete(m.clients, userID)
	}
	metrics.ActiveClients.Dec()
}

func (m *Manager) Send(n Notification) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	if chans, ok := m.clients[n.UserID]; ok {
		for _, ch := range chans {
			select {
			case ch <- n:
				metrics.MessageCount.WithLabelValues(n.UserID).Inc()
			default: // Don't block if the client is slow
				metrics.DroppedMessages.Inc()
			}
		}
	}
}
