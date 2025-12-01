package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WishlistItem represents a product saved to customer's wishlist
type WishlistItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
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
