package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/niaga-platform/service-customer/internal/domain/customer"
	"github.com/niaga-platform/service-customer/internal/domain/shared"
	"gorm.io/gorm"
)

// CustomerRepository interface for customer domain.
type CustomerRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error)
	GetByEmail(ctx context.Context, email string) (*customer.Customer, error)
	List(ctx context.Context, page, pageSize int, status *string) ([]*customer.Customer, int64, error)
	Save(ctx context.Context, c *customer.Customer) error
	Update(ctx context.Context, c *customer.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// GormCustomerRepository is the GORM implementation of CustomerRepository.
type GormCustomerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new GormCustomerRepository.
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &GormCustomerRepository{db: db}
}

// GetByID retrieves a customer by ID.
func (r *GormCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error) {
	var model CustomerModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, err
	}
	return r.toDomain(&model)
}

// GetByEmail retrieves a customer by email.
func (r *GormCustomerRepository) GetByEmail(ctx context.Context, email string) (*customer.Customer, error) {
	var model CustomerModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, err
	}
	return r.toDomain(&model)
}

// List retrieves customers with pagination.
func (r *GormCustomerRepository) List(ctx context.Context, page, pageSize int, status *string) ([]*customer.Customer, int64, error) {
	var models []CustomerModel
	var total int64

	query := r.db.WithContext(ctx).Model(&CustomerModel{})
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	customers := make([]*customer.Customer, len(models))
	for i, m := range models {
		c, err := r.toDomain(&m)
		if err != nil {
			return nil, 0, err
		}
		customers[i] = c
	}
	return customers, total, nil
}

// Save creates a new customer.
func (r *GormCustomerRepository) Save(ctx context.Context, c *customer.Customer) error {
	model := r.toModel(c)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update updates an existing customer.
func (r *GormCustomerRepository) Update(ctx context.Context, c *customer.Customer) error {
	model := r.toModel(c)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete soft-deletes a customer.
func (r *GormCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&CustomerModel{}, "id = ?", id).Error
}

// toDomain converts model to domain entity.
func (r *GormCustomerRepository) toDomain(m *CustomerModel) (*customer.Customer, error) {
	return customer.NewCustomer(customer.CustomerParams{
		ID:        m.ID,
		Email:     m.Email,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Phone:     m.Phone,
	})
}

// toModel converts domain entity to model.
func (r *GormCustomerRepository) toModel(c *customer.Customer) *CustomerModel {
	phone := ""
	if !c.Phone().IsEmpty() {
		phone = c.Phone().Value()
	}

	status, _ := shared.ParseCustomerStatus(string(c.Status()))

	return &CustomerModel{
		ID:          c.ID(),
		Email:       c.Email().Value(),
		FirstName:   c.Name().FirstName(),
		LastName:    c.Name().LastName(),
		Phone:       phone,
		AvatarURL:   c.AvatarURL(),
		Status:      status.String(),
		TotalOrders: c.TotalOrders(),
		TotalSpent:  c.TotalSpent(),
		CreatedAt:   c.CreatedAt(),
		UpdatedAt:   c.UpdatedAt(),
	}
}
