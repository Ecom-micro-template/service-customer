package persistence

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MeasurementModel is the GORM persistence model for CustomerMeasurement.
type MeasurementModel struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name   *string   `gorm:"type:varchar(100)" json:"name,omitempty"`
	Gender string    `gorm:"type:varchar(20);not null" json:"gender"`

	// Upper body measurements (cm)
	Bust          *float64 `gorm:"type:decimal(5,1)" json:"bust,omitempty"`
	Chest         *float64 `gorm:"type:decimal(5,1)" json:"chest,omitempty"`
	Waist         *float64 `gorm:"type:decimal(5,1)" json:"waist,omitempty"`
	Hip           *float64 `gorm:"type:decimal(5,1)" json:"hip,omitempty"`
	ShoulderWidth *float64 `gorm:"type:decimal(5,1)" json:"shoulder_width,omitempty"`
	ArmLength     *float64 `gorm:"type:decimal(5,1)" json:"arm_length,omitempty"`

	// Lower body measurements (cm)
	Inseam  *float64 `gorm:"type:decimal(5,1)" json:"inseam,omitempty"`
	Outseam *float64 `gorm:"type:decimal(5,1)" json:"outseam,omitempty"`
	Thigh   *float64 `gorm:"type:decimal(5,1)" json:"thigh,omitempty"`

	// Additional measurements (cm/kg)
	Neck   *float64 `gorm:"type:decimal(5,1)" json:"neck,omitempty"`
	Wrist  *float64 `gorm:"type:decimal(5,1)" json:"wrist,omitempty"`
	Height *float64 `gorm:"type:decimal(5,1)" json:"height,omitempty"`
	Weight *float64 `gorm:"type:decimal(5,1)" json:"weight,omitempty"`

	Notes     *string   `gorm:"type:text" json:"notes,omitempty"`
	IsDefault bool      `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name.
func (MeasurementModel) TableName() string {
	return "crm.customer_measurements"
}

// BeforeCreate hook to generate UUID if not provided.
func (m *MeasurementModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
