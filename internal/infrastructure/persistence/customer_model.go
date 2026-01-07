// Package persistence contains GORM models and repository implementations.
package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomerModel is the GORM persistence model for Customer.
type CustomerModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	FirstName   string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName    string         `gorm:"type:varchar(100)" json:"last_name"`
	Phone       string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	AvatarURL   string         `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	TotalOrders int            `gorm:"default:0" json:"total_orders"`
	TotalSpent  float64        `gorm:"type:decimal(12,2);default:0" json:"total_spent"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name.
func (CustomerModel) TableName() string {
	return "public.customers"
}

// BeforeCreate hook to generate UUID if not provided.
func (m *CustomerModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// CustomerNoteModel is the GORM model for CustomerNote.
type CustomerNoteModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CustomerID uuid.UUID  `gorm:"type:uuid;index" json:"customer_id"`
	Note       string     `gorm:"type:text" json:"note"`
	IsPrivate  bool       `gorm:"default:false" json:"is_private"`
	CreatedBy  *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TableName specifies the table name.
func (CustomerNoteModel) TableName() string {
	return "public.customer_notes"
}

// BeforeCreate hook.
func (m *CustomerNoteModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// CustomerActivityModel is the GORM model for CustomerActivity.
type CustomerActivityModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CustomerID uuid.UUID `gorm:"type:uuid;index" json:"customer_id"`
	Type       string    `gorm:"type:varchar(50)" json:"type"`
	Title      string    `gorm:"type:varchar(255)" json:"title"`
	Details    string    `gorm:"type:text" json:"details,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the table name.
func (CustomerActivityModel) TableName() string {
	return "public.customer_activities"
}

// BeforeCreate hook.
func (m *CustomerActivityModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// CustomerSegmentModel is the GORM model for CustomerSegment.
type CustomerSegmentModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Color       string    `gorm:"type:varchar(7)" json:"color,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name.
func (CustomerSegmentModel) TableName() string {
	return "public.customer_segments"
}

// BeforeCreate hook.
func (m *CustomerSegmentModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
