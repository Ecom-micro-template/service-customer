package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WishlistItemModel is the GORM persistence model for WishlistItem.
type WishlistItemModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`

	// Variant-specific fields
	VariantID   *uuid.UUID `gorm:"type:uuid" json:"variant_id,omitempty"`
	VariantSKU  *string    `gorm:"type:varchar(50)" json:"variant_sku,omitempty"`
	VariantName *string    `gorm:"type:varchar(100)" json:"variant_name,omitempty"`

	// Price tracking for price drop alerts
	PriceAtAdd   float64 `gorm:"type:decimal(10,2);default:0" json:"price_at_add"`
	NotifyOnSale bool    `gorm:"default:false" json:"notify_on_sale"`

	// Denormalized product info for display
	ProductName  *string `gorm:"type:varchar(255)" json:"product_name,omitempty"`
	ProductSlug  *string `gorm:"type:varchar(255)" json:"product_slug,omitempty"`
	ProductImage *string `gorm:"type:varchar(500)" json:"product_image,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name.
func (WishlistItemModel) TableName() string {
	return "customer.wishlist_items"
}

// BeforeCreate hook to generate UUID if not provided.
func (m *WishlistItemModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// GetUniqueKey returns a unique key for deduplication.
func (m *WishlistItemModel) GetUniqueKey() string {
	if m.VariantID != nil {
		return m.ProductID.String() + "-" + m.VariantID.String()
	}
	return m.ProductID.String() + "-nil"
}
