package customer

import (
	"time"

	"github.com/google/uuid"
)

// Event is the base interface for all customer domain events.
type Event interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() uuid.UUID
}

// baseEvent contains common event fields.
type baseEvent struct {
	occurredAt  time.Time
	aggregateID uuid.UUID
}

func (e baseEvent) OccurredAt() time.Time  { return e.occurredAt }
func (e baseEvent) AggregateID() uuid.UUID { return e.aggregateID }

// CustomerCreatedEvent is raised when a new customer is created.
type CustomerCreatedEvent struct {
	baseEvent
	Email string
	Name  string
}

func (e CustomerCreatedEvent) EventType() string { return "customer.created" }

// NewCustomerCreatedEvent creates a new CustomerCreatedEvent.
func NewCustomerCreatedEvent(customerID uuid.UUID, email, name string) CustomerCreatedEvent {
	return CustomerCreatedEvent{
		baseEvent: baseEvent{occurredAt: time.Now(), aggregateID: customerID},
		Email:     email,
		Name:      name,
	}
}

// CustomerUpdatedEvent is raised when a customer is updated.
type CustomerUpdatedEvent struct {
	baseEvent
}

func (e CustomerUpdatedEvent) EventType() string { return "customer.updated" }

// NewCustomerUpdatedEvent creates a new CustomerUpdatedEvent.
func NewCustomerUpdatedEvent(customerID uuid.UUID) CustomerUpdatedEvent {
	return CustomerUpdatedEvent{
		baseEvent: baseEvent{occurredAt: time.Now(), aggregateID: customerID},
	}
}

// CustomerStatusChangedEvent is raised when customer status changes.
type CustomerStatusChangedEvent struct {
	baseEvent
	NewStatus string
}

func (e CustomerStatusChangedEvent) EventType() string { return "customer.status_changed" }

// NewCustomerStatusChangedEvent creates a new CustomerStatusChangedEvent.
func NewCustomerStatusChangedEvent(customerID uuid.UUID, newStatus string) CustomerStatusChangedEvent {
	return CustomerStatusChangedEvent{
		baseEvent: baseEvent{occurredAt: time.Now(), aggregateID: customerID},
		NewStatus: newStatus,
	}
}

// CustomerDeletedEvent is raised when a customer is deleted.
type CustomerDeletedEvent struct {
	baseEvent
}

func (e CustomerDeletedEvent) EventType() string { return "customer.deleted" }

// NewCustomerDeletedEvent creates a new CustomerDeletedEvent.
func NewCustomerDeletedEvent(customerID uuid.UUID) CustomerDeletedEvent {
	return CustomerDeletedEvent{
		baseEvent: baseEvent{occurredAt: time.Now(), aggregateID: customerID},
	}
}
