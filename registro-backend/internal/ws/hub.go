package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	SchoolID string
	Role     string
}

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	Type         string      `json:"type"`
	Payload      interface{} `json:"payload"`
	Recipient    string      `json:"recipient,omitempty"`
	SchoolID     string      `json:"school_id,omitempty"`
	AllowedRoles []string    `json:"allowed_roles,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func isRoleAllowed(role string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, r := range allowed {
		if r == role {
			return true
		}
	}
	return false
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; !ok {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if userClients, ok := h.clients[client.UserID]; ok {
				if _, ok := userClients[client]; ok {
					delete(userClients, client)
					close(client.Send)
					if len(userClients) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			// Use full Lock (not RLock) because we may delete stale clients from the map.
			h.mu.Lock()
			if message.Recipient != "" {
				if clients, ok := h.clients[message.Recipient]; ok {
					bytes, _ := json.Marshal(message)
					for client := range clients {
						if isRoleAllowed(client.Role, message.AllowedRoles) {
							select {
							case client.Send <- bytes:
							default:
								close(client.Send)
								delete(clients, client)
							}
						}
					}
				}
			} else if message.SchoolID != "" {
				bytes, _ := json.Marshal(message)
				for _, clients := range h.clients {
					for client := range clients {
						if client.SchoolID == message.SchoolID && isRoleAllowed(client.Role, message.AllowedRoles) {
							select {
							case client.Send <- bytes:
							default:
								close(client.Send)
								delete(clients, client)
							}
						}
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) BroadcastToUser(userID string, msgType string, payload interface{}) {
	h.broadcast <- Message{
		Type:      msgType,
		Payload:   payload,
		Recipient: userID,
	}
}

func (h *Hub) BroadcastToSchool(schoolID string, msgType string, payload interface{}) {
	h.broadcast <- Message{
		Type:     msgType,
		Payload:  payload,
		SchoolID: schoolID,
	}
}

func (h *Hub) BroadcastToSchoolRoles(schoolID string, allowedRoles []string, msgType string, payload interface{}) {
	h.broadcast <- Message{
		Type:         msgType,
		Payload:      payload,
		SchoolID:     schoolID,
		AllowedRoles: allowedRoles,
	}
}
