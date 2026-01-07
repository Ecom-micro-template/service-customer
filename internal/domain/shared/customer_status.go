package shared

import (
	"errors"
	"fmt"
)

// CustomerStatus represents the status of a customer account.
type CustomerStatus string

// Customer status constants
const (
	StatusActive    CustomerStatus = "active"
	StatusInactive  CustomerStatus = "inactive"
	StatusSuspended CustomerStatus = "suspended"
	StatusBlocked   CustomerStatus = "blocked"
)

// ErrInvalidCustomerStatus is returned for invalid status values.
var ErrInvalidCustomerStatus = errors.New("invalid customer status")

// AllCustomerStatuses returns all valid statuses.
func AllCustomerStatuses() []CustomerStatus {
	return []CustomerStatus{StatusActive, StatusInactive, StatusSuspended, StatusBlocked}
}

// IsValid returns true if the status is valid.
func (s CustomerStatus) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusSuspended, StatusBlocked:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (s CustomerStatus) String() string {
	return string(s)
}

// Label returns a human-readable label.
func (s CustomerStatus) Label() string {
	switch s {
	case StatusActive:
		return "Active"
	case StatusInactive:
		return "Inactive"
	case StatusSuspended:
		return "Suspended"
	case StatusBlocked:
		return "Blocked"
	default:
		return "Unknown"
	}
}

// CanLogin returns true if the customer can log in with this status.
func (s CustomerStatus) CanLogin() bool {
	return s == StatusActive
}

// CanPurchase returns true if the customer can make purchases.
func (s CustomerStatus) CanPurchase() bool {
	return s == StatusActive
}

// CanBeActivated returns true if this status can transition to active.
func (s CustomerStatus) CanBeActivated() bool {
	return s == StatusInactive || s == StatusSuspended
}

// CanBeSuspended returns true if this status can transition to suspended.
func (s CustomerStatus) CanBeSuspended() bool {
	return s == StatusActive
}

// CanBeBlocked returns true if this status can transition to blocked.
func (s CustomerStatus) CanBeBlocked() bool {
	return s != StatusBlocked
}

// ParseCustomerStatus parses a string into a CustomerStatus.
func ParseCustomerStatus(s string) (CustomerStatus, error) {
	status := CustomerStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidCustomerStatus, s)
	}
	return status, nil
}
