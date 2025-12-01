package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/niaga-platform/service-customer/internal/middleware"
	"github.com/niaga-platform/service-customer/internal/repository"
	"gorm.io/gorm"
)

// WishlistHandler handles wishlist-related requests
type WishlistHandler struct {
	repo *repository.WishlistRepository
}

// NewWishlistHandler creates a new wishlist handler
func NewWishlistHandler(db *gorm.DB) *WishlistHandler {
	return &WishlistHandler{
		repo: repository.NewWishlistRepository(db),
	}
}

// AddToWishlistRequest represents the request body for adding to wishlist
type AddToWishlistRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
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
		"items": items,
		"count": len(items),
	})
}

// AddToWishlist adds a product to the wishlist
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

	if err := h.repo.Add(c.Request.Context(), userID, req.ProductID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Product added to wishlist",
		"product_id": req.ProductID,
	})
}

// RemoveFromWishlist removes a product from the wishlist
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

	if err := h.repo.Remove(c.Request.Context(), userID, productID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not in wishlist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from wishlist"})
}
