package shared

import (
	"errors"
	"strings"
	"unicode"
)

// PersonName errors
var (
	ErrEmptyFirstName = errors.New("first name cannot be empty")
	ErrEmptyLastName  = errors.New("last name cannot be empty")
)

// PersonName represents a person's name.
type PersonName struct {
	firstName string
	lastName  string
}

// NewPersonName creates a new PersonName with validation.
func NewPersonName(firstName, lastName string) (PersonName, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if firstName == "" {
		return PersonName{}, ErrEmptyFirstName
	}
	if lastName == "" {
		return PersonName{}, ErrEmptyLastName
	}

	return PersonName{
		firstName: firstName,
		lastName:  lastName,
	}, nil
}

// NewPersonNameOptionalLast creates a name where last name is optional.
func NewPersonNameOptionalLast(firstName, lastName string) (PersonName, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if firstName == "" {
		return PersonName{}, ErrEmptyFirstName
	}

	return PersonName{
		firstName: firstName,
		lastName:  lastName,
	}, nil
}

// MustPersonName creates a PersonName, panicking on error.
func MustPersonName(firstName, lastName string) PersonName {
	n, err := NewPersonName(firstName, lastName)
	if err != nil {
		panic(err)
	}
	return n
}

// FirstName returns the first name.
func (n PersonName) FirstName() string {
	return n.firstName
}

// LastName returns the last name.
func (n PersonName) LastName() string {
	return n.lastName
}

// FullName returns the full name.
func (n PersonName) FullName() string {
	if n.lastName == "" {
		return n.firstName
	}
	return n.firstName + " " + n.lastName
}

// String returns the string representation.
func (n PersonName) String() string {
	return n.FullName()
}

// Initials returns the initials.
func (n PersonName) Initials() string {
	var initials strings.Builder

	if n.firstName != "" {
		for _, r := range n.firstName {
			if unicode.IsLetter(r) {
				initials.WriteRune(unicode.ToUpper(r))
				break
			}
		}
	}

	if n.lastName != "" {
		for _, r := range n.lastName {
			if unicode.IsLetter(r) {
				initials.WriteRune(unicode.ToUpper(r))
				break
			}
		}
	}

	return initials.String()
}

// IsEmpty returns true if the name is empty.
func (n PersonName) IsEmpty() bool {
	return n.firstName == "" && n.lastName == ""
}

// Equals compares two names (case-insensitive).
func (n PersonName) Equals(other PersonName) bool {
	return strings.EqualFold(n.firstName, other.firstName) &&
		strings.EqualFold(n.lastName, other.lastName)
}

// FormalName returns formal name format (Last, First).
func (n PersonName) FormalName() string {
	if n.lastName == "" {
		return n.firstName
	}
	return n.lastName + ", " + n.firstName
}
