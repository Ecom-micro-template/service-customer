package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProfileRepository handles profile data operations
type ProfileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetByUserID retrieves a profile by user ID
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Upsert creates or updates a profile
func (r *ProfileRepository) Upsert(ctx context.Context, profile *domain.Profile) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"full_name", "email", "phone", "date_of_birth", "gender", "profile_picture", "updated_at"}),
	}).Create(profile).Error
}

// Create creates a new profile
func (r *ProfileRepository) Create(ctx context.Context, profile *domain.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// Update updates an existing profile
func (r *ProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}
