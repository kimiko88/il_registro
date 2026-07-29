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
	"time"

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
// Bug 89: se il canale Redis si chiude inaspettatamente (restart, network blip),
// tentiamo la riconnessione con backoff esponenziale invece di uscire silenziosamente.
func (h *Hub) redisListener(ctx context.Context) {
	backoff := 500 * time.Millisecond
	const maxBackoff = 30 * time.Second

	for {
		ch := h.pubsub.Channel()
		running := true
		for running {
			select {
			case <-ctx.Done():
				return
			case redisMsg, ok := <-ch:
				if !ok {
					// Bug 89: canale chiuso — prova a riconnettersi
					log.Printf("ws.Hub: redis pubsub channel closed, reconnecting in %v...", backoff)
					running = false
				} else {
					backoff = 500 * time.Millisecond // reset on success
					var msg Message
					if err := json.Unmarshal([]byte(redisMsg.Payload), &msg); err != nil {
						log.Printf("ws.Hub: failed to unmarshal redis message: %v", err)
						continue
					}
					h.deliverLocally(msg)
				}
			}
		}

		// Wait before reconnect
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}

		// Re-subscribe
		_ = h.pubsub.Close()
		h.pubsub = h.rdb.Subscribe(ctx, redisPubSubChannel)

		// Increase backoff (capped)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

// deliverLocally consegna il messaggio ai client WebSocket di questa istanza.
func (h *Hub) deliverLocally(msg Message) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	var targetClients []*Client
	h.mu.RLock()
	if msg.Recipient != "" {
		if clients, ok := h.clients[msg.Recipient]; ok {
			for client := range clients {
				// Bug 87: verifica SchoolID quando il messaggio ha un contesto scolastico,
				// per prevenire consegne cross-tenant se gli UserID non fossero globalmente unici.
				if msg.SchoolID != "" && client.SchoolID != msg.SchoolID {
					continue
				}
				if isRoleAllowed(client.Role, msg.AllowedRoles) {
					targetClients = append(targetClients, client)
				}
			}
		}
	} else if msg.SchoolID != "" {
		for _, clients := range h.clients {
			for client := range clients {
				if client.SchoolID == msg.SchoolID && isRoleAllowed(client.Role, msg.AllowedRoles) {
					targetClients = append(targetClients, client)
				}
			}
		}
	}
	h.mu.RUnlock()

	for _, client := range targetClients {
		select {
		case client.Send <- bytes:
		default:
			select {
			case h.unregister <- client:
			default:
			}
		}
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

func (h *Hub) sendBroadcast(msg Message) {
	select {
	case h.localBroadcast <- msg:
	default:
		// Bug 88: log dettagliato per facilitare il debugging in produzione.
		// NOTA: per garantire zero drop su messaggi critici (GRADE_ADDED, ABSENCE_RECORDED)
		// è necessario un message broker persistente (es. Redis Streams o DB inbox).
		log.Printf("[WARN] ws.Hub: localBroadcast channel full — message DROPPED (type=%s recipient=%q schoolID=%q)",
			msg.Type, msg.Recipient, msg.SchoolID)
	}
}

func (h *Hub) BroadcastToUser(userID, msgType string, payload interface{}) {
	h.sendBroadcast(Message{Type: msgType, Payload: payload, Recipient: userID})
}

func (h *Hub) BroadcastToSchool(schoolID, msgType string, payload interface{}) {
	h.sendBroadcast(Message{Type: msgType, Payload: payload, SchoolID: schoolID})
}

func (h *Hub) BroadcastToSchoolRoles(schoolID string, allowedRoles []string, msgType string, payload interface{}) {
	h.sendBroadcast(Message{Type: msgType, Payload: payload, SchoolID: schoolID, AllowedRoles: allowedRoles})
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
