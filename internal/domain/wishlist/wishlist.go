package wishlist

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Domain errors for Wishlist aggregate
var (
	ErrWishlistNotFound  = errors.New("wishlist not found")
	ErrItemNotFound      = errors.New("item not found in wishlist")
	ErrItemAlreadyExists = errors.New("item already in wishlist")
)

// Wishlist is the aggregate root for customer wishlists.
type Wishlist struct {
	userID    uuid.UUID
	items     []WishlistItem
	updatedAt time.Time
}

// NewWishlist creates a new Wishlist aggregate.
func NewWishlist(userID uuid.UUID) *Wishlist {
	return &Wishlist{
		userID:    userID,
		items:     make([]WishlistItem, 0),
		updatedAt: time.Now(),
	}
}

// Getters
func (w *Wishlist) UserID() uuid.UUID     { return w.userID }
func (w *Wishlist) Items() []WishlistItem { return w.items }
func (w *Wishlist) UpdatedAt() time.Time  { return w.updatedAt }

// ItemCount returns the number of items in the wishlist.
func (w *Wishlist) ItemCount() int {
	return len(w.items)
}

// IsEmpty returns true if the wishlist is empty.
func (w *Wishlist) IsEmpty() bool {
	return len(w.items) == 0
}

// --- Behavior Methods ---

// AddItem adds an item to the wishlist.
func (w *Wishlist) AddItem(params WishlistItemParams) error {
	// Check for duplicates
	for _, item := range w.items {
		if item.MatchesProduct(params.ProductID, params.VariantID) {
			return ErrItemAlreadyExists
		}
	}

	item := NewWishlistItem(params)
	w.items = append(w.items, item)
	w.updatedAt = time.Now()
	return nil
}

// RemoveItem removes an item from the wishlist.
func (w *Wishlist) RemoveItem(productID uuid.UUID, variantID *uuid.UUID) error {
	for i, item := range w.items {
		if item.MatchesProduct(productID, variantID) {
			w.items = append(w.items[:i], w.items[i+1:]...)
			w.updatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

// RemoveItemByID removes an item by its ID.
func (w *Wishlist) RemoveItemByID(itemID uuid.UUID) error {
	for i, item := range w.items {
		if item.ID() == itemID {
			w.items = append(w.items[:i], w.items[i+1:]...)
			w.updatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

// Clear removes all items from the wishlist.
func (w *Wishlist) Clear() {
	w.items = make([]WishlistItem, 0)
	w.updatedAt = time.Now()
}

// ContainsProduct checks if a product is in the wishlist.
func (w *Wishlist) ContainsProduct(productID uuid.UUID, variantID *uuid.UUID) bool {
	for _, item := range w.items {
		if item.MatchesProduct(productID, variantID) {
			return true
		}
	}
	return false
}

// GetItem returns an item by product/variant.
func (w *Wishlist) GetItem(productID uuid.UUID, variantID *uuid.UUID) *WishlistItem {
	for _, item := range w.items {
		if item.MatchesProduct(productID, variantID) {
			return &item
		}
	}
	return nil
}

// SetNotifyOnSale updates notification setting for an item.
func (w *Wishlist) SetNotifyOnSale(productID uuid.UUID, variantID *uuid.UUID, notify bool) error {
	for i, item := range w.items {
		if item.MatchesProduct(productID, variantID) {
			w.items[i] = item.WithNotifyOnSale(notify)
			w.updatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

// ProductIDs returns all product IDs in the wishlist.
func (w *Wishlist) ProductIDs() []uuid.UUID {
	seen := make(map[uuid.UUID]bool)
	var ids []uuid.UUID

	for _, item := range w.items {
		if !seen[item.ProductID()] {
			ids = append(ids, item.ProductID())
			seen[item.ProductID()] = true
		}
	}
	return ids
}

// ItemsForNotification returns items with notification enabled.
func (w *Wishlist) ItemsForNotification() []WishlistItem {
	var items []WishlistItem
	for _, item := range w.items {
		if item.NotifyOnSale() {
			items = append(items, item)
		}
	}
	return items
}
