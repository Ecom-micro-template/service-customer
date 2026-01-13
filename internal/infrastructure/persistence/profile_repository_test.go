package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate models
	err = db.AutoMigrate(&models.Profile{})
	require.NoError(t, err)

	return db
}

func TestProfileRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)
	ctx := context.Background()

	profile := &models.Profile{
		ID:       uuid.New(),
		FullName: "John Doe",
		Email:    "john@example.com",
		Phone:    "+1234567890",
		Gender:   "male",
	}

	err := repo.Create(ctx, profile)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, profile.ID)
}

func TestProfileRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	profile := &models.Profile{
		ID:       userID,
		FullName: "Jane Doe",
		Email:    "jane@example.com",
		Phone:    "+0987654321",
	}

	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	retrieved, err := repo.GetByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "Jane Doe", retrieved.FullName)
	assert.Equal(t, "jane@example.com", retrieved.Email)
}

func TestProfileRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	profile := &models.Profile{
		ID:       userID,
		FullName: "Original Name",
		Email:    "original@example.com",
	}

	err := repo.Create(ctx, profile)
	require.NoError(t, err)

	// Update profile
	profile.FullName = "Updated Name"
	profile.Email = "updated@example.com"
	err = repo.Update(ctx, profile)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.FullName)
	assert.Equal(t, "updated@example.com", retrieved.Email)
}

func TestProfileRepository_Upsert(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProfileRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	profile := &models.Profile{
		ID:       userID,
		FullName: "First Insert",
		Email:    "first@example.com",
	}

	// First upsert (insert)
	err := repo.Upsert(ctx, profile)
	assert.NoError(t, err)

	// Second upsert (update)
	profile.FullName = "Second Insert"
	profile.Email = "second@example.com"
	err = repo.Upsert(ctx, profile)
	assert.NoError(t, err)

	// Verify only one record exists with updated data
	retrieved, err := repo.GetByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, "Second Insert", retrieved.FullName)
	assert.Equal(t, "second@example.com", retrieved.Email)
}
