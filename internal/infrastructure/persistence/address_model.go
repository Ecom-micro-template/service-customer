package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddressModel is the GORM persistence model for Address.
type AddressModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Label         string    `gorm:"type:varchar(50)" json:"label"`
	RecipientName string    `gorm:"type:varchar(200);not null" json:"recipient_name"`
	Phone         string    `gorm:"type:varchar(50);not null" json:"phone"`
	AddressLine1  string    `gorm:"type:varchar(500);not null" json:"address_line1"`
	AddressLine2  string    `gorm:"type:varchar(500)" json:"address_line2,omitempty"`
	City          string    `gorm:"type:varchar(100);not null" json:"city"`
	State         string    `gorm:"type:varchar(100);not null" json:"state"`
	Postcode      string    `gorm:"type:varchar(20);not null" json:"postcode"`
	Country       string    `gorm:"type:varchar(100);not null;default:'Malaysia'" json:"country"`
	IsDefault     bool      `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName specifies the table name.
func (AddressModel) TableName() string {
	return "customer.addresses"
}

// BeforeCreate hook to generate UUID if not provided.
func (m *AddressModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
