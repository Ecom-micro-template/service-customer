package measurement

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain/shared"
)

// Domain errors for Measurement entity
var (
	ErrMeasurementNotFound = errors.New("measurement not found")
	ErrInvalidMeasurement  = errors.New("invalid measurement data")
)

// CustomerMeasurement is the entity for customer body measurements.
type CustomerMeasurement struct {
	id          uuid.UUID
	userID      uuid.UUID
	measurement shared.BodyMeasurement
	isDefault   bool
	createdAt   time.Time
	updatedAt   time.Time
}

// MeasurementParams contains parameters for creating a CustomerMeasurement.
type MeasurementParams struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Gender        string
	Name          string
	Bust          *float64
	Chest         *float64
	Waist         *float64
	Hip           *float64
	ShoulderWidth *float64
	ArmLength     *float64
	Inseam        *float64
	Outseam       *float64
	Thigh         *float64
	Neck          *float64
	Wrist         *float64
	Height        *float64
	Weight        *float64
	Notes         string
	IsDefault     bool
}

// NewCustomerMeasurement creates a new CustomerMeasurement entity.
func NewCustomerMeasurement(params MeasurementParams) (*CustomerMeasurement, error) {
	if params.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}

	gender := shared.GenderMen
	if params.Gender == "women" {
		gender = shared.GenderWomen
	}

	measurement := shared.NewBodyMeasurement(shared.BodyMeasurementParams{
		Gender:        gender,
		Name:          params.Name,
		Bust:          params.Bust,
		Chest:         params.Chest,
		Waist:         params.Waist,
		Hip:           params.Hip,
		ShoulderWidth: params.ShoulderWidth,
		ArmLength:     params.ArmLength,
		Inseam:        params.Inseam,
		Outseam:       params.Outseam,
		Thigh:         params.Thigh,
		Neck:          params.Neck,
		Wrist:         params.Wrist,
		Height:        params.Height,
		Weight:        params.Weight,
		Notes:         params.Notes,
	})

	id := params.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	now := time.Now()
	return &CustomerMeasurement{
		id:          id,
		userID:      params.UserID,
		measurement: measurement,
		isDefault:   params.IsDefault,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Getters
func (m *CustomerMeasurement) ID() uuid.UUID                       { return m.id }
func (m *CustomerMeasurement) UserID() uuid.UUID                   { return m.userID }
func (m *CustomerMeasurement) Measurement() shared.BodyMeasurement { return m.measurement }
func (m *CustomerMeasurement) IsDefault() bool                     { return m.isDefault }
func (m *CustomerMeasurement) CreatedAt() time.Time                { return m.createdAt }
func (m *CustomerMeasurement) UpdatedAt() time.Time                { return m.updatedAt }

// Convenience getters
func (m *CustomerMeasurement) Name() string          { return m.measurement.Name() }
func (m *CustomerMeasurement) Gender() shared.Gender { return m.measurement.Gender() }

// --- Behavior Methods ---

// Update updates the measurement.
func (m *CustomerMeasurement) Update(params MeasurementParams) {
	gender := shared.GenderMen
	if params.Gender == "women" {
		gender = shared.GenderWomen
	}

	m.measurement = shared.NewBodyMeasurement(shared.BodyMeasurementParams{
		Gender:        gender,
		Name:          params.Name,
		Bust:          params.Bust,
		Chest:         params.Chest,
		Waist:         params.Waist,
		Hip:           params.Hip,
		ShoulderWidth: params.ShoulderWidth,
		ArmLength:     params.ArmLength,
		Inseam:        params.Inseam,
		Outseam:       params.Outseam,
		Thigh:         params.Thigh,
		Neck:          params.Neck,
		Wrist:         params.Wrist,
		Height:        params.Height,
		Weight:        params.Weight,
		Notes:         params.Notes,
	})
	m.updatedAt = time.Now()
}

// SetDefault sets this measurement as the default.
func (m *CustomerMeasurement) SetDefault() {
	m.isDefault = true
	m.updatedAt = time.Now()
}

// ClearDefault clears the default flag.
func (m *CustomerMeasurement) ClearDefault() {
	m.isDefault = false
	m.updatedAt = time.Now()
}

// IsComplete returns true if the measurement has all standard measurements.
func (m *CustomerMeasurement) IsComplete() bool {
	return m.measurement.IsComplete()
}

// BMI returns the BMI if calculable.
func (m *CustomerMeasurement) BMI() *float64 {
	return m.measurement.BMI()
}
