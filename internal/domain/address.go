// Package domain contains GORM persistence models for the customer service.
//
// Deprecated: This package is being migrated to DDD architecture.
// For new development, use:
//   - Domain models: github.com/Ecom-micro-template/service-customer/internal/domain/address
//   - Persistence: github.com/Ecom-micro-template/service-customer/internal/infrastructure/persistence
//
// Existing code can continue using this package during the transition period.
package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Address represents a customer shipping/billing address
type Address struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Label         string    `gorm:"type:varchar(50)" json:"label"` // Home, Office, Other
	RecipientName string    `gorm:"type:varchar(200);not null" json:"recipient_name"`
	Phone         string    `gorm:"type:varchar(50);not null" json:"phone"`
	AddressLine1  string    `gorm:"type:varchar(500);not null" json:"address_line1"`
	AddressLine2  string    `gorm:"type:varchar(500)" json:"address_line2,omitempty"`
	City          string    `gorm:"type:varchar(100);not null" json:"city"`
	State         string    `gorm:"type:varchar(100);not null" json:"state"`
	Postcode      string    `gorm:"type:varchar(20);not null" json:"postcode"`
	Country       string    `gorm:"type:varchar(100);not null;default:'USA'" json:"country"`
	IsDefault     bool      `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName specifies the table name for Address
func (Address) TableName() string {
	return "customer.addresses"
}

// BeforeCreate hook to ensure UUID is set
func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
