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

func setupAddressTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Address{})
	require.NoError(t, err)

	return db
}

func TestAddressRepository_Create(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	address := &models.Address{
		UserID:        userID,
		Label:         "Home",
		RecipientName: "John Doe",
		Phone:         "+1234567890",
		AddressLine1:  "123 Main St",
		City:          "New York",
		State:         "NY",
		Postcode:      "10001",
		Country:       "USA",
		IsDefault:     true,
	}

	err := repo.Create(ctx, address)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, address.ID)
}

func TestAddressRepository_ListByUserID(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)
	ctx := context.Background()

	userID := uuid.New()

	// Create multiple addresses
	addresses := []*models.Address{
		{
			UserID:        userID,
			Label:         "Home",
			RecipientName: "John Doe",
			Phone:         "+1111111111",
			AddressLine1:  "123 Home St",
			City:          "City1",
			State:         "ST",
			Postcode:      "10001",
			Country:       "USA",
			IsDefault:     true,
		},
		{
			UserID:        userID,
			Label:         "Office",
			RecipientName: "John Doe",
			Phone:         "+2222222222",
			AddressLine1:  "456 Work Ave",
			City:          "City2",
			State:         "ST",
			Postcode:      "10002",
			Country:       "USA",
			IsDefault:     false,
		},
	}

	for _, addr := range addresses {
		err := repo.Create(ctx, addr)
		require.NoError(t, err)
	}

	// List addresses
	list, err := repo.ListByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.True(t, list[0].IsDefault) // Default should be first
}

func TestAddressRepository_SetDefault(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)
	ctx := context.Background()

	userID := uuid.New()

	// Create first address as default
	addr1 := &models.Address{
		UserID:        userID,
		Label:         "Home",
		RecipientName: "John Doe",
		Phone:         "+1111111111",
		AddressLine1:  "123 Home St",
		City:          "City1",
		State:         "ST",
		Postcode:      "10001",
		Country:       "USA",
		IsDefault:     true,
	}
	err := repo.Create(ctx, addr1)
	require.NoError(t, err)

	// Create second address
	addr2 := &models.Address{
		UserID:        userID,
		Label:         "Office",
		RecipientName: "John Doe",
		Phone:         "+2222222222",
		AddressLine1:  "456 Work Ave",
		City:          "City2",
		State:         "ST",
		Postcode:      "10002",
		Country:       "USA",
		IsDefault:     false,
	}
	err = repo.Create(ctx, addr2)
	require.NoError(t, err)

	// Set second address as default
	err = repo.SetDefault(ctx, addr2.ID, userID)
	assert.NoError(t, err)

	// Verify only addr2 is default
	list, err := repo.ListByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	for _, addr := range list {
		if addr.ID == addr2.ID {
			assert.True(t, addr.IsDefault)
		} else {
			assert.False(t, addr.IsDefault)
		}
	}
}

func TestAddressRepository_Delete(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	address := &models.Address{
		UserID:        userID,
		Label:         "Home",
		RecipientName: "John Doe",
		Phone:         "+1234567890",
		AddressLine1:  "123 Main St",
		City:          "New York",
		State:         "NY",
		Postcode:      "10001",
		Country:       "USA",
	}

	err := repo.Create(ctx, address)
	require.NoError(t, err)

	// Delete address
	err = repo.Delete(ctx, address.ID, userID)
	assert.NoError(t, err)

	// Verify deleted
	_, err = repo.GetByID(ctx, address.ID, userID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestAddressRepository_Update(t *testing.T) {
	db := setupAddressTestDB(t)
	repo := NewAddressRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	address := &models.Address{
		UserID:        userID,
		Label:         "Home",
		RecipientName: "John Doe",
		Phone:         "+1234567890",
		AddressLine1:  "123 Main St",
		City:          "New York",
		State:         "NY",
		Postcode:      "10001",
		Country:       "USA",
	}

	err := repo.Create(ctx, address)
	require.NoError(t, err)

	// Update address
	address.Label = "Home Sweet Home"
	address.City = "Los Angeles"
	err = repo.Update(ctx, address)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.GetByID(ctx, address.ID, userID)
	assert.NoError(t, err)
	assert.Equal(t, "Home Sweet Home", retrieved.Label)
	assert.Equal(t, "Los Angeles", retrieved.City)
}
