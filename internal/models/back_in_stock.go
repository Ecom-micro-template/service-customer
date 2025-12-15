package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HI-001: Back-in-Stock Subscription Model

// BackInStockSubscription represents a customer's subscription to be notified
// when an out-of-stock product becomes available again
type BackInStockSubscription struct {
	ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CustomerID         uuid.UUID      `gorm:"type:uuid;not null;index:idx_bis_customer" json:"customerId"`
	ProductID          uuid.UUID      `gorm:"type:uuid;not null;index:idx_bis_product" json:"productId"`
	VariantID          *uuid.UUID     `gorm:"type:uuid;index:idx_bis_variant" json:"variantId,omitempty"`

	// Denormalized product info for quick access
	ProductName        string         `gorm:"size:255" json:"productName"`
	ProductSlug        string         `gorm:"size:255" json:"productSlug"`
	ProductImage       string         `gorm:"size:500" json:"productImage,omitempty"`
	VariantSKU         string         `gorm:"size:100" json:"variantSku,omitempty"`
	VariantName        string         `gorm:"size:255" json:"variantName,omitempty"`

	// Notification tracking
	IsNotified         bool           `gorm:"default:false" json:"isNotified"`
	NotificationSentAt *time.Time     `json:"notificationSentAt,omitempty"`

	// Timestamps
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Customer           *Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (BackInStockSubscription) TableName() string {
	return "customer.back_in_stock_subscriptions"
}

// BackInStockSubscribeInput is the request body for subscribing
type BackInStockSubscribeInput struct {
	ProductID    string `json:"productId" binding:"required"`
	VariantID    string `json:"variantId,omitempty"`
	ProductName  string `json:"productName"`
	ProductSlug  string `json:"productSlug"`
	ProductImage string `json:"productImage,omitempty"`
	VariantSKU   string `json:"variantSku,omitempty"`
	VariantName  string `json:"variantName,omitempty"`
}

// BackInStockStats represents statistics about back-in-stock subscriptions
type BackInStockStats struct {
	TotalSubscriptions   int64 `json:"totalSubscriptions"`
	PendingNotifications int64 `json:"pendingNotifications"`
	SentNotifications    int64 `json:"sentNotifications"`
	UniqueProducts       int64 `json:"uniqueProducts"`
	UniqueCustomers      int64 `json:"uniqueCustomers"`
}

// BackInStockNotification is the data sent to notification service
type BackInStockNotification struct {
	SubscriptionID string `json:"subscriptionId"`
	CustomerID     string `json:"customerId"`
	CustomerEmail  string `json:"customerEmail"`
	CustomerName   string `json:"customerName"`
	ProductID      string `json:"productId"`
	ProductName    string `json:"productName"`
	ProductSlug    string `json:"productSlug"`
	ProductImage   string `json:"productImage"`
	VariantID      string `json:"variantId,omitempty"`
	VariantSKU     string `json:"variantSku,omitempty"`
	VariantName    string `json:"variantName,omitempty"`
	StockQuantity  int    `json:"stockQuantity"`
}
