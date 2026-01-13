package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"gorm.io/gorm"
)

// AddressRepository handles address data operations
type AddressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// ListByUserID retrieves all addresses for a user
func (r *AddressRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Address, error) {
	var addresses []domain.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// GetByID retrieves an address by ID with ownership check
func (r *AddressRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Create creates a new address
func (r *AddressRepository) Create(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// If this address is set as default, clear other defaults first
		if address.IsDefault {
			if err := tx.Model(&domain.Address{}).
				Where("user_id = ? AND is_default = ?", address.UserID, true).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(address).Error
	})
}

// Update updates an existing address
func (r *AddressRepository) Update(ctx context.Context, address *domain.Address) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// If this address is set as default, clear other defaults first
		if address.IsDefault {
			if err := tx.Model(&domain.Address{}).
				Where("user_id = ? AND id != ? AND is_default = ?", address.UserID, address.ID, true).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(address).Error
	})
}

// Delete deletes an address with ownership check
func (r *AddressRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Address{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// SetDefault sets an address as the default address
func (r *AddressRepository) SetDefault(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verify address exists and belongs to user
		var address domain.Address
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
			return err
		}

		// Clear all other defaults for this user
		if err := tx.Model(&domain.Address{}).
			Where("user_id = ? AND is_default = ?", userID, true).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Set this address as default
		return tx.Model(&domain.Address{}).
			Where("id = ?", id).
			Update("is_default", true).Error
	})
}
