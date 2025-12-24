package services

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type NotificationService struct {
	redis *redis.Client
}

func NewNotificationService(redis *redis.Client) *NotificationService {
	return &NotificationService{redis: redis}
}

type Notification struct {
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *NotificationService) Send(ctx context.Context, channel string, notification *Notification) error {
	if s.redis == nil {
		return nil
	}
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return s.redis.Publish(ctx, channel, data).Err()
}

func (s *NotificationService) SendOrderUpdate(ctx context.Context, courierID, orderID, status string) error {
	return s.Send(ctx, "courier:"+courierID, &Notification{
		Type: "order_update", Title: "Order Updated",
		Message: "Order status changed to " + status,
		Data:    map[string]string{"orderId": orderID, "status": status},
	})
}

func (s *NotificationService) SendNewOrder(ctx context.Context, courierID, orderID, customerName string) error {
	return s.Send(ctx, "courier:"+courierID, &Notification{
		Type: "new_order", Title: "New Order",
		Message: "New order from " + customerName,
		Data:    map[string]string{"orderId": orderID},
	})
}
