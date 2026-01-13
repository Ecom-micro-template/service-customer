package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"gorm.io/gorm"
)

// MeasurementRepository handles database operations for customer measurements
type MeasurementRepository struct {
	db *gorm.DB
}

// NewMeasurementRepository creates a new measurement repository
func NewMeasurementRepository(db *gorm.DB) *MeasurementRepository {
	return &MeasurementRepository{db: db}
}

// Create creates a new customer measurement
func (r *MeasurementRepository) Create(ctx context.Context, measurement *models.CustomerMeasurement) error {
	return r.db.WithContext(ctx).Create(measurement).Error
}

// GetByID retrieves a measurement by ID with user ownership check (IDOR protection)
func (r *MeasurementRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.CustomerMeasurement, error) {
	var measurement models.CustomerMeasurement
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&measurement).Error
	if err != nil {
		return nil, err
	}
	return &measurement, nil
}

// GetByUserID retrieves all measurements for a user
func (r *MeasurementRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.CustomerMeasurement, error) {
	var measurements []models.CustomerMeasurement
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&measurements).Error
	return measurements, err
}

// GetDefaultByUserID retrieves the default measurement for a user
func (r *MeasurementRepository) GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*models.CustomerMeasurement, error) {
	var measurement models.CustomerMeasurement
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&measurement).Error
	if err != nil {
		return nil, err
	}
	return &measurement, nil
}

// Update updates a measurement
func (r *MeasurementRepository) Update(ctx context.Context, measurement *models.CustomerMeasurement) error {
	return r.db.WithContext(ctx).Save(measurement).Error
}

// Delete deletes a measurement with user ownership check (IDOR protection)
func (r *MeasurementRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.CustomerMeasurement{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// SetDefault sets a measurement as default and unsets others for the user
func (r *MeasurementRepository) SetDefault(ctx context.Context, userID, measurementID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all default measurements for this user
		if err := tx.Model(&models.CustomerMeasurement{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Set the new default
		return tx.Model(&models.CustomerMeasurement{}).
			Where("id = ? AND user_id = ?", measurementID, userID).
			Update("is_default", true).Error
	})
}
