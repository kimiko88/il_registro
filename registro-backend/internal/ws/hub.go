// internal/ws/hub.go
// Hub distribuito via Redis Pub/Sub — compatibile con deploy multi-istanza.
// Ogni istanza pubblica messaggi su Redis e li riceve tramite subscribe.
// Le connessioni WebSocket rimangono locali all'istanza, ma i messaggi
// vengono instradati correttamente indipendentemente da quale istanza li origina.

package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const redisPubSubChannel = "ws:broadcast"

// Client rappresenta una singola connessione WebSocket su questa istanza.
type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	UserID    string
	SchoolID  string
	Role      string
	closeOnce sync.Once
}

func (c *Client) CloseSend() {
	c.closeOnce.Do(func() { close(c.Send) })
}

// Message è il formato condiviso sia per i canali Go interni
// sia per la serializzazione JSON su Redis.
type Message struct {
	Type         string      `json:"type"`
	Payload      interface{} `json:"payload"`
	Recipient    string      `json:"recipient,omitempty"`
	SchoolID     string      `json:"school_id,omitempty"`
	AllowedRoles []string    `json:"allowed_roles,omitempty"`
}

// Hub gestisce le connessioni locali e un bridge Redis per il multi-istanza.
type Hub struct {
	// clients: userID -> set di connessioni locali
	clients map[string]map[*Client]bool
	mu      sync.RWMutex

	// canali interni (solo per messaggi originati da questa istanza)
	localBroadcast chan Message
	register       chan *Client
	unregister     chan *Client

	// Redis (opzionale: se redisURL non è fornito, funziona in modalità in-memory locale)
	rdb    *redis.Client
	pubsub *redis.PubSub
}

// NewHub crea l'Hub e stabilisce la connessione Redis se redisURL non è vuoto.
// redisURL esempio: "redis://localhost:6379" oppure l'URL di Upstash.
func NewHub(redisURL string) *Hub {
	h := &Hub{
		clients:        make(map[string]map[*Client]bool),
		localBroadcast: make(chan Message, 512),
		register:       make(chan *Client, 64),
		unregister:     make(chan *Client, 64),
	}

	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Printf("ws.Hub: invalid REDIS_URL: %v. Falling back to local in-memory hub.", err)
		} else {
			h.rdb = redis.NewClient(opt)
			log.Println("ws.Hub: initialized with Redis Pub/Sub support")
		}
	} else {
		log.Println("ws.Hub: REDIS_URL not provided, running in local in-memory mode")
	}

	return h
}

// Run avvia il loop principale + il subscriber Redis (se abilitato).
// Da chiamare con go hub.Run(ctx) all'avvio del server.
func (h *Hub) Run(ctx context.Context) {
	if h.rdb != nil {
		// Subscribe al canale Redis condiviso
		h.pubsub = h.rdb.Subscribe(ctx, redisPubSubChannel)
		// Goroutine: legge da Redis e delivera ai client locali
		go h.redisListener(ctx)
	}

	// Loop principale: gestisce register/unregister + messaggi locali
	for {
		select {
		case <-ctx.Done():
			if h.pubsub != nil {
				_ = h.pubsub.Close()
			}
			if h.rdb != nil {
				_ = h.rdb.Close()
			}
			return

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
				if _, ok2 := userClients[client]; ok2 {
					delete(userClients, client)
					client.CloseSend()
					if len(userClients) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.localBroadcast:
			if h.rdb != nil {
				// Pubblica su Redis -> raggiunge TUTTE le istanze (inclusa questa)
				h.publishToRedis(ctx, msg)
			} else {
				// Modalità locale in-memory senza Redis
				h.deliverLocally(msg)
			}
		}
	}
}

// redisListener riceve tutti i messaggi dal canale Redis e li consegna
// ai client WebSocket connessi a questa istanza.
func (h *Hub) redisListener(ctx context.Context) {
	ch := h.pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case redisMsg, ok := <-ch:
			if !ok {
				return
			}
			var msg Message
			if err := json.Unmarshal([]byte(redisMsg.Payload), &msg); err != nil {
				log.Printf("ws.Hub: failed to unmarshal redis message: %v", err)
				continue
			}
			h.deliverLocally(msg)
		}
	}
}

// deliverLocally consegna il messaggio ai client WebSocket di questa istanza.
func (h *Hub) deliverLocally(msg Message) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if msg.Recipient != "" {
		// Unicast: solo al destinatario specifico
		clients, ok := h.clients[msg.Recipient]
		if !ok {
			return
		}
		for client := range clients {
			if isRoleAllowed(client.Role, msg.AllowedRoles) {
				h.sendToClient(client, clients, bytes)
			}
		}
		if len(clients) == 0 {
			delete(h.clients, msg.Recipient)
		}
	} else if msg.SchoolID != "" {
		// Broadcast per scuola (con filtro opzionale per ruolo)
		for userID, clients := range h.clients {
			for client := range clients {
				if client.SchoolID == msg.SchoolID && isRoleAllowed(client.Role, msg.AllowedRoles) {
					h.sendToClient(client, clients, bytes)
				}
			}
			if len(clients) == 0 {
				delete(h.clients, userID)
			}
		}
	}
}

// sendToClient invia bytes al client; se il canale è pieno, chiude il client.
func (h *Hub) sendToClient(c *Client, siblings map[*Client]bool, bytes []byte) {
	select {
	case c.Send <- bytes:
	default:
		c.CloseSend()
		delete(siblings, c)
	}
}

// publishToRedis serializza il messaggio e lo pubblica su Redis.
func (h *Hub) publishToRedis(ctx context.Context, msg Message) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ws.Hub: marshal error: %v", err)
		return
	}
	if err := h.rdb.Publish(ctx, redisPubSubChannel, bytes).Err(); err != nil {
		log.Printf("ws.Hub: redis publish error: %v", err)
	}
}

// --- API pubblica ---

func (h *Hub) BroadcastToUser(userID, msgType string, payload interface{}) {
	h.localBroadcast <- Message{Type: msgType, Payload: payload, Recipient: userID}
}

func (h *Hub) BroadcastToSchool(schoolID, msgType string, payload interface{}) {
	h.localBroadcast <- Message{Type: msgType, Payload: payload, SchoolID: schoolID}
}

func (h *Hub) BroadcastToSchoolRoles(schoolID string, allowedRoles []string, msgType string, payload interface{}) {
	h.localBroadcast <- Message{Type: msgType, Payload: payload, SchoolID: schoolID, AllowedRoles: allowedRoles}
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
