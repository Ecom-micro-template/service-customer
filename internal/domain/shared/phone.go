package shared

import (
	"errors"
	"regexp"
	"strings"
)

// Phone errors
var (
	ErrInvalidPhone = errors.New("invalid phone number format")
	ErrEmptyPhone   = errors.New("phone number cannot be empty")
)

// phoneRegex validates phone numbers (allows +, digits, spaces, dashes, parentheses)
var phoneRegex = regexp.MustCompile(`^[\+]?[(]?[0-9]{1,4}[)]?[-\s\./0-9]*$`)

// Phone represents a validated phone number.
type Phone struct {
	value   string
	country string
}

// NewPhone creates a new Phone with validation.
func NewPhone(phone string) (Phone, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return Phone{}, ErrEmptyPhone
	}
	if !phoneRegex.MatchString(phone) {
		return Phone{}, ErrInvalidPhone
	}
	return Phone{value: phone}, nil
}

// NewPhoneWithCountry creates a Phone with country code.
func NewPhoneWithCountry(phone, country string) (Phone, error) {
	p, err := NewPhone(phone)
	if err != nil {
		return Phone{}, err
	}
	p.country = country
	return p, nil
}

// MustPhone creates a Phone, panicking on error.
func MustPhone(phone string) Phone {
	p, err := NewPhone(phone)
	if err != nil {
		panic(err)
	}
	return p
}

// EmptyPhone returns an empty phone (for optional fields).
func EmptyPhone() Phone {
	return Phone{}
}

// Value returns the phone string.
func (p Phone) Value() string {
	return p.value
}

// String returns the string representation.
func (p Phone) String() string {
	return p.value
}

// Country returns the country code.
func (p Phone) Country() string {
	return p.country
}

// IsEmpty returns true if the phone is empty.
func (p Phone) IsEmpty() bool {
	return p.value == ""
}

// Equals compares two phones.
func (p Phone) Equals(other Phone) bool {
	return p.value == other.value
}

// Normalized returns phone with only digits (and optional + prefix).
func (p Phone) Normalized() string {
	var result strings.Builder
	for i, ch := range p.value {
		if ch >= '0' && ch <= '9' {
			result.WriteRune(ch)
		} else if ch == '+' && i == 0 {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// MaskedPhone returns a masked version for display.
// e.g., "+60123456789" -> "+60****6789"
func (p Phone) MaskedPhone() string {
	normalized := p.Normalized()
	if len(normalized) <= 6 {
		return normalized
	}
	visible := 4
	return normalized[:len(normalized)-visible-4] + "****" + normalized[len(normalized)-visible:]
}
