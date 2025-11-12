package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// In production, check against allowed origins
		return true
	},
}

type WebSocketHandler struct {
	*BaseHandler
	clients   map[*websocket.Conn]string // conn -> userID
	broadcast chan []byte
	mu        sync.RWMutex
}

func NewWebSocketHandler(baseHandler *BaseHandler) *WebSocketHandler {
	handler := &WebSocketHandler{
		BaseHandler: baseHandler,
		clients:     make(map[*websocket.Conn]string),
		broadcast:   make(chan []byte, 256),
	}

	// Start broadcast goroutine
	go handler.handleBroadcasts()

	return handler
}

func (h *WebSocketHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/ws", h.handleWebSocket)
}

func (h *WebSocketHandler) handleWebSocket(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := h.GetUserID(c)
	if userID == "" {
		log.Println("[WebSocket] No userID - unauthorized connection attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] Failed to upgrade connection: %v\n", err)
		return
	}

	// Register client
	h.mu.Lock()
	h.clients[conn] = userID
	h.mu.Unlock()

	log.Printf("[WebSocket] ✅ Client connected: userID=%s, total=%d\n", userID, len(h.clients))

	// Handle client disconnect
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		clientCount := len(h.clients)
		h.mu.Unlock()

		conn.Close()
		log.Printf("[WebSocket] 🔌 Client disconnected: userID=%s, total=%d\n", userID, clientCount)
	}()

	// Configure connection
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Read messages from client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] Error reading message: %v\n", err)
			}
			break
		}

		// Parse message
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[WebSocket] Failed to parse message: %v\n", err)
			continue
		}

		// Handle ping
		if event, ok := msg["event"].(string); ok && event == "ping" {
			h.sendToClient(conn, map[string]interface{}{
				"event": "pong",
				"data":  map[string]interface{}{},
			})
			continue
		}

		log.Printf("[WebSocket] 📨 Received from %s: %v\n", userID, msg)
	}
}

// Send message to specific client
func (h *WebSocketHandler) sendToClient(conn *websocket.Conn, message map[string]interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteMessage(websocket.TextMessage, data)
}

// Broadcast message to all clients
func (h *WebSocketHandler) BroadcastMessage(event string, data interface{}) {
	message := map[string]interface{}{
		"event": event,
		"data":  data,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("[WebSocket] Failed to marshal broadcast message: %v\n", err)
		return
	}

	h.broadcast <- msgBytes
}

// Handle broadcast messages
func (h *WebSocketHandler) handleBroadcasts() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-h.broadcast:
			h.mu.RLock()
			for conn := range h.clients {
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("[WebSocket] Failed to send broadcast: %v\n", err)
					conn.Close()
				}
			}
			h.mu.RUnlock()

		case <-ticker.C:
			// Send ping to all clients
			h.mu.RLock()
			for conn := range h.clients {
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("[WebSocket] Failed to send ping: %v\n", err)
					conn.Close()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Send message to specific user
func (h *WebSocketHandler) SendToUser(userID string, event string, data interface{}) {
	message := map[string]interface{}{
		"event": event,
		"data":  data,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("[WebSocket] Failed to marshal message: %v\n", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for conn, connUserID := range h.clients {
		if connUserID == userID {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				log.Printf("[WebSocket] Failed to send to user %s: %v\n", userID, err)
			}
			break
		}
	}
}

// Send message to session participants
func (h *WebSocketHandler) SendToSession(sessionID string, event string, data interface{}) {
	// This would require session participant tracking
	// For now, broadcast to all clients
	h.BroadcastMessage(event, data)
}
