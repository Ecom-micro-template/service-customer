package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/niaga-platform/service-customer/internal/models"
	"github.com/niaga-platform/service-customer/internal/repository"
	"gorm.io/gorm"
)

// MeasurementHandler handles customer measurement-related requests
type MeasurementHandler struct {
	repo *repository.MeasurementRepository
}

// NewMeasurementHandler creates a new measurement handler
func NewMeasurementHandler(db *gorm.DB) *MeasurementHandler {
	return &MeasurementHandler{
		repo: repository.NewMeasurementRepository(db),
	}
}

// CreateMeasurementRequest represents the request body
type CreateMeasurementRequest struct {
	Name          *string  `json:"name"`
	Gender        string   `json:"gender" binding:"required,oneof=men women"`
	Bust          *float64 `json:"bust"`
	Chest         *float64 `json:"chest"`
	Waist         *float64 `json:"waist"`
	Hip           *float64 `json:"hip"`
	ShoulderWidth *float64 `json:"shoulder_width"`
	ArmLength     *float64 `json:"arm_length"`
	Inseam        *float64 `json:"inseam"`
	Outseam       *float64 `json:"outseam"`
	Thigh         *float64 `json:"thigh"`
	Neck          *float64 `json:"neck"`
	Wrist         *float64 `json:"wrist"`
	Height        *float64 `json:"height"`
	Weight        *float64 `json:"weight"`
	Notes         *string  `json:"notes"`
	IsDefault     *bool    `json:"is_default"`
}

// Create handles measurement creation
func (h *MeasurementHandler) Create(c *gin.Context) {
	// TODO: Get user ID from auth context
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isDefault := false
	if req.IsDefault != nil {
		isDefault = *req.IsDefault
	}

	measurement := &models.CustomerMeasurement{
		UserID:        userID,
		Name:          req.Name,
		Gender:        req.Gender,
		Bust:          req.Bust,
		Chest:         req.Chest,
		Waist:         req.Waist,
		Hip:           req.Hip,
		ShoulderWidth: req.ShoulderWidth,
		ArmLength:     req.ArmLength,
		Inseam:        req.Inseam,
		Outseam:       req.Outseam,
		Thigh:         req.Thigh,
		Neck:          req.Neck,
		Wrist:         req.Wrist,
		Height:        req.Height,
		Weight:        req.Weight,
		Notes:         req.Notes,
		IsDefault:     isDefault,
	}

	if err := h.repo.Create(c.Request.Context(), measurement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create measurement"})
		return
	}

	// If set as default, update other measurements
	if isDefault {
		h.repo.SetDefault(c.Request.Context(), userID, measurement.ID)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Measurement created successfully",
		"measurement": measurement,
	})
}

// GetByID retrieves a measurement by ID (with IDOR protection)
func (h *MeasurementHandler) GetByID(c *gin.Context) {
	// Get user ID from auth context for ownership check
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid measurement ID"})
		return
	}

	// IDOR protection: only fetch if owned by user
	measurement, err := h.repo.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve measurement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"measurement": measurement})
}

// List retrieves all measurements for the authenticated user
func (h *MeasurementHandler) List(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	measurements, err := h.repo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve measurements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements,
		"count":        len(measurements),
	})
}

// Update updates a measurement (with IDOR protection)
func (h *MeasurementHandler) Update(c *gin.Context) {
	// Get user ID from auth context for ownership check
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid measurement ID"})
		return
	}

	var req CreateMeasurementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// IDOR protection: only fetch if owned by user
	measurement, err := h.repo.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve measurement"})
		return
	}

	// Update fields
	if req.Name != nil {
		measurement.Name = req.Name
	}
	if req.Gender != "" {
		measurement.Gender = req.Gender
	}
	measurement.Bust = req.Bust
	measurement.Chest = req.Chest
	measurement.Waist = req.Waist
	measurement.Hip = req.Hip
	measurement.ShoulderWidth = req.ShoulderWidth
	measurement.ArmLength = req.ArmLength
	measurement.Inseam = req.Inseam
	measurement.Outseam = req.Outseam
	measurement.Thigh = req.Thigh
	measurement.Neck = req.Neck
	measurement.Wrist = req.Wrist
	measurement.Height = req.Height
	measurement.Weight = req.Weight
	measurement.Notes = req.Notes
	
	if req.IsDefault != nil {
		measurement.IsDefault = *req.IsDefault
	}

	if err := h.repo.Update(c.Request.Context(), measurement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update measurement"})
		return
	}

	// If set as default, update other measurements
	if measurement.IsDefault {
		h.repo.SetDefault(c.Request.Context(), measurement.UserID, measurement.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Measurement updated successfully",
		"measurement": measurement,
	})
}

// Delete deletes a measurement (with IDOR protection)
func (h *MeasurementHandler) Delete(c *gin.Context) {
	// Get user ID from auth context for ownership check
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid measurement ID"})
		return
	}

	// IDOR protection: only delete if owned by user
	if err := h.repo.Delete(c.Request.Context(), id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete measurement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Measurement deleted successfully"})
}

// SetDefault sets a measurement as default
func (h *MeasurementHandler) SetDefault(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid measurement ID"})
		return
	}

	if err := h.repo.SetDefault(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default measurement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default measurement set successfully"})
}
