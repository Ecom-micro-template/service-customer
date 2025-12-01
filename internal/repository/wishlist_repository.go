package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/niaga-platform/service-customer/internal/models"
	"gorm.io/gorm"
)

// WishlistRepository handles wishlist data operations
type WishlistRepository struct {
	db *gorm.DB
}

// NewWishlistRepository creates a new wishlist repository
func NewWishlistRepository(db *gorm.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

// ListByUserID retrieves all wishlist items for a user
func (r *WishlistRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.WishlistItem, error) {
	var items []models.WishlistItem
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// Add adds a product to the wishlist (handles duplicates)
func (r *WishlistRepository) Add(ctx context.Context, userID, productID uuid.UUID) error {
	// Check if already exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.WishlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error; err != nil {
		return err
	}

	// If already exists, just return success
	if count > 0 {
		return nil
	}

	// Create new wishlist item
	item := &models.WishlistItem{
		UserID:    userID,
		ProductID: productID,
	}
	return r.db.WithContext(ctx).Create(item).Error
}

// Remove removes a product from the wishlist
func (r *WishlistRepository) Remove(ctx context.Context, userID, productID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.WishlistItem{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Exists checks if a product is in the user's wishlist
func (r *WishlistRepository) Exists(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.WishlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}
