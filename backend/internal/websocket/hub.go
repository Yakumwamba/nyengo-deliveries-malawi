package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Client represents a WebSocket connection
type Client struct {
	ID            string
	CourierID     uuid.UUID
	Conn          *websocket.Conn
	Send          chan []byte
	Subscriptions map[string]bool // order IDs subscribed to
	mu            sync.RWMutex
}

// Hub maintains active WebSocket connections and manages subscriptions
type Hub struct {
	// Registered clients
	clients map[string]*Client

	// Clients subscribed to specific orders (for tracking)
	orderSubscriptions map[string]map[string]*Client // orderID -> clientID -> client

	// Channels
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	// Redis for cross-instance pub/sub
	redis *redis.Client

	mu sync.RWMutex
}

// Message types
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type LocationPayload struct {
	OrderID   string  `json:"orderId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed,omitempty"`
	Heading   float64 `json:"heading,omitempty"`
	Timestamp int64   `json:"timestamp"`
}

type SubscribePayload struct {
	OrderID string `json:"orderId"`
	Action  string `json:"action"` // "subscribe" or "unsubscribe"
}

type TrackingUpdate struct {
	OrderID           string    `json:"orderId"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	Speed             float64   `json:"speed,omitempty"`
	Heading           float64   `json:"heading,omitempty"`
	DistanceRemaining float64   `json:"distanceRemaining"`
	ETAMinutes        int       `json:"etaMinutes"`
	Timestamp         time.Time `json:"timestamp"`
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:            make(map[string]*Client),
		orderSubscriptions: make(map[string]map[string]*Client),
		broadcast:          make(chan []byte),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
	}
}

// SetRedis sets the Redis client for cross-instance communication
func (h *Hub) SetRedis(redis *redis.Client) {
	h.redis = redis

	if redis != nil {
		// Subscribe to tracking events from Redis
		go h.subscribeToRedis()
	}
}

// Run starts the hub event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				// Remove from all order subscriptions
				for orderID := range client.Subscriptions {
					if clients, exists := h.orderSubscriptions[orderID]; exists {
						delete(clients, client.ID)
						if len(clients) == 0 {
							delete(h.orderSubscriptions, orderID)
						}
					}
				}
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SubscribeToOrder subscribes a client to order tracking updates
func (h *Hub) SubscribeToOrder(clientID, orderID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, exists := h.clients[clientID]
	if !exists {
		return
	}

	// Add to order subscriptions
	if h.orderSubscriptions[orderID] == nil {
		h.orderSubscriptions[orderID] = make(map[string]*Client)
	}
	h.orderSubscriptions[orderID][clientID] = client

	// Track in client's subscriptions
	client.mu.Lock()
	client.Subscriptions[orderID] = true
	client.mu.Unlock()
}

// UnsubscribeFromOrder removes a client from order tracking
func (h *Hub) UnsubscribeFromOrder(clientID, orderID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, exists := h.orderSubscriptions[orderID]; exists {
		if client, ok := clients[clientID]; ok {
			client.mu.Lock()
			delete(client.Subscriptions, orderID)
			client.mu.Unlock()
		}
		delete(clients, clientID)
		if len(clients) == 0 {
			delete(h.orderSubscriptions, orderID)
		}
	}
}

// BroadcastToOrder sends a message to all clients subscribed to an order
func (h *Hub) BroadcastToOrder(orderID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, exists := h.orderSubscriptions[orderID]; exists {
		for _, client := range clients {
			select {
			case client.Send <- message:
			default:
				// Client buffer full, skip
			}
		}
	}

	// Also publish to Redis for other instances
	if h.redis != nil {
		h.redis.Publish(context.Background(), "tracking:broadcast:"+orderID, message)
	}
}

// BroadcastLocationUpdate sends a location update to all subscribed clients
func (h *Hub) BroadcastLocationUpdate(update *TrackingUpdate) {
	msg := WSMessage{
		Type: "location_update",
	}

	payload, _ := json.Marshal(update)
	msg.Payload = payload

	message, _ := json.Marshal(msg)
	h.BroadcastToOrder(update.OrderID, message)
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(clientID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[clientID]; ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, clientID)
		}
	}
}

// SendToCourier sends a message to a courier by their ID
func (h *Hub) SendToCourier(courierID uuid.UUID, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.CourierID == courierID {
			select {
			case client.Send <- message:
			default:
			}
		}
	}
}

// subscribeToRedis listens for tracking events from other server instances
func (h *Hub) subscribeToRedis() {
	ctx := context.Background()
	pubsub := h.redis.PSubscribe(ctx, "tracking:broadcast:*")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		// Extract order ID from channel name
		orderID := msg.Channel[len("tracking:broadcast:"):]

		h.mu.RLock()
		if clients, exists := h.orderSubscriptions[orderID]; exists {
			for _, client := range clients {
				select {
				case client.Send <- []byte(msg.Payload):
				default:
				}
			}
		}
		h.mu.RUnlock()
	}
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		courierID := c.Locals("courier_id").(uuid.UUID)

		client := &Client{
			ID:            uuid.New().String(),
			CourierID:     courierID,
			Conn:          c,
			Send:          make(chan []byte, 256),
			Subscriptions: make(map[string]bool),
		}

		hub.register <- client

		defer func() {
			hub.unregister <- client
			c.Close()
		}()

		// Writer goroutine
		go func() {
			for message := range client.Send {
				if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
					break
				}
			}
		}()

		// Reader loop
		for {
			_, rawMessage, err := c.ReadMessage()
			if err != nil {
				break
			}

			var msg WSMessage
			if err := json.Unmarshal(rawMessage, &msg); err != nil {
				continue
			}

			// Handle different message types
			switch msg.Type {
			case "location_update":
				// Driver sending location - broadcast to subscribers
				var loc LocationPayload
				if json.Unmarshal(msg.Payload, &loc) == nil {
					update := &TrackingUpdate{
						OrderID:   loc.OrderID,
						Latitude:  loc.Latitude,
						Longitude: loc.Longitude,
						Speed:     loc.Speed,
						Heading:   loc.Heading,
						Timestamp: time.Now(),
					}
					hub.BroadcastLocationUpdate(update)
				}

			case "subscribe":
				// Client subscribing to order updates
				var sub SubscribePayload
				if json.Unmarshal(msg.Payload, &sub) == nil {
					if sub.Action == "subscribe" {
						hub.SubscribeToOrder(client.ID, sub.OrderID)
					} else {
						hub.UnsubscribeFromOrder(client.ID, sub.OrderID)
					}
				}

			case "ping":
				// Respond with pong
				pong, _ := json.Marshal(WSMessage{Type: "pong"})
				client.Send <- pong
			}
		}
	})
}

// GetActiveSubscribers returns count of clients watching an order
func (h *Hub) GetActiveSubscribers(orderID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, exists := h.orderSubscriptions[orderID]; exists {
		return len(clients)
	}
	return 0
}
