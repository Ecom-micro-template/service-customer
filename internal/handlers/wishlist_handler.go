package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/middleware"
	"github.com/Ecom-micro-template/service-customer/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// WishlistHandler handles wishlist-related requests
type WishlistHandler struct {
	repo *persistence.WishlistRepository
}

// NewWishlistHandler creates a new wishlist handler
func NewWishlistHandler(db *gorm.DB) *WishlistHandler {
	return &WishlistHandler{
		repo: persistence.NewWishlistRepository(db),
	}
}

// AddToWishlistRequest represents the request body for adding to wishlist
type AddToWishlistRequest struct {
	ProductID    uuid.UUID  `json:"product_id" binding:"required"`
	VariantID    *uuid.UUID `json:"variant_id,omitempty"`
	VariantSKU   *string    `json:"variant_sku,omitempty"`
	VariantName  *string    `json:"variant_name,omitempty"` // e.g., "Red / Large"
	PriceAtAdd   float64    `json:"price_at_add,omitempty"`
	NotifyOnSale *bool      `json:"notify_on_sale,omitempty"`
	ProductName  *string    `json:"product_name,omitempty"`
	ProductSlug  *string    `json:"product_slug,omitempty"`
	ProductImage *string    `json:"product_image,omitempty"`
}

// UpdateWishlistItemRequest represents the request body for updating a wishlist item
type UpdateWishlistItemRequest struct {
	NotifyOnSale *bool `json:"notify_on_sale"`
}

// GetWishlist retrieves the customer's wishlist
// GET /api/v1/customer/wishlist
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	items, err := h.repo.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": items,
			"count": len(items),
		},
	})
}

// AddToWishlist adds a product/variant to the wishlist
// POST /api/v1/customer/wishlist
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notifyOnSale := false
	if req.NotifyOnSale != nil {
		notifyOnSale = *req.NotifyOnSale
	}

	input := persistence.AddWishlistItemInput{
		ProductID:    req.ProductID,
		VariantID:    req.VariantID,
		VariantSKU:   req.VariantSKU,
		VariantName:  req.VariantName,
		PriceAtAdd:   req.PriceAtAdd,
		NotifyOnSale: notifyOnSale,
		ProductName:  req.ProductName,
		ProductSlug:  req.ProductSlug,
		ProductImage: req.ProductImage,
	}

	if err := h.repo.AddWithVariant(c.Request.Context(), userID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"message":    "Added to wishlist",
		"product_id": req.ProductID,
		"variant_id": req.VariantID,
	})
}

// RemoveFromWishlist removes a product from the wishlist (all variants)
// DELETE /api/v1/customer/wishlist/:productId
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Check if variant_id is provided as query parameter
	variantIDStr := c.Query("variant_id")
	var variantID *uuid.UUID
	if variantIDStr != "" {
		parsed, err := uuid.Parse(variantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
			return
		}
		variantID = &parsed
	}

	var removeErr error
	if variantID != nil {
		removeErr = h.repo.RemoveWithVariant(c.Request.Context(), userID, productID, variantID)
	} else {
		removeErr = h.repo.Remove(c.Request.Context(), userID, productID)
	}

	if removeErr != nil {
		if removeErr == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not in wishlist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Removed from wishlist",
	})
}

// RemoveWishlistItem removes a wishlist item by ID
// DELETE /api/v1/customer/wishlist/items/:itemId
func (h *WishlistHandler) RemoveWishlistItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := h.repo.RemoveByID(c.Request.Context(), userID, itemID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item removed from wishlist",
	})
}

// UpdateWishlistItem updates a wishlist item (e.g., notify_on_sale)
// PATCH /api/v1/customer/wishlist/items/:itemId
func (h *WishlistHandler) UpdateWishlistItem(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req UpdateWishlistItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NotifyOnSale != nil {
		if err := h.repo.UpdateNotifyOnSale(c.Request.Context(), userID, itemID, *req.NotifyOnSale); err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Wishlist item updated",
	})
}

// CheckWishlist checks if a product/variant is in the wishlist
// GET /api/v1/customer/wishlist/check/:productId
func (h *WishlistHandler) CheckWishlist(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Check if variant_id is provided as query parameter
	variantIDStr := c.Query("variant_id")
	var variantID *uuid.UUID
	if variantIDStr != "" {
		parsed, err := uuid.Parse(variantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
			return
		}
		variantID = &parsed
	}

	var exists bool
	var checkErr error
	if variantID != nil {
		exists, checkErr = h.repo.ExistsWithVariant(c.Request.Context(), userID, productID, variantID)
	} else {
		exists, checkErr = h.repo.Exists(c.Request.Context(), userID, productID)
	}

	if checkErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"in_wishlist": exists,
		"product_id": productID,
		"variant_id": variantID,
	})
}

// GetWishlistCount returns the count of items in the wishlist
// GET /api/v1/customer/wishlist/count
func (h *WishlistHandler) GetWishlistCount(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	count, err := h.repo.CountByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   count,
	})
}
