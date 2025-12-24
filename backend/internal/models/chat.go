package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageType represents the type of chat message
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeLocation MessageType = "location"
	MessageTypeSystem   MessageType = "system"
)

// ChatConversation represents a chat thread between courier and customer
type ChatConversation struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	OrderID       uuid.UUID  `json:"orderId" db:"order_id"`
	CourierID     uuid.UUID  `json:"courierId" db:"courier_id"`
	CustomerPhone string     `json:"customerPhone" db:"customer_phone"`
	CustomerName  string     `json:"customerName" db:"customer_name"`
	IsActive      bool       `json:"isActive" db:"is_active"`
	LastMessageAt *time.Time `json:"lastMessageAt,omitempty" db:"last_message_at"`
	UnreadCount   int        `json:"unreadCount" db:"unread_count"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

// ChatMessage represents a single message in a conversation
type ChatMessage struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	ConversationID uuid.UUID   `json:"conversationId" db:"conversation_id"`
	SenderType     string      `json:"senderType" db:"sender_type"` // "courier" or "customer"
	SenderID       string      `json:"senderId" db:"sender_id"`
	MessageType    MessageType `json:"messageType" db:"message_type"`
	Content        string      `json:"content" db:"content"`
	MediaURL       string      `json:"mediaUrl,omitempty" db:"media_url"`
	Latitude       *float64    `json:"latitude,omitempty" db:"latitude"`
	Longitude      *float64    `json:"longitude,omitempty" db:"longitude"`
	IsRead         bool        `json:"isRead" db:"is_read"`
	ReadAt         *time.Time  `json:"readAt,omitempty" db:"read_at"`
	CreatedAt      time.Time   `json:"createdAt" db:"created_at"`
}

// SendMessageRequest is the request for sending a chat message
type SendMessageRequest struct {
	ConversationID uuid.UUID   `json:"conversationId" validate:"required"`
	MessageType    MessageType `json:"messageType" validate:"required,oneof=text image location"`
	Content        string      `json:"content" validate:"required_if=MessageType text"`
	MediaURL       string      `json:"mediaUrl,omitempty"`
	Latitude       *float64    `json:"latitude,omitempty"`
	Longitude      *float64    `json:"longitude,omitempty"`
}

// ConversationListItem is a simplified view for listing conversations
type ConversationListItem struct {
	ID            uuid.UUID  `json:"id"`
	OrderNumber   string     `json:"orderNumber"`
	CustomerName  string     `json:"customerName"`
	LastMessage   string     `json:"lastMessage"`
	LastMessageAt *time.Time `json:"lastMessageAt"`
	UnreadCount   int        `json:"unreadCount"`
	IsActive      bool       `json:"isActive"`
}

// QuickReply represents a predefined quick reply option
type QuickReply struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Message string `json:"message"`
}

// GetDefaultQuickReplies returns common quick reply options
func GetDefaultQuickReplies() []QuickReply {
	return []QuickReply{
		{ID: "otw", Label: "On my way", Message: "I'm on my way to pick up your package."},
		{ID: "arrived_pickup", Label: "Arrived at pickup", Message: "I've arrived at the pickup location."},
		{ID: "picked_up", Label: "Package picked up", Message: "I've picked up your package and heading to you now."},
		{ID: "nearby", Label: "Almost there", Message: "I'm almost at your location, please be ready."},
		{ID: "arrived_delivery", Label: "Arrived", Message: "I've arrived at the delivery location."},
		{ID: "cant_find", Label: "Can't find location", Message: "I'm having trouble finding the exact location. Could you share more details?"},
		{ID: "delay", Label: "Slight delay", Message: "I'm experiencing a slight delay due to traffic. I'll be there shortly."},
		{ID: "delivered", Label: "Delivered", Message: "Your package has been delivered. Thank you!"},
	}
}
