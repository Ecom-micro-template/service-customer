package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/middleware"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"github.com/Ecom-micro-template/service-customer/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// HI-001: Back-in-Stock Handler

// BackInStockHandler handles back-in-stock subscription requests
type BackInStockHandler struct {
	repo *repository.BackInStockRepository
}

// NewBackInStockHandler creates a new back-in-stock handler
func NewBackInStockHandler(db *gorm.DB) *BackInStockHandler {
	return &BackInStockHandler{
		repo: repository.NewBackInStockRepository(db),
	}
}

// Subscribe subscribes a customer to back-in-stock notifications
// POST /api/v1/customer/back-in-stock
func (h *BackInStockHandler) Subscribe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var input models.BackInStockSubscribeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.repo.Subscribe(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Subscribed to back-in-stock notification",
		"data":    subscription,
	})
}

// Unsubscribe removes a subscription by product/variant
// DELETE /api/v1/customer/back-in-stock/:productId
func (h *BackInStockHandler) Unsubscribe(c *gin.Context) {
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

	// Optional variant_id
	var variantID *uuid.UUID
	if variantIDStr := c.Query("variant_id"); variantIDStr != "" {
		parsed, err := uuid.Parse(variantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
			return
		}
		variantID = &parsed
	}

	if err := h.repo.Unsubscribe(c.Request.Context(), userID, productID, variantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Unsubscribed from back-in-stock notification",
	})
}

// UnsubscribeByID removes a subscription by ID
// DELETE /api/v1/customer/back-in-stock/subscriptions/:id
func (h *BackInStockHandler) UnsubscribeByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	subscriptionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := h.repo.UnsubscribeByID(c.Request.Context(), userID, subscriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Subscription removed",
	})
}

// GetSubscriptions returns all subscriptions for the current customer
// GET /api/v1/customer/back-in-stock
func (h *BackInStockHandler) GetSubscriptions(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	subscriptions, err := h.repo.GetByCustomer(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"subscriptions": subscriptions,
			"count":         len(subscriptions),
		},
	})
}

// IsSubscribed checks if customer is subscribed to a product
// GET /api/v1/customer/back-in-stock/check/:productId
func (h *BackInStockHandler) IsSubscribed(c *gin.Context) {
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

	// Optional variant_id
	var variantID *uuid.UUID
	if variantIDStr := c.Query("variant_id"); variantIDStr != "" {
		parsed, err := uuid.Parse(variantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
			return
		}
		variantID = &parsed
	}

	subscribed, err := h.repo.IsSubscribed(c.Request.Context(), userID, productID, variantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"subscribed": subscribed,
		"product_id": productID,
		"variant_id": variantID,
	})
}

// Admin Handler

// AdminBackInStockHandler handles admin back-in-stock operations
type AdminBackInStockHandler struct {
	repo *repository.BackInStockRepository
}

// NewAdminBackInStockHandler creates a new admin handler
func NewAdminBackInStockHandler(db *gorm.DB) *AdminBackInStockHandler {
	return &AdminBackInStockHandler{
		repo: repository.NewBackInStockRepository(db),
	}
}

// GetStats returns subscription statistics
// GET /api/v1/admin/back-in-stock/stats
func (h *AdminBackInStockHandler) GetStats(c *gin.Context) {
	stats, err := h.repo.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ListSubscriptions returns all subscriptions with pagination
// GET /api/v1/admin/back-in-stock/subscriptions
func (h *AdminBackInStockHandler) ListSubscriptions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	pendingOnly := c.Query("pending_only") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	subscriptions, total, err := h.repo.ListAll(c.Request.Context(), page, limit, pendingOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions"})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"subscriptions": subscriptions,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// GetByProduct returns subscriptions for a specific product
// GET /api/v1/admin/back-in-stock/products/:productId/subscriptions
func (h *AdminBackInStockHandler) GetByProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Optional variant_id
	var variantID *uuid.UUID
	if variantIDStr := c.Query("variant_id"); variantIDStr != "" {
		parsed, err := uuid.Parse(variantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant ID"})
			return
		}
		variantID = &parsed
	}

	subscriptions, err := h.repo.GetByProduct(c.Request.Context(), productID, variantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"subscriptions": subscriptions,
			"count":         len(subscriptions),
			"product_id":    productID,
			"variant_id":    variantID,
		},
	})
}

// MarkAsNotified marks subscriptions as notified (after sending notifications)
// POST /api/v1/admin/back-in-stock/mark-notified
func (h *AdminBackInStockHandler) MarkAsNotified(c *gin.Context) {
	var req struct {
		SubscriptionIDs []string `json:"subscription_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ids []uuid.UUID
	for _, idStr := range req.SubscriptionIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID: " + idStr})
			return
		}
		ids = append(ids, id)
	}

	if err := h.repo.MarkMultipleAsNotified(c.Request.Context(), ids); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as notified"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Subscriptions marked as notified",
		"count":   len(ids),
	})
}

// Cleanup deletes old notified subscriptions
// DELETE /api/v1/admin/back-in-stock/cleanup
func (h *AdminBackInStockHandler) Cleanup(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("older_than_days", "30"))
	if days < 1 {
		days = 30
	}

	deleted, err := h.repo.DeleteOldNotified(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cleanup completed",
		"deleted": deleted,
	})
}
