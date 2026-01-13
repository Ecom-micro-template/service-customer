package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/models"
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

// AddWishlistItemInput contains all fields for adding a wishlist item
type AddWishlistItemInput struct {
	ProductID    uuid.UUID
	VariantID    *uuid.UUID
	VariantSKU   *string
	VariantName  *string
	PriceAtAdd   float64
	NotifyOnSale bool
	ProductName  *string
	ProductSlug  *string
	ProductImage *string
}

// Add adds a product to the wishlist (handles duplicates)
// Deprecated: Use AddWithVariant instead for variant support
func (r *WishlistRepository) Add(ctx context.Context, userID, productID uuid.UUID) error {
	return r.AddWithVariant(ctx, userID, AddWishlistItemInput{
		ProductID: productID,
	})
}

// AddWithVariant adds a product/variant to the wishlist with full details
func (r *WishlistRepository) AddWithVariant(ctx context.Context, userID uuid.UUID, input AddWishlistItemInput) error {
	// Build query to check for existing item
	query := r.db.WithContext(ctx).Model(&models.WishlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, input.ProductID)

	// Check variant-specific or product-level
	if input.VariantID != nil {
		query = query.Where("variant_id = ?", *input.VariantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}

	// If already exists, just return success
	if count > 0 {
		return nil
	}

	// Create new wishlist item
	item := &models.WishlistItem{
		UserID:       userID,
		ProductID:    input.ProductID,
		VariantID:    input.VariantID,
		VariantSKU:   input.VariantSKU,
		VariantName:  input.VariantName,
		PriceAtAdd:   input.PriceAtAdd,
		NotifyOnSale: input.NotifyOnSale,
		ProductName:  input.ProductName,
		ProductSlug:  input.ProductSlug,
		ProductImage: input.ProductImage,
	}
	return r.db.WithContext(ctx).Create(item).Error
}

// Remove removes a product from the wishlist (any variant)
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

// RemoveWithVariant removes a specific product/variant from the wishlist
func (r *WishlistRepository) RemoveWithVariant(ctx context.Context, userID, productID uuid.UUID, variantID *uuid.UUID) error {
	query := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	result := query.Delete(&models.WishlistItem{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// RemoveByID removes a wishlist item by its ID
func (r *WishlistRepository) RemoveByID(ctx context.Context, userID, itemID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", itemID, userID).
		Delete(&models.WishlistItem{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Exists checks if a product is in the user's wishlist (any variant)
func (r *WishlistRepository) Exists(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.WishlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}

// ExistsWithVariant checks if a specific product/variant is in the user's wishlist
func (r *WishlistRepository) ExistsWithVariant(ctx context.Context, userID, productID uuid.UUID, variantID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.WishlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", *variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// GetByProductID retrieves all wishlist items for a specific product (all variants)
func (r *WishlistRepository) GetByProductID(ctx context.Context, userID, productID uuid.UUID) ([]models.WishlistItem, error) {
	var items []models.WishlistItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// UpdateNotifyOnSale updates the price drop notification setting
func (r *WishlistRepository) UpdateNotifyOnSale(ctx context.Context, userID, itemID uuid.UUID, notify bool) error {
	result := r.db.WithContext(ctx).
		Model(&models.WishlistItem{}).
		Where("id = ? AND user_id = ?", itemID, userID).
		Update("notify_on_sale", notify)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetItemsForPriceDropAlert retrieves items where notify_on_sale is true
func (r *WishlistRepository) GetItemsForPriceDropAlert(ctx context.Context) ([]models.WishlistItem, error) {
	var items []models.WishlistItem
	err := r.db.WithContext(ctx).
		Where("notify_on_sale = ?", true).
		Find(&items).Error
	return items, err
}

// CountByUserID returns the count of wishlist items for a user
func (r *WishlistRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.WishlistItem{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
