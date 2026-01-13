package address

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain/shared"
)

// Domain errors for Address aggregate
var (
	ErrAddressNotFound = errors.New("address not found")
	ErrInvalidAddress  = errors.New("invalid address data")
	ErrMaxAddresses    = errors.New("maximum number of addresses reached")
)

// Address is the aggregate root for customer addresses.
type Address struct {
	id            uuid.UUID
	userID        uuid.UUID
	label         AddressType
	recipientName string
	phone         shared.Phone
	addressLine1  string
	addressLine2  string
	city          string
	state         string
	postcode      string
	country       string
	isDefault     bool
	createdAt     time.Time
	updatedAt     time.Time
}

// AddressParams contains parameters for creating an Address.
type AddressParams struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Label         string
	RecipientName string
	Phone         string
	AddressLine1  string
	AddressLine2  string
	City          string
	State         string
	Postcode      string
	Country       string
	IsDefault     bool
}

// NewAddress creates a new Address aggregate.
func NewAddress(params AddressParams) (*Address, error) {
	if params.UserID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if strings.TrimSpace(params.RecipientName) == "" {
		return nil, errors.New("recipient name is required")
	}
	if strings.TrimSpace(params.AddressLine1) == "" {
		return nil, errors.New("address line 1 is required")
	}
	if strings.TrimSpace(params.City) == "" {
		return nil, errors.New("city is required")
	}
	if strings.TrimSpace(params.State) == "" {
		return nil, errors.New("state is required")
	}
	if strings.TrimSpace(params.Postcode) == "" {
		return nil, errors.New("postcode is required")
	}

	phone, _ := shared.NewPhone(params.Phone)

	label, err := ParseAddressType(params.Label)
	if err != nil {
		label = TypeHome // Default
	}

	id := params.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	country := params.Country
	if country == "" {
		country = "Malaysia"
	}

	now := time.Now()
	return &Address{
		id:            id,
		userID:        params.UserID,
		label:         label,
		recipientName: strings.TrimSpace(params.RecipientName),
		phone:         phone,
		addressLine1:  strings.TrimSpace(params.AddressLine1),
		addressLine2:  strings.TrimSpace(params.AddressLine2),
		city:          strings.TrimSpace(params.City),
		state:         strings.TrimSpace(params.State),
		postcode:      strings.TrimSpace(params.Postcode),
		country:       country,
		isDefault:     params.IsDefault,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// Getters
func (a *Address) ID() uuid.UUID         { return a.id }
func (a *Address) UserID() uuid.UUID     { return a.userID }
func (a *Address) Label() AddressType    { return a.label }
func (a *Address) RecipientName() string { return a.recipientName }
func (a *Address) Phone() shared.Phone   { return a.phone }
func (a *Address) AddressLine1() string  { return a.addressLine1 }
func (a *Address) AddressLine2() string  { return a.addressLine2 }
func (a *Address) City() string          { return a.city }
func (a *Address) State() string         { return a.state }
func (a *Address) Postcode() string      { return a.postcode }
func (a *Address) Country() string       { return a.country }
func (a *Address) IsDefault() bool       { return a.isDefault }
func (a *Address) CreatedAt() time.Time  { return a.createdAt }
func (a *Address) UpdatedAt() time.Time  { return a.updatedAt }

// FullAddress returns the formatted full address.
func (a *Address) FullAddress() string {
	parts := []string{a.addressLine1}
	if a.addressLine2 != "" {
		parts = append(parts, a.addressLine2)
	}
	parts = append(parts, a.city, a.state, a.postcode, a.country)
	return strings.Join(parts, ", ")
}

// --- Behavior Methods ---

// Update updates the address details.
func (a *Address) Update(params AddressParams) error {
	if params.RecipientName != "" {
		a.recipientName = strings.TrimSpace(params.RecipientName)
	}
	if params.Phone != "" {
		phone, err := shared.NewPhone(params.Phone)
		if err == nil {
			a.phone = phone
		}
	}
	if params.AddressLine1 != "" {
		a.addressLine1 = strings.TrimSpace(params.AddressLine1)
	}
	a.addressLine2 = strings.TrimSpace(params.AddressLine2)
	if params.City != "" {
		a.city = strings.TrimSpace(params.City)
	}
	if params.State != "" {
		a.state = strings.TrimSpace(params.State)
	}
	if params.Postcode != "" {
		a.postcode = strings.TrimSpace(params.Postcode)
	}
	if params.Country != "" {
		a.country = params.Country
	}
	if params.Label != "" {
		label, err := ParseAddressType(params.Label)
		if err == nil {
			a.label = label
		}
	}

	a.updatedAt = time.Now()
	return nil
}

// SetDefault sets this address as the default.
func (a *Address) SetDefault() {
	a.isDefault = true
	a.updatedAt = time.Now()
}

// ClearDefault clears the default flag.
func (a *Address) ClearDefault() {
	a.isDefault = false
	a.updatedAt = time.Now()
}

// SetLabel sets the address label.
func (a *Address) SetLabel(label AddressType) {
	if label.IsValid() {
		a.label = label
		a.updatedAt = time.Now()
	}
}
