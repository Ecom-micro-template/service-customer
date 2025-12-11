package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Customer represents a customer in the system
type Customer struct {
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

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (Customer) TableName() string {
	return "public.customers"
}

// CreateCustomerRequest represents a request to create a customer
type CreateCustomerRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone,omitempty"`
}

// UpdateCustomerRequest represents a request to update a customer
type UpdateCustomerRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Status    *string `json:"status,omitempty"`
}

// CustomerNote represents a note on a customer
type CustomerNote struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	CustomerID uuid.UUID  `gorm:"type:uuid;index" json:"customer_id"`
	Note       string     `gorm:"type:text" json:"note"`
	IsPrivate  bool       `gorm:"default:false" json:"is_private"`
	CreatedBy  *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (n *CustomerNote) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

func (CustomerNote) TableName() string {
	return "public.customer_notes"
}

// CustomerActivity represents a customer activity log
type CustomerActivity struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CustomerID uuid.UUID `gorm:"type:uuid;index" json:"customer_id"`
	Type       string    `gorm:"type:varchar(50)" json:"type"`
	Title      string    `gorm:"type:varchar(255)" json:"title"`
	Details    string    `gorm:"type:text" json:"details,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *CustomerActivity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (CustomerActivity) TableName() string {
	return "public.customer_activities"
}

// CustomerSegment represents a customer segment
type CustomerSegment struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Color       string    `gorm:"type:varchar(7)" json:"color,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *CustomerSegment) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (CustomerSegment) TableName() string {
	return "public.customer_segments"
}

// CustomerSegmentAssignment represents assignment of a customer to a segment
type CustomerSegmentAssignment struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CustomerID uuid.UUID `gorm:"type:uuid;index" json:"customer_id"`
	SegmentID  uuid.UUID `gorm:"type:uuid;index" json:"segment_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (a *CustomerSegmentAssignment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (CustomerSegmentAssignment) TableName() string {
	return "public.customer_segment_assignments"
}

// CustomerListFilter represents filters for customer listing
type CustomerListFilter struct {
	Status    string     `form:"status"`
	Segment   string     `form:"segment"`
	DateFrom  *time.Time `form:"date_from"`
	DateTo    *time.Time `form:"date_to"`
	OrdersMin *int       `form:"orders_min"`
	OrdersMax *int       `form:"orders_max"`
	SpentMin  *float64   `form:"spent_min"`
	SpentMax  *float64   `form:"spent_max"`
	Search    string     `form:"search"`
	Page      int        `form:"page"`
	Limit     int        `form:"limit"`
	SortBy    string     `form:"sort_by"`
	SortOrder string     `form:"sort_order"`
}
