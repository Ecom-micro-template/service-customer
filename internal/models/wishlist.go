// Package models contains GORM persistence models for the customer service.
//
// Deprecated: This package is being migrated to DDD architecture.
// For new development, use:
//   - Domain models: github.com/niaga-platform/service-customer/internal/domain/wishlist
//   - Persistence: github.com/niaga-platform/service-customer/internal/infrastructure/persistence
//
// Existing code can continue using this package during the transition period.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WishlistItem represents a product (optionally with specific variant) saved to customer's wishlist
type WishlistItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`

	// Variant-specific fields (CUS-001)
	VariantID   *uuid.UUID `gorm:"type:uuid" json:"variant_id,omitempty"` // Optional: specific variant (null = any variant)
	VariantSKU  *string    `gorm:"type:varchar(50)" json:"variant_sku,omitempty"`
	VariantName *string    `gorm:"type:varchar(100)" json:"variant_name,omitempty"` // e.g., "Red / Large"

	// Price tracking for price drop alerts
	PriceAtAdd   float64 `gorm:"type:decimal(10,2);default:0" json:"price_at_add"`
	NotifyOnSale bool    `gorm:"default:false" json:"notify_on_sale"`

	// Denormalized product info for display without joining
	ProductName  *string `gorm:"type:varchar(255)" json:"product_name,omitempty"`
	ProductSlug  *string `gorm:"type:varchar(255)" json:"product_slug,omitempty"`
	ProductImage *string `gorm:"type:varchar(500)" json:"product_image,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for WishlistItem
func (WishlistItem) TableName() string {
	return "customer.wishlist_items"
}

// BeforeCreate hook to ensure UUID is set
func (w *WishlistItem) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// GetUniqueKey returns a unique key for deduplication (product + variant combination)
func (w *WishlistItem) GetUniqueKey() string {
	if w.VariantID != nil {
		return w.ProductID.String() + "-" + w.VariantID.String()
	}
	return w.ProductID.String() + "-nil"
}
