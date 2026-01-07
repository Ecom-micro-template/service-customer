// Package shared provides shared value objects for the customer domain.
package shared

import (
	"errors"
	"regexp"
	"strings"
)

// Email errors
var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrEmptyEmail   = errors.New("email cannot be empty")
)

// emailRegex is a simple email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email represents a validated email address.
type Email struct {
	value string
}

// NewEmail creates a new Email with validation.
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return Email{}, ErrEmptyEmail
	}
	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: email}, nil
}

// MustEmail creates an Email, panicking on error.
func MustEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return e
}

// Value returns the email string.
func (e Email) Value() string {
	return e.value
}

// String returns the string representation.
func (e Email) String() string {
	return e.value
}

// Domain returns the domain part of the email.
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// LocalPart returns the local part of the email (before @).
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

// IsEmpty returns true if the email is empty.
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Equals compares two emails.
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// MaskedEmail returns a masked version for display.
// e.g., "john.doe@example.com" -> "j***e@example.com"
func (e Email) MaskedEmail() string {
	local := e.LocalPart()
	domain := e.Domain()
	if len(local) <= 2 {
		return local + "@" + domain
	}
	return string(local[0]) + "***" + string(local[len(local)-1]) + "@" + domain
}
