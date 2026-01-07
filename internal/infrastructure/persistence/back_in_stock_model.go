package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BackInStockSubscriptionModel is the GORM persistence model for back-in-stock subscriptions.
type BackInStockSubscriptionModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CustomerID uuid.UUID  `gorm:"type:uuid;not null;index:idx_bis_customer" json:"customerId"`
	ProductID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_bis_product" json:"productId"`
	VariantID  *uuid.UUID `gorm:"type:uuid;index:idx_bis_variant" json:"variantId,omitempty"`

	// Denormalized product info for quick access
	ProductName  string `gorm:"size:255" json:"productName"`
	ProductSlug  string `gorm:"size:255" json:"productSlug"`
	ProductImage string `gorm:"size:500" json:"productImage,omitempty"`
	VariantSKU   string `gorm:"size:100" json:"variantSku,omitempty"`
	VariantName  string `gorm:"size:255" json:"variantName,omitempty"`

	// Notification tracking
	IsNotified         bool       `gorm:"default:false" json:"isNotified"`
	NotificationSentAt *time.Time `json:"notificationSentAt,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Customer *CustomerModel `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

// TableName specifies the table name.
func (BackInStockSubscriptionModel) TableName() string {
	return "customer.back_in_stock_subscriptions"
}

// BeforeCreate hook to generate UUID if not provided.
func (m *BackInStockSubscriptionModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
