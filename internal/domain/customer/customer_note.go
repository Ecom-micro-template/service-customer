package customer

import (
	"time"

	"github.com/google/uuid"
)

// CustomerNote represents a note on a customer.
type CustomerNote struct {
	id         uuid.UUID
	customerID uuid.UUID
	note       string
	isPrivate  bool
	createdBy  *uuid.UUID
	createdAt  time.Time
}

// NewCustomerNote creates a new CustomerNote.
func NewCustomerNote(customerID uuid.UUID, note string, isPrivate bool, createdBy *uuid.UUID) CustomerNote {
	return CustomerNote{
		id:         uuid.New(),
		customerID: customerID,
		note:       note,
		isPrivate:  isPrivate,
		createdBy:  createdBy,
		createdAt:  time.Now(),
	}
}

// Getters
func (n CustomerNote) ID() uuid.UUID         { return n.id }
func (n CustomerNote) CustomerID() uuid.UUID { return n.customerID }
func (n CustomerNote) Note() string          { return n.note }
func (n CustomerNote) IsPrivate() bool       { return n.isPrivate }
func (n CustomerNote) CreatedBy() *uuid.UUID { return n.createdBy }
func (n CustomerNote) CreatedAt() time.Time  { return n.createdAt }
