package wishlist

import (
	"github.com/google/uuid"
)

// WishlistItem represents a product in a wishlist.
// This is a value object - immutable once created.
type WishlistItem struct {
	id           uuid.UUID
	productID    uuid.UUID
	variantID    *uuid.UUID
	variantSKU   string
	variantName  string
	priceAtAdd   float64
	notifyOnSale bool

	// Denormalized product info
	productName  string
	productSlug  string
	productImage string
}

// WishlistItemParams contains parameters for creating a WishlistItem.
type WishlistItemParams struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	VariantID    *uuid.UUID
	VariantSKU   string
	VariantName  string
	PriceAtAdd   float64
	NotifyOnSale bool
	ProductName  string
	ProductSlug  string
	ProductImage string
}

// NewWishlistItem creates a new WishlistItem.
func NewWishlistItem(params WishlistItemParams) WishlistItem {
	id := params.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	return WishlistItem{
		id:           id,
		productID:    params.ProductID,
		variantID:    params.VariantID,
		variantSKU:   params.VariantSKU,
		variantName:  params.VariantName,
		priceAtAdd:   params.PriceAtAdd,
		notifyOnSale: params.NotifyOnSale,
		productName:  params.ProductName,
		productSlug:  params.ProductSlug,
		productImage: params.ProductImage,
	}
}

// Getters
func (i WishlistItem) ID() uuid.UUID         { return i.id }
func (i WishlistItem) ProductID() uuid.UUID  { return i.productID }
func (i WishlistItem) VariantID() *uuid.UUID { return i.variantID }
func (i WishlistItem) VariantSKU() string    { return i.variantSKU }
func (i WishlistItem) VariantName() string   { return i.variantName }
func (i WishlistItem) PriceAtAdd() float64   { return i.priceAtAdd }
func (i WishlistItem) NotifyOnSale() bool    { return i.notifyOnSale }
func (i WishlistItem) ProductName() string   { return i.productName }
func (i WishlistItem) ProductSlug() string   { return i.productSlug }
func (i WishlistItem) ProductImage() string  { return i.productImage }

// HasVariant returns true if this item refers to a specific variant.
func (i WishlistItem) HasVariant() bool {
	return i.variantID != nil
}

// UniqueKey returns a unique key for this item (product + variant).
func (i WishlistItem) UniqueKey() string {
	if i.variantID != nil {
		return i.productID.String() + "-" + i.variantID.String()
	}
	return i.productID.String() + "-nil"
}

// MatchesProduct checks if this item matches a product/variant.
func (i WishlistItem) MatchesProduct(productID uuid.UUID, variantID *uuid.UUID) bool {
	if i.productID != productID {
		return false
	}
	if i.variantID == nil && variantID == nil {
		return true
	}
	if i.variantID == nil || variantID == nil {
		return false
	}
	return *i.variantID == *variantID
}

// WithNotifyOnSale returns a new item with updated notification setting.
func (i WishlistItem) WithNotifyOnSale(notify bool) WishlistItem {
	return WishlistItem{
		id:           i.id,
		productID:    i.productID,
		variantID:    i.variantID,
		variantSKU:   i.variantSKU,
		variantName:  i.variantName,
		priceAtAdd:   i.priceAtAdd,
		notifyOnSale: notify,
		productName:  i.productName,
		productSlug:  i.productSlug,
		productImage: i.productImage,
	}
}
