package address

import (
	"errors"
	"fmt"
)

// AddressType represents the type of address.
type AddressType string

const (
	TypeHome   AddressType = "home"
	TypeOffice AddressType = "office"
	TypeOther  AddressType = "other"
)

// ErrInvalidAddressType is returned for invalid address types.
var ErrInvalidAddressType = errors.New("invalid address type")

// IsValid returns true if the type is valid.
func (t AddressType) IsValid() bool {
	switch t {
	case TypeHome, TypeOffice, TypeOther:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (t AddressType) String() string {
	return string(t)
}

// Label returns a human-readable label.
func (t AddressType) Label() string {
	switch t {
	case TypeHome:
		return "Home"
	case TypeOffice:
		return "Office"
	case TypeOther:
		return "Other"
	default:
		return "Unknown"
	}
}

// ParseAddressType parses a string into an AddressType.
func ParseAddressType(s string) (AddressType, error) {
	t := AddressType(s)
	if !t.IsValid() {
		return "", fmt.Errorf("%w: %s", ErrInvalidAddressType, s)
	}
	return t, nil
}
