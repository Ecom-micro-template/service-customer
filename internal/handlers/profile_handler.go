package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niaga-platform/service-customer/internal/middleware"
	"github.com/niaga-platform/service-customer/internal/models"
	"github.com/niaga-platform/service-customer/internal/repository"
	"gorm.io/gorm"
)

// ProfileHandler handles profile-related requests
type ProfileHandler struct {
	repo *repository.ProfileRepository
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{
		repo: repository.NewProfileRepository(db),
	}
}

// UpdateProfileRequest represents the request body for updating profile
type UpdateProfileRequest struct {
	FullName       string     `json:"full_name"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Gender         string     `json:"gender"`
	ProfilePicture string     `json:"profile_picture"`
}

// GetProfile retrieves the customer's profile
// GET /api/v1/customer/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	profile, err := h.repo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return empty profile if not found
			c.JSON(http.StatusOK, gin.H{
				"profile": gin.H{
					"id":    userID,
					"email": "",
				},
				"message": "Profile not found, please update your profile",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// UpdateProfile creates or updates the customer's profile
// PUT /api/v1/customer/profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing profile or create new one
	profile, err := h.repo.GetByUserID(c.Request.Context(), userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	// Create new profile if doesn't exist
	if profile == nil {
		profile = &models.Profile{
			ID: userID,
		}
	}

	// Update fields
	if req.FullName != "" {
		profile.FullName = req.FullName
	}
	if req.Email != "" {
		profile.Email = req.Email
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		profile.Gender = req.Gender
	}
	if req.ProfilePicture != "" {
		profile.ProfilePicture = req.ProfilePicture
	}

	// Upsert profile
	if err := h.repo.Upsert(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}
