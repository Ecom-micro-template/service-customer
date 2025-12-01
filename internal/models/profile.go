package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Profile represents a customer profile
type Profile struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	FullName       string     `gorm:"type:varchar(200)" json:"full_name"`
	Email          string     `gorm:"type:varchar(200);uniqueIndex" json:"email"`
	Phone          string     `gorm:"type:varchar(50)" json:"phone"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	Gender         string     `gorm:"type:varchar(20)" json:"gender,omitempty"` // male, female, other
	ProfilePicture string     `gorm:"type:varchar(500)" json:"profile_picture,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// TableName specifies the table name for Profile
func (Profile) TableName() string {
	return "customer.profiles"
}

// BeforeCreate hook to ensure UUID is set
func (p *Profile) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
