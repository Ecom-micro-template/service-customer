package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/middleware"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"github.com/Ecom-micro-template/service-customer/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// AddressHandler handles address-related requests
type AddressHandler struct {
	repo *repository.AddressRepository
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(db *gorm.DB) *AddressHandler {
	return &AddressHandler{
		repo: repository.NewAddressRepository(db),
	}
}

// CreateAddressRequest represents the request body for creating an address
type CreateAddressRequest struct {
	Label         string `json:"label" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	AddressLine1  string `json:"address_line1" binding:"required"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city" binding:"required"`
	State         string `json:"state" binding:"required"`
	Postcode      string `json:"postcode" binding:"required"`
	Country       string `json:"country" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}

// UpdateAddressRequest represents the request body for updating an address
type UpdateAddressRequest struct {
	Label         string `json:"label"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Postcode      string `json:"postcode"`
	Country       string `json:"country"`
	IsDefault     *bool  `json:"is_default"`
}

// ListAddresses retrieves all addresses for the customer
// GET /api/v1/customer/addresses
func (h *AddressHandler) ListAddresses(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	addresses, err := h.repo.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"addresses": addresses,
		"count":     len(addresses),
	})
}

// CreateAddress creates a new address
// POST /api/v1/customer/addresses
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := &models.Address{
		UserID:        userID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		City:          req.City,
		State:         req.State,
		Postcode:      req.Postcode,
		Country:       req.Country,
		IsDefault:     req.IsDefault,
	}

	if err := h.repo.Create(c.Request.Context(), address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Address created successfully",
		"address": address,
	})
}

// UpdateAddress updates an existing address
// PUT /api/v1/customer/addresses/:id
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing address
	address, err := h.repo.GetByID(c.Request.Context(), addressID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve address"})
		return
	}

	// Update fields
	if req.Label != "" {
		address.Label = req.Label
	}
	if req.RecipientName != "" {
		address.RecipientName = req.RecipientName
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.AddressLine1 != "" {
		address.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != "" {
		address.AddressLine2 = req.AddressLine2
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.State != "" {
		address.State = req.State
	}
	if req.Postcode != "" {
		address.Postcode = req.Postcode
	}
	if req.Country != "" {
		address.Country = req.Country
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	if err := h.repo.Update(c.Request.Context(), address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address updated successfully",
		"address": address,
	})
}

// DeleteAddress deletes an address
// DELETE /api/v1/customer/addresses/:id
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), addressID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

// SetDefaultAddress sets an address as the default
// PUT /api/v1/customer/addresses/:id/default
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	if err := h.repo.SetDefault(c.Request.Context(), addressID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default address set successfully"})
}
