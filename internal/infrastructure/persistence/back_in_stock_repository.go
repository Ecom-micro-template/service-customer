package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"gorm.io/gorm"
)

// HI-001: Back-in-Stock Repository

// BackInStockRepository handles back-in-stock subscription database operations
type BackInStockRepository struct {
	db *gorm.DB
}

// NewBackInStockRepository creates a new repository
func NewBackInStockRepository(db *gorm.DB) *BackInStockRepository {
	return &BackInStockRepository{db: db}
}

// Subscribe creates a new subscription or returns existing one
func (r *BackInStockRepository) Subscribe(ctx context.Context, customerID uuid.UUID, input domain.BackInStockSubscribeInput) (*domain.BackInStockSubscription, error) {
	productID, err := uuid.Parse(input.ProductID)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	var variantID *uuid.UUID
	if input.VariantID != "" {
		vid, err := uuid.Parse(input.VariantID)
		if err != nil {
			return nil, errors.New("invalid variant ID")
		}
		variantID = &vid
	}

	// Check if subscription already exists
	var existing domain.BackInStockSubscription
	query := r.db.WithContext(ctx).Where("customer_id = ? AND product_id = ?", customerID, productID)
	if variantID != nil {
		query = query.Where("variant_id = ?", variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	if err := query.First(&existing).Error; err == nil {
		// Already subscribed, just return it
		return &existing, nil
	}

	// Create new subscription
	subscription := domain.BackInStockSubscription{
		CustomerID:   customerID,
		ProductID:    productID,
		VariantID:    variantID,
		ProductName:  input.ProductName,
		ProductSlug:  input.ProductSlug,
		ProductImage: input.ProductImage,
		VariantSKU:   input.VariantSKU,
		VariantName:  input.VariantName,
		IsNotified:   false,
	}

	if err := r.db.WithContext(ctx).Create(&subscription).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

// Unsubscribe removes a subscription
func (r *BackInStockRepository) Unsubscribe(ctx context.Context, customerID, productID uuid.UUID, variantID *uuid.UUID) error {
	query := r.db.WithContext(ctx).
		Where("customer_id = ? AND product_id = ?", customerID, productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	return query.Delete(&domain.BackInStockSubscription{}).Error
}

// UnsubscribeByID removes a subscription by ID
func (r *BackInStockRepository) UnsubscribeByID(ctx context.Context, customerID, subscriptionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND customer_id = ?", subscriptionID, customerID).
		Delete(&domain.BackInStockSubscription{}).Error
}

// GetByCustomer returns all subscriptions for a customer
func (r *BackInStockRepository) GetByCustomer(ctx context.Context, customerID uuid.UUID) ([]domain.BackInStockSubscription, error) {
	var subscriptions []domain.BackInStockSubscription
	err := r.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&subscriptions).Error
	return subscriptions, err
}

// GetByProduct returns all pending subscriptions for a product
func (r *BackInStockRepository) GetByProduct(ctx context.Context, productID uuid.UUID, variantID *uuid.UUID) ([]domain.BackInStockSubscription, error) {
	var subscriptions []domain.BackInStockSubscription
	query := r.db.WithContext(ctx).
		Preload("Customer").
		Where("product_id = ? AND is_notified = false", productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", variantID)
	}

	err := query.Find(&subscriptions).Error
	return subscriptions, err
}

// GetPendingNotifications returns all subscriptions that haven't been notified
func (r *BackInStockRepository) GetPendingNotifications(ctx context.Context, limit int) ([]domain.BackInStockSubscription, error) {
	var subscriptions []domain.BackInStockSubscription
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Where("is_notified = false").
		Limit(limit).
		Find(&subscriptions).Error
	return subscriptions, err
}

// MarkAsNotified marks a subscription as notified
func (r *BackInStockRepository) MarkAsNotified(ctx context.Context, subscriptionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.BackInStockSubscription{}).
		Where("id = ?", subscriptionID).
		Updates(map[string]interface{}{
			"is_notified":          true,
			"notification_sent_at": gorm.Expr("NOW()"),
		}).Error
}

// MarkMultipleAsNotified marks multiple subscriptions as notified
func (r *BackInStockRepository) MarkMultipleAsNotified(ctx context.Context, subscriptionIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.BackInStockSubscription{}).
		Where("id IN ?", subscriptionIDs).
		Updates(map[string]interface{}{
			"is_notified":          true,
			"notification_sent_at": gorm.Expr("NOW()"),
		}).Error
}

// IsSubscribed checks if a customer is subscribed to a product
func (r *BackInStockRepository) IsSubscribed(ctx context.Context, customerID, productID uuid.UUID, variantID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&domain.BackInStockSubscription{}).
		Where("customer_id = ? AND product_id = ?", customerID, productID)

	if variantID != nil {
		query = query.Where("variant_id = ?", variantID)
	} else {
		query = query.Where("variant_id IS NULL")
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetStats returns statistics about subscriptions
func (r *BackInStockRepository) GetStats(ctx context.Context) (*domain.BackInStockStats, error) {
	var stats domain.BackInStockStats

	// Total subscriptions
	r.db.WithContext(ctx).Model(&domain.BackInStockSubscription{}).Count(&stats.TotalSubscriptions)

	// Pending notifications
	r.db.WithContext(ctx).Model(&domain.BackInStockSubscription{}).
		Where("is_notified = false").Count(&stats.PendingNotifications)

	// Sent notifications
	r.db.WithContext(ctx).Model(&domain.BackInStockSubscription{}).
		Where("is_notified = true").Count(&stats.SentNotifications)

	// Unique products
	r.db.WithContext(ctx).Model(&domain.BackInStockSubscription{}).
		Distinct("product_id").Count(&stats.UniqueProducts)

	// Unique customers
	r.db.WithContext(ctx).Model(&domain.BackInStockSubscription{}).
		Distinct("customer_id").Count(&stats.UniqueCustomers)

	return &stats, nil
}

// Admin methods

// ListAll returns all subscriptions with pagination (admin)
func (r *BackInStockRepository) ListAll(ctx context.Context, page, limit int, pendingOnly bool) ([]domain.BackInStockSubscription, int64, error) {
	var subscriptions []domain.BackInStockSubscription
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.BackInStockSubscription{})

	if pendingOnly {
		query = query.Where("is_notified = false")
	}

	// Count total
	query.Count(&total)

	// Get paginated results
	offset := (page - 1) * limit
	err := query.
		Preload("Customer").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&subscriptions).Error

	return subscriptions, total, err
}

// DeleteOldNotified deletes old notified subscriptions (cleanup)
func (r *BackInStockRepository) DeleteOldNotified(ctx context.Context, olderThanDays int) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("is_notified = true AND notification_sent_at < NOW() - INTERVAL '? days'", olderThanDays).
		Delete(&domain.BackInStockSubscription{})
	return result.RowsAffected, result.Error
}
