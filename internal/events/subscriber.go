package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/niaga-platform/service-customer/internal/models"
	"github.com/niaga-platform/service-customer/internal/repository"
	"go.uber.org/zap"
)

// HI-001: Back-in-Stock Event Subscriber

// ProductRestockedEvent represents a product back in stock event from inventory service
type ProductRestockedEvent struct {
	ProductID   string  `json:"product_id"`
	VariantID   string  `json:"variant_id,omitempty"`
	WarehouseID string  `json:"warehouse_id"`
	Quantity    float64 `json:"quantity"`
	ProductName string  `json:"product_name,omitempty"`
	ProductSlug string  `json:"product_slug,omitempty"`
}

// BackInStockSubscriber handles back-in-stock event subscriptions
type BackInStockSubscriber struct {
	nc                 *nats.Conn
	backInStockRepo    *repository.BackInStockRepository
	notificationClient NotificationClient
	logger             *zap.Logger
}

// NotificationClient interface for sending notifications
type NotificationClient interface {
	SendBackInStockNotification(notification models.BackInStockNotification) error
}

// NewBackInStockSubscriber creates a new subscriber
func NewBackInStockSubscriber(
	nc *nats.Conn,
	backInStockRepo *repository.BackInStockRepository,
	notificationClient NotificationClient,
	logger *zap.Logger,
) *BackInStockSubscriber {
	return &BackInStockSubscriber{
		nc:                 nc,
		backInStockRepo:    backInStockRepo,
		notificationClient: notificationClient,
		logger:             logger,
	}
}

// Subscribe starts listening for restock events
func (s *BackInStockSubscriber) Subscribe() error {
	_, err := s.nc.Subscribe("inventory.product.restocked", func(msg *nats.Msg) {
		s.handleRestockedEvent(msg.Data)
	})
	if err != nil {
		s.logger.Error("Failed to subscribe to inventory.product.restocked", zap.Error(err))
		return err
	}

	s.logger.Info("Subscribed to inventory.product.restocked events")
	return nil
}

// handleRestockedEvent processes a product restocked event
func (s *BackInStockSubscriber) handleRestockedEvent(data []byte) {
	var event ProductRestockedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		s.logger.Error("Failed to unmarshal restocked event", zap.Error(err))
		return
	}

	s.logger.Info("Processing product restocked event",
		zap.String("product_id", event.ProductID),
		zap.String("variant_id", event.VariantID),
		zap.Float64("quantity", event.Quantity))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse product ID
	productID, err := uuid.Parse(event.ProductID)
	if err != nil {
		s.logger.Error("Invalid product ID in event", zap.Error(err))
		return
	}

	// Parse variant ID if present
	var variantID *uuid.UUID
	if event.VariantID != "" {
		vid, err := uuid.Parse(event.VariantID)
		if err != nil {
			s.logger.Error("Invalid variant ID in event", zap.Error(err))
			return
		}
		variantID = &vid
	}

	// Get all pending subscriptions for this product/variant
	subscriptions, err := s.backInStockRepo.GetByProduct(ctx, productID, variantID)
	if err != nil {
		s.logger.Error("Failed to get subscriptions for product",
			zap.String("product_id", event.ProductID),
			zap.Error(err))
		return
	}

	if len(subscriptions) == 0 {
		s.logger.Debug("No pending subscriptions for restocked product",
			zap.String("product_id", event.ProductID))
		return
	}

	s.logger.Info("Found subscriptions to notify",
		zap.String("product_id", event.ProductID),
		zap.Int("count", len(subscriptions)))

	// Send notifications and mark as notified
	var notifiedIDs []uuid.UUID
	for _, sub := range subscriptions {
		// Build notification
		notification := models.BackInStockNotification{
			SubscriptionID: sub.ID.String(),
			CustomerID:     sub.CustomerID.String(),
			ProductID:      sub.ProductID.String(),
			ProductName:    sub.ProductName,
			ProductSlug:    sub.ProductSlug,
			ProductImage:   sub.ProductImage,
			StockQuantity:  int(event.Quantity),
		}

		if sub.VariantID != nil {
			notification.VariantID = sub.VariantID.String()
		}
		notification.VariantSKU = sub.VariantSKU
		notification.VariantName = sub.VariantName

		// Get customer info if available
		if sub.Customer != nil {
			notification.CustomerEmail = sub.Customer.Email
			notification.CustomerName = sub.Customer.FirstName + " " + sub.Customer.LastName
		}

		// Send notification
		if s.notificationClient != nil {
			if err := s.notificationClient.SendBackInStockNotification(notification); err != nil {
				s.logger.Error("Failed to send notification",
					zap.String("subscription_id", sub.ID.String()),
					zap.Error(err))
				continue
			}
		}

		notifiedIDs = append(notifiedIDs, sub.ID)
	}

	// Mark subscriptions as notified in batch
	if len(notifiedIDs) > 0 {
		if err := s.backInStockRepo.MarkMultipleAsNotified(ctx, notifiedIDs); err != nil {
			s.logger.Error("Failed to mark subscriptions as notified", zap.Error(err))
		} else {
			s.logger.Info("Marked subscriptions as notified",
				zap.Int("count", len(notifiedIDs)))
		}
	}
}

// SimpleNotificationClient is a basic HTTP client for notifications
type SimpleNotificationClient struct {
	baseURL string
	logger  *zap.Logger
}

// NewSimpleNotificationClient creates a new notification client
func NewSimpleNotificationClient(baseURL string, logger *zap.Logger) *SimpleNotificationClient {
	return &SimpleNotificationClient{
		baseURL: baseURL,
		logger:  logger,
	}
}

// SendBackInStockNotification sends a back-in-stock notification
func (c *SimpleNotificationClient) SendBackInStockNotification(notification models.BackInStockNotification) error {
	// In a real implementation, this would make an HTTP call to the notification service
	// For now, we'll log the notification
	c.logger.Info("Sending back-in-stock notification",
		zap.String("customer_email", notification.CustomerEmail),
		zap.String("product_name", notification.ProductName),
		zap.Int("stock_quantity", notification.StockQuantity))

	// TODO: Implement actual HTTP call to notification service
	// POST to c.baseURL + "/api/v1/notifications/back-in-stock"

	return nil
}
