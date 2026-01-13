package persistence

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

func setupWishlistTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.WishlistItem{})
	require.NoError(t, err)

	return db
}

func TestWishlistRepository_Add(t *testing.T) {
	db := setupWishlistTestDB(t)
	repo := NewWishlistRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	productID := uuid.New()

	err := repo.Add(ctx, userID, productID)
	assert.NoError(t, err)

	// Verify added
	exists, err := repo.Exists(ctx, userID, productID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestWishlistRepository_AddDuplicate(t *testing.T) {
	db := setupWishlistTestDB(t)
	repo := NewWishlistRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	productID := uuid.New()

	// Add first time
	err := repo.Add(ctx, userID, productID)
	assert.NoError(t, err)

	// Add again (should not error)
	err = repo.Add(ctx, userID, productID)
	assert.NoError(t, err)

	// Verify only one item exists
	items, err := repo.ListByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestWishlistRepository_ListByUserID(t *testing.T) {
	db := setupWishlistTestDB(t)
	repo := NewWishlistRepository(db)
	ctx := context.Background()

	userID := uuid.New()

	// Add multiple products
	productIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	for _, productID := range productIDs {
		err := repo.Add(ctx, userID, productID)
		require.NoError(t, err)
	}

	// List wishlist
	items, err := repo.ListByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, items, 3)
}

func TestWishlistRepository_Remove(t *testing.T) {
	db := setupWishlistTestDB(t)
	repo := NewWishlistRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	productID := uuid.New()

	// Add product
	err := repo.Add(ctx, userID, productID)
	require.NoError(t, err)

	// Remove product
	err = repo.Remove(ctx, userID, productID)
	assert.NoError(t, err)

	// Verify removed
	exists, err := repo.Exists(ctx, userID, productID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestWishlistRepository_Exists(t *testing.T) {
	db := setupWishlistTestDB(t)
	repo := NewWishlistRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	productID := uuid.New()

	// Should not exist initially
	exists, err := repo.Exists(ctx, userID, productID)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Add product
	err = repo.Add(ctx, userID, productID)
	require.NoError(t, err)

	// Should exist now
	exists, err = repo.Exists(ctx, userID, productID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestWishlistRepository_RemoveNonExistent(t *testing.T) {
	db := setupWishlistTestDB(t)
	repo := NewWishlistRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	productID := uuid.New()

	// Remove non-existent product
	err := repo.Remove(ctx, userID, productID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
