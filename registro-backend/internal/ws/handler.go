package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8192
)

var (
	allowedOriginsOnce   sync.Once
	cachedAllowedOrigins map[string]bool
)

// allowedOrigins returns the set of permitted WebSocket origins from the
// ALLOWED_ORIGINS environment variable (comma-separated). Falls back to
// rejecting all cross-origin requests if the variable is not set.
func allowedOrigins() map[string]bool {
	allowedOriginsOnce.Do(func() {
		raw := os.Getenv("ALLOWED_ORIGINS")
		cachedAllowedOrigins = make(map[string]bool)
		for _, o := range strings.Split(raw, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				cachedAllowedOrigins[o] = true
			}
		}
	})
	return cachedAllowedOrigins
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowed := allowedOrigins()
		if len(allowed) == 0 {
			// Default local dev fallback: allow only exact localhost/127.0.0.1 hosts.
			// We parse the origin URL and compare the host exactly — no substring match.
			if origin == "" {
				return true
			}
			parsed, err := url.Parse(origin)
			if err != nil {
				return false
			}
			host := parsed.Hostname() // strips port
			return host == "localhost" || host == "127.0.0.1"
		}
		if allowed["*"] {
			isProd := os.Getenv("GIN_MODE") == "release" || os.Getenv("APP_ENV") == "production"
			if isProd {
				log.Println("WARNING: Wildcard '*' in ALLOWED_ORIGINS is forbidden in production environment; rejecting WebSocket connection")
				return false
			}
			return true
		}
		if origin == "" {
			// Reject missing Origin header when restrictive allowed origins are configured
			return false
		}
		return allowed[origin]
	},
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws error: %v", err)
			}
			break
		}
		// Use strict JSON matching so that only {"type":"PING"} triggers a PONG,
		// not any payload that happens to contain the substring "PING".
		var cm wsClientMsg
		if err := json.Unmarshal(msg, &cm); err == nil && cm.Type == "PING" {
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			_ = c.Conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"PONG"}`))
		}
	}
}

// wsClientMsg is the minimal shape of messages sent by the browser client.
type wsClientMsg struct {
	Type string `json:"type"`
}

// writeSingleMessage sends data as its own WebSocket text frame.
// Using a dedicated NextWriter/Close pair per message guarantees that each
// JSON object is delivered as a separate frame, preventing concatenation bugs
// like {"type":"A"}{"type":"B"} that break client-side JSON.parse().
func (c *Client) writeSingleMessage(data []byte) error {
	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	if _, err = w.Write(data); err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.writeSingleMessage(message); err != nil {
				return
			}
			// Drain any additionally queued messages — each as its own frame.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				if err := c.writeSingleMessage(<-c.Send); err != nil {
					return
				}
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) Listen(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Safe extraction with ok-pattern to avoid panic on missing/wrong-typed values.
	role := ""
	if v, ok := c.Get("role"); ok {
		if s, ok := v.(string); ok {
			role = s
		}
	}
	schoolID := ""
	if v, ok := c.Get("school_id"); ok {
		if s, ok := v.(string); ok {
			schoolID = s
		}
	}
	if role != "superadmin" && role != "system_auditor" && schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required for websocket connection"})
		return
	}

	// The auth middleware handles token validation before this point.
	// We no longer negotiate Sec-WebSocket-Protocol for token delivery;
	// the token is passed via Bearer header or ?token= query param instead.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ws upgrade error:", err)
		return
	}

	client := &Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Role:     role,
		SchoolID: schoolID,
	}

	client.Hub.register <- client
	go client.writePump()
	go client.readPump()
}
