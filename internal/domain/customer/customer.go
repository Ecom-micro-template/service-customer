package customer

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain/shared"
)

// Domain errors for Customer aggregate
var (
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCustomer    = errors.New("invalid customer data")
	ErrCannotModify       = errors.New("customer cannot be modified in current state")
)

// Customer is the aggregate root for customer domain.
type Customer struct {
	id          uuid.UUID
	email       shared.Email
	name        shared.PersonName
	phone       shared.Phone
	avatarURL   string
	status      shared.CustomerStatus
	totalOrders int
	totalSpent  float64
	createdAt   time.Time
	updatedAt   time.Time

	// Related entities
	notes      []CustomerNote
	activities []CustomerActivity

	// Domain events
	events []Event
}

// CustomerParams contains parameters for creating a Customer.
type CustomerParams struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	Phone     string
}

// NewCustomer creates a new Customer aggregate.
func NewCustomer(params CustomerParams) (*Customer, error) {
	email, err := shared.NewEmail(params.Email)
	if err != nil {
		return nil, err
	}

	name, err := shared.NewPersonName(params.FirstName, params.LastName)
	if err != nil {
		return nil, err
	}

	phone := shared.EmptyPhone()
	if params.Phone != "" {
		phone, _ = shared.NewPhone(params.Phone)
	}

	id := params.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	now := time.Now()
	customer := &Customer{
		id:          id,
		email:       email,
		name:        name,
		phone:       phone,
		status:      shared.StatusActive,
		totalOrders: 0,
		totalSpent:  0,
		createdAt:   now,
		updatedAt:   now,
		notes:       make([]CustomerNote, 0),
		activities:  make([]CustomerActivity, 0),
		events:      make([]Event, 0),
	}

	customer.addEvent(NewCustomerCreatedEvent(id, email.Value(), name.FullName()))

	return customer, nil
}

// Getters
func (c *Customer) ID() uuid.UUID                  { return c.id }
func (c *Customer) Email() shared.Email            { return c.email }
func (c *Customer) Name() shared.PersonName        { return c.name }
func (c *Customer) Phone() shared.Phone            { return c.phone }
func (c *Customer) AvatarURL() string              { return c.avatarURL }
func (c *Customer) Status() shared.CustomerStatus  { return c.status }
func (c *Customer) TotalOrders() int               { return c.totalOrders }
func (c *Customer) TotalSpent() float64            { return c.totalSpent }
func (c *Customer) CreatedAt() time.Time           { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time           { return c.updatedAt }
func (c *Customer) Notes() []CustomerNote          { return c.notes }
func (c *Customer) Activities() []CustomerActivity { return c.activities }

// --- Behavior Methods ---

// UpdateProfile updates the customer's profile.
func (c *Customer) UpdateProfile(firstName, lastName string, phone string) error {
	if !c.status.CanLogin() && c.status != shared.StatusInactive {
		return ErrCannotModify
	}

	if firstName != "" && lastName != "" {
		name, err := shared.NewPersonName(firstName, lastName)
		if err != nil {
			return err
		}
		c.name = name
	}

	if phone != "" {
		p, err := shared.NewPhone(phone)
		if err != nil {
			return err
		}
		c.phone = p
	}

	c.updatedAt = time.Now()
	c.addEvent(NewCustomerUpdatedEvent(c.id))
	return nil
}

// SetAvatarURL sets the avatar URL.
func (c *Customer) SetAvatarURL(url string) {
	c.avatarURL = url
	c.updatedAt = time.Now()
}

// Activate activates the customer.
func (c *Customer) Activate() error {
	if !c.status.CanBeActivated() {
		return ErrCannotModify
	}
	c.status = shared.StatusActive
	c.updatedAt = time.Now()
	c.addEvent(NewCustomerStatusChangedEvent(c.id, string(c.status)))
	return nil
}

// Suspend suspends the customer.
func (c *Customer) Suspend(reason string) error {
	if !c.status.CanBeSuspended() {
		return ErrCannotModify
	}
	c.status = shared.StatusSuspended
	c.updatedAt = time.Now()
	c.addEvent(NewCustomerStatusChangedEvent(c.id, string(c.status)))
	c.AddNote("Suspended: "+reason, true, nil)
	return nil
}

// Block blocks the customer.
func (c *Customer) Block(reason string) error {
	if !c.status.CanBeBlocked() {
		return ErrCannotModify
	}
	c.status = shared.StatusBlocked
	c.updatedAt = time.Now()
	c.addEvent(NewCustomerStatusChangedEvent(c.id, string(c.status)))
	c.AddNote("Blocked: "+reason, true, nil)
	return nil
}

// RecordOrder records an order for the customer.
func (c *Customer) RecordOrder(orderTotal float64) {
	c.totalOrders++
	c.totalSpent += orderTotal
	c.updatedAt = time.Now()
	c.RecordActivity("order", "Order Placed", "")
}

// AddNote adds a note to the customer.
func (c *Customer) AddNote(note string, isPrivate bool, createdBy *uuid.UUID) {
	customerNote := NewCustomerNote(c.id, note, isPrivate, createdBy)
	c.notes = append(c.notes, customerNote)
}

// RecordActivity records an activity.
func (c *Customer) RecordActivity(activityType, title, details string) {
	activity := NewCustomerActivity(c.id, activityType, title, details)
	c.activities = append(c.activities, activity)
}

// CanLogin returns true if the customer can log in.
func (c *Customer) CanLogin() bool {
	return c.status.CanLogin()
}

// CanPurchase returns true if the customer can make purchases.
func (c *Customer) CanPurchase() bool {
	return c.status.CanPurchase()
}

// IsActive returns true if customer is active.
func (c *Customer) IsActive() bool {
	return c.status == shared.StatusActive
}

// Events returns and clears the collected domain events.
func (c *Customer) Events() []Event {
	events := c.events
	c.events = make([]Event, 0)
	return events
}

func (c *Customer) addEvent(event Event) {
	c.events = append(c.events, event)
}
