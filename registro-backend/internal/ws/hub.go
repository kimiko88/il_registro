package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	SchoolID string
	Role     string
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients map[UserID]map[*Client]bool (allows multiple devices per user)
	clients map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Lock for map access (though channels handle most sync, map reads might need it if extended)
	mu sync.RWMutex
}

// Message defines the structure of WebSocket messages
type Message struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Recipient string      `json:"recipient,omitempty"` // UserID
	SchoolID  string      `json:"school_id,omitempty"` // Broadcast to school
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
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
			var deadClients []*Client

			h.mu.RLock()
			// Direct message to specific user
			if message.Recipient != "" {
				if clients, ok := h.clients[message.Recipient]; ok {
					bytes, _ := json.Marshal(message)
					for client := range clients {
						select {
						case client.Send <- bytes:
						default:
							deadClients = append(deadClients, client)
						}
					}
				}
			} else if message.SchoolID != "" {
				// Broadcast to all in school (naive implementation, could be optimized)
				bytes, _ := json.Marshal(message)
				for _, clients := range h.clients {
					for client := range clients {
						if client.SchoolID == message.SchoolID {
							select {
							case client.Send <- bytes:
							default:
								deadClients = append(deadClients, client)
							}
						}
					}
				}
			}
			h.mu.RUnlock()

			if len(deadClients) > 0 {
				h.mu.Lock()
				for _, client := range deadClients {
					if userClients, ok := h.clients[client.UserID]; ok {
						if _, exists := userClients[client]; exists {
							delete(userClients, client)
							close(client.Send)
							if len(userClients) == 0 {
								delete(h.clients, client.UserID)
							}
						}
					}
				}
				h.mu.Unlock()
			}
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
