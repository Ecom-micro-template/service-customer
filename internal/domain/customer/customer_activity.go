package customer

import (
	"time"

	"github.com/google/uuid"
)

// CustomerActivity represents a customer activity log entry.
type CustomerActivity struct {
	id           uuid.UUID
	customerID   uuid.UUID
	activityType string
	title        string
	details      string
	createdAt    time.Time
}

// Activity type constants
const (
	ActivityTypeLogin       = "login"
	ActivityTypeOrder       = "order"
	ActivityTypeWishlist    = "wishlist"
	ActivityTypeAddress     = "address"
	ActivityTypeProfile     = "profile"
	ActivityTypeMeasurement = "measurement"
)

// NewCustomerActivity creates a new CustomerActivity.
func NewCustomerActivity(customerID uuid.UUID, activityType, title, details string) CustomerActivity {
	return CustomerActivity{
		id:           uuid.New(),
		customerID:   customerID,
		activityType: activityType,
		title:        title,
		details:      details,
		createdAt:    time.Now(),
	}
}

// Getters
func (a CustomerActivity) ID() uuid.UUID         { return a.id }
func (a CustomerActivity) CustomerID() uuid.UUID { return a.customerID }
func (a CustomerActivity) Type() string          { return a.activityType }
func (a CustomerActivity) Title() string         { return a.title }
func (a CustomerActivity) Details() string       { return a.details }
func (a CustomerActivity) CreatedAt() time.Time  { return a.createdAt }
