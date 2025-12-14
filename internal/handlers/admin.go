package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/KilangDesaMurniBatik/service-customer/internal/models"
)

// AdminCustomerHandler handles admin-specific customer operations
type AdminCustomerHandler struct {
	db *gorm.DB
}

// NewAdminCustomerHandler creates a new admin customer handler
func NewAdminCustomerHandler(db *gorm.DB) *AdminCustomerHandler {
	return &AdminCustomerHandler{db: db}
}

// GetCustomers handles GET /api/v1/admin/customers
func (h *AdminCustomerHandler) GetCustomers(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	search := c.Query("search")
	status := c.Query("status")
	segment := c.Query("segment")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Build query
	query := h.db.Model(&models.Customer{})

	// Apply filters
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ? OR phone LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if segment != "" {
		query = query.Where("segment = ?", segment)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Apply pagination and sorting
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Handle special sort cases
	if sortBy == "totalSpent" {
		sortBy = "total_spent"
	} else if sortBy == "orderCount" {
		sortBy = "order_count"
	}

	query = query.Order(sortBy + " " + sortOrder)

	// Fetch customers with statistics
	var customers []struct {
		models.Customer
		OrderCount  int64   `json:"orderCount"`
		TotalSpent  float64 `json:"totalSpent"`
		LastOrderAt *time.Time `json:"lastOrderAt,omitempty"`
	}

	// Use subquery to get customer statistics
	subQuery := h.db.Table("orders").
		Select("customer_id, COUNT(*) as order_count, SUM(total_amount) as total_spent, MAX(created_at) as last_order_at").
		Where("status NOT IN ('cancelled', 'refunded')").
		Group("customer_id")

	if err := query.
		Select("customers.*, COALESCE(stats.order_count, 0) as order_count, COALESCE(stats.total_spent, 0) as total_spent, stats.last_order_at").
		Joins("LEFT JOIN (?) as stats ON customers.id = stats.customer_id", subQuery).
		Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch customers",
			"code":  "DB_ERROR",
		})
		return
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    customers,
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetCustomerStats handles GET /api/v1/admin/customers/stats
func (h *AdminCustomerHandler) GetCustomerStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	// Calculate date range based on period
	var startDate time.Time
	now := time.Now()

	switch period {
	case "day":
		startDate = now.AddDate(0, 0, -1)
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "year":
		startDate = now.AddDate(-1, 0, 0)
	default:
		startDate = now.AddDate(0, -1, 0)
	}

	// Get current period stats
	var currentStats struct {
		TotalCustomers   int64 `gorm:"column:total_customers"`
		NewCustomers     int64 `gorm:"column:new_customers"`
		ActiveCustomers  int64 `gorm:"column:active_customers"`
		ReturningCustomers int64 `gorm:"column:returning_customers"`
	}

	h.db.Model(&models.Customer{}).
		Select(
			"COUNT(*) as total_customers",
			"COUNT(CASE WHEN created_at >= ? THEN 1 END) as new_customers",
			"COUNT(CASE WHEN status = 'active' THEN 1 END) as active_customers",
		).
		Where("created_at <= ?", now).
		Scan(&currentStats, startDate)

	// Get customers with multiple orders
	h.db.Table("customers").
		Select("COUNT(DISTINCT customer_id)").
		Joins("JOIN orders ON customers.id = orders.customer_id").
		Where("orders.created_at >= ?", startDate).
		Group("customer_id").
		Having("COUNT(*) > 1").
		Count(&currentStats.ReturningCustomers)

	// Get previous period stats for comparison
	prevStartDate := startDate
	switch period {
	case "day":
		prevStartDate = startDate.AddDate(0, 0, -1)
		startDate = startDate.AddDate(0, 0, -1)
	case "week":
		prevStartDate = startDate.AddDate(0, 0, -7)
		startDate = startDate.AddDate(0, 0, -7)
	case "month":
		prevStartDate = startDate.AddDate(0, -1, 0)
		startDate = startDate.AddDate(0, -1, 0)
	case "year":
		prevStartDate = startDate.AddDate(-1, 0, 0)
		startDate = startDate.AddDate(-1, 0, 0)
	}

	var previousStats struct {
		NewCustomers int64 `gorm:"column:new_customers"`
	}

	h.db.Model(&models.Customer{}).
		Select("COUNT(*) as new_customers").
		Where("created_at >= ? AND created_at < ?", prevStartDate, startDate).
		Scan(&previousStats)

	// Calculate percentage change
	customerChange := float64(0)
	if previousStats.NewCustomers > 0 {
		customerChange = ((float64(currentStats.NewCustomers) - float64(previousStats.NewCustomers)) / float64(previousStats.NewCustomers)) * 100
	}

	// Get top customers
	var topCustomers []struct {
		models.Customer
		TotalSpent float64 `json:"totalSpent"`
		OrderCount int64   `json:"orderCount"`
	}

	h.db.Table("customers").
		Select("customers.*, COALESCE(SUM(orders.total_amount), 0) as total_spent, COUNT(orders.id) as order_count").
		Joins("LEFT JOIN orders ON customers.id = orders.customer_id").
		Where("orders.status NOT IN ('cancelled', 'refunded')").
		Group("customers.id").
		Order("total_spent DESC").
		Limit(5).
		Find(&topCustomers)

	// Get customer segments
	var segments []struct {
		Segment string `json:"segment"`
		Count   int64  `json:"count"`
	}

	h.db.Model(&models.Customer{}).
		Select("segment, COUNT(*) as count").
		Group("segment").
		Find(&segments)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"current": gin.H{
				"totalCustomers":     currentStats.TotalCustomers,
				"newCustomers":       currentStats.NewCustomers,
				"activeCustomers":    currentStats.ActiveCustomers,
				"returningCustomers": currentStats.ReturningCustomers,
			},
			"changes": gin.H{
				"customerChange": customerChange,
			},
			"topCustomers": topCustomers,
			"segments":     segments,
		},
	})
}

// GetCustomerByID handles GET /api/v1/admin/customers/:id
func (h *AdminCustomerHandler) GetCustomerByID(c *gin.Context) {
	customerID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer ID",
			"code":  "INVALID_ID",
		})
		return
	}

	var customer models.Customer
	if err := h.db.First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Customer not found",
				"code":  "NOT_FOUND",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch customer",
			"code":  "DB_ERROR",
		})
		return
	}

	// Get customer statistics
	var stats struct {
		OrderCount     int64     `json:"orderCount"`
		TotalSpent     float64   `json:"totalSpent"`
		AverageOrder   float64   `json:"averageOrder"`
		LastOrderAt    *time.Time `json:"lastOrderAt,omitempty"`
		FirstOrderAt   *time.Time `json:"firstOrderAt,omitempty"`
	}

	h.db.Table("orders").
		Select(
			"COUNT(*) as order_count",
			"COALESCE(SUM(total_amount), 0) as total_spent",
			"COALESCE(AVG(total_amount), 0) as average_order",
			"MAX(created_at) as last_order_at",
			"MIN(created_at) as first_order_at",
		).
		Where("customer_id = ? AND status NOT IN ('cancelled', 'refunded')", id).
		Scan(&stats)

	// Get addresses
	var addresses []models.Address
	h.db.Where("customer_id = ?", id).Find(&addresses)

	// Get wishlist items count
	var wishlistCount int64
	h.db.Model(&models.WishlistItem{}).Where("customer_id = ?", id).Count(&wishlistCount)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"customer":      customer,
			"statistics":    stats,
			"addresses":     addresses,
			"wishlistCount": wishlistCount,
		},
	})
}

// GetCustomerOrders handles GET /api/v1/admin/customers/:id/orders
func (h *AdminCustomerHandler) GetCustomerOrders(c *gin.Context) {
	customerID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Parse UUID
	id, err := uuid.Parse(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer ID",
			"code":  "INVALID_ID",
		})
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Check if customer exists
	var customer models.Customer
	if err := h.db.First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Customer not found",
				"code":  "NOT_FOUND",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch customer",
			"code":  "DB_ERROR",
		})
		return
	}

	// Fetch orders
	var orders []struct {
		ID          uuid.UUID `json:"id"`
		OrderNumber string    `json:"orderNumber"`
		Status      string    `json:"status"`
		TotalAmount float64   `json:"totalAmount"`
		ItemCount   int       `json:"itemCount"`
		CreatedAt   time.Time `json:"createdAt"`
	}

	query := h.db.Table("orders").
		Select("orders.id, orders.order_number, orders.status, orders.total_amount, orders.created_at, COUNT(order_items.id) as item_count").
		Joins("LEFT JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.customer_id = ?", id).
		Group("orders.id")

	// Count total
	var total int64
	query.Count(&total)

	// Apply pagination
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("orders.created_at DESC")

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch orders",
			"code":  "DB_ERROR",
		})
		return
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"customer": gin.H{
				"id":    customer.ID,
				"name":  customer.Name,
				"email": customer.Email,
			},
			"orders": orders,
		},
		"meta": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// UpdateCustomerStatus handles PUT /api/v1/admin/customers/:id/status
func (h *AdminCustomerHandler) UpdateCustomerStatus(c *gin.Context) {
	customerID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer ID",
			"code":  "INVALID_ID",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive suspended"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	// Update customer status
	var customer models.Customer
	if err := h.db.First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Customer not found",
				"code":  "NOT_FOUND",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch customer",
			"code":  "DB_ERROR",
		})
		return
	}

	oldStatus := customer.Status
	customer.Status = req.Status

	if err := h.db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update customer status",
			"code":  "DB_ERROR",
		})
		return
	}

	// Log the status change
	auditLog := models.AuditLog{
		ID:         uuid.New(),
		EntityType: "customer",
		EntityID:   customer.ID,
		Action:     "status_change",
		UserID:     c.GetString("user_id"),
		Details: map[string]interface{}{
			"old_status": oldStatus,
			"new_status": req.Status,
			"reason":     req.Reason,
		},
		CreatedAt: time.Now(),
	}

	h.db.Create(&auditLog)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"customer":  customer,
			"oldStatus": oldStatus,
			"newStatus": req.Status,
		},
		"message": "Customer status updated successfully",
	})
}

// UpdateCustomerSegment handles PUT /api/v1/admin/customers/:id/segment
func (h *AdminCustomerHandler) UpdateCustomerSegment(c *gin.Context) {
	customerID := c.Param("id")

	// Parse UUID
	id, err := uuid.Parse(customerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid customer ID",
			"code":  "INVALID_ID",
		})
		return
	}

	var req struct {
		Segment string `json:"segment" binding:"required,oneof=vip gold silver bronze regular"`
		Note    string `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	// Update customer segment
	var customer models.Customer
	if err := h.db.First(&customer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Customer not found",
				"code":  "NOT_FOUND",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch customer",
			"code":  "DB_ERROR",
		})
		return
	}

	oldSegment := customer.Segment
	customer.Segment = req.Segment

	if err := h.db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update customer segment",
			"code":  "DB_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"customer":   customer,
			"oldSegment": oldSegment,
			"newSegment": req.Segment,
		},
		"message": "Customer segment updated successfully",
	})
}

// BulkUpdateCustomers handles POST /api/v1/admin/customers/bulk
func (h *AdminCustomerHandler) BulkUpdateCustomers(c *gin.Context) {
	var req struct {
		CustomerIDs []string `json:"customer_ids" binding:"required,min=1"`
		Action      string   `json:"action" binding:"required,oneof=updateStatus updateSegment export"`
		Status      string   `json:"status"`
		Segment     string   `json:"segment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	// Convert string IDs to UUIDs
	var customerIDs []uuid.UUID
	for _, idStr := range req.CustomerIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid customer ID: " + idStr,
				"code":  "INVALID_ID",
			})
			return
		}
		customerIDs = append(customerIDs, id)
	}

	switch req.Action {
	case "updateStatus":
		if req.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Status is required for updateStatus action",
				"code":  "MISSING_STATUS",
			})
			return
		}

		result := h.db.Model(&models.Customer{}).
			Where("id IN ?", customerIDs).
			Update("status", req.Status)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update customers",
				"code":  "DB_ERROR",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"updated": result.RowsAffected,
			},
			"message": "Customers updated successfully",
		})

	case "updateSegment":
		if req.Segment == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Segment is required for updateSegment action",
				"code":  "MISSING_SEGMENT",
			})
			return
		}

		result := h.db.Model(&models.Customer{}).
			Where("id IN ?", customerIDs).
			Update("segment", req.Segment)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update customers",
				"code":  "DB_ERROR",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"updated": result.RowsAffected,
			},
			"message": "Customer segments updated successfully",
		})

	case "export":
		var customers []models.Customer
		h.db.Where("id IN ?", customerIDs).Find(&customers)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"customers": customers,
			},
			"message": "Customers exported successfully",
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid action",
			"code":  "INVALID_ACTION",
		})
	}
}