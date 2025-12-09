package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/niaga-platform/lib-common/response"
	"github.com/niaga-platform/service-customer/internal/models"
	"github.com/niaga-platform/service-customer/internal/repository"
	"go.uber.org/zap"
)

type AdminCustomerHandler struct {
	customerRepo repository.CustomerRepository
	logger       *zap.Logger
}

func NewAdminCustomerHandler(customerRepo repository.CustomerRepository, logger *zap.Logger) *AdminCustomerHandler {
	return &AdminCustomerHandler{
		customerRepo: customerRepo,
		logger:       logger,
	}
}

// CustomerListFilter represents filters for admin customer listing
type CustomerListFilter struct {
	Status    string     `form:"status"`
	Segment   string     `form:"segment"`
	DateFrom  *time.Time `form:"date_from"`
	DateTo    *time.Time `form:"date_to"`
	OrdersMin *int       `form:"orders_min"`
	OrdersMax *int       `form:"orders_max"`
	SpentMin  *float64   `form:"spent_min"`
	SpentMax  *float64   `form:"spent_max"`
	Search    string     `form:"search"`
	Page      int        `form:"page"`
	Limit     int        `form:"limit"`
	SortBy    string     `form:"sort_by"`
	SortOrder string     `form:"sort_order"`
}

// GetCustomers handles GET /admin/customers
func (h *AdminCustomerHandler) GetCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := CustomerListFilter{
		Status:    c.Query("status"),
		Segment:   c.Query("segment"),
		Search:    c.Query("search"),
		Page:      page,
		Limit:     limit,
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	// Parse date filters
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			filter.DateFrom = &dateFrom
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = dateTo.Add(24*time.Hour - time.Second)
			filter.DateTo = &dateTo
		}
	}

	// Parse order count filters
	if ordersMinStr := c.Query("orders_min"); ordersMinStr != "" {
		if ordersMin, err := strconv.Atoi(ordersMinStr); err == nil {
			filter.OrdersMin = &ordersMin
		}
	}
	if ordersMaxStr := c.Query("orders_max"); ordersMaxStr != "" {
		if ordersMax, err := strconv.Atoi(ordersMaxStr); err == nil {
			filter.OrdersMax = &ordersMax
		}
	}

	// Parse spending filters
	if spentMinStr := c.Query("spent_min"); spentMinStr != "" {
		if spentMin, err := strconv.ParseFloat(spentMinStr, 64); err == nil {
			filter.SpentMin = &spentMin
		}
	}
	if spentMaxStr := c.Query("spent_max"); spentMaxStr != "" {
		if spentMax, err := strconv.ParseFloat(spentMaxStr, 64); err == nil {
			filter.SpentMax = &spentMax
		}
	}

	customers, total, err := h.customerRepo.ListAdmin(filter)
	if err != nil {
		h.logger.Error("Failed to list customers", zap.Error(err))
		response.InternalServerError(c, "Failed to retrieve customers")
		return
	}

	response.Paginated(c, customers, page, limit, total)
}

// GetCustomer handles GET /admin/customers/:id
func (h *AdminCustomerHandler) GetCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	customer, err := h.customerRepo.GetByID(customerID)
	if err != nil {
		h.logger.Error("Failed to get customer", zap.Error(err))
		response.NotFound(c, "Customer not found")
		return
	}

	response.OK(c, "Customer retrieved", customer)
}

// CreateCustomer handles POST /admin/customers
func (h *AdminCustomerHandler) CreateCustomer(c *gin.Context) {
	var req models.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	// Get admin user ID
	var createdBy *uuid.UUID
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			createdBy = &uid
		}
	}

	customer, err := h.customerRepo.Create(&req, createdBy)
	if err != nil {
		h.logger.Error("Failed to create customer", zap.Error(err))
		response.InternalServerError(c, "Failed to create customer")
		return
	}

	response.Created(c, "Customer created successfully", customer)
}

// UpdateCustomer handles PUT /admin/customers/:id
func (h *AdminCustomerHandler) UpdateCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	var req models.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	customer, err := h.customerRepo.Update(customerID, &req)
	if err != nil {
		h.logger.Error("Failed to update customer", zap.Error(err))
		response.InternalServerError(c, "Failed to update customer")
		return
	}

	response.Updated(c, "Customer updated successfully", customer)
}

// DeleteCustomer handles DELETE /admin/customers/:id
func (h *AdminCustomerHandler) DeleteCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	if err := h.customerRepo.Delete(customerID); err != nil {
		h.logger.Error("Failed to delete customer", zap.Error(err))
		response.InternalServerError(c, "Failed to delete customer")
		return
	}

	response.Deleted(c, "Customer deleted successfully")
}

// GetCustomerOrders handles GET /admin/customers/:id/orders
func (h *AdminCustomerHandler) GetCustomerOrders(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	orders, total, err := h.customerRepo.GetCustomerOrders(customerID, page, limit)
	if err != nil {
		h.logger.Error("Failed to get customer orders", zap.Error(err))
		response.InternalServerError(c, "Failed to retrieve customer orders")
		return
	}

	response.Paginated(c, orders, page, limit, total)
}

// AddCustomerNote handles POST /admin/customers/:id/notes
func (h *AdminCustomerHandler) AddCustomerNote(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	var req struct {
		Note      string `json:"note" binding:"required"`
		IsPrivate bool   `json:"is_private"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	// Get admin user ID
	var createdBy uuid.UUID
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			createdBy = uid
		}
	}

	note, err := h.customerRepo.AddNote(customerID, req.Note, req.IsPrivate, createdBy)
	if err != nil {
		h.logger.Error("Failed to add customer note", zap.Error(err))
		response.InternalServerError(c, "Failed to add customer note")
		return
	}

	response.Created(c, "Customer note added successfully", note)
}

// GetCustomerNotes handles GET /admin/customers/:id/notes
func (h *AdminCustomerHandler) GetCustomerNotes(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	notes, err := h.customerRepo.GetNotes(customerID)
	if err != nil {
		h.logger.Error("Failed to get customer notes", zap.Error(err))
		response.InternalServerError(c, "Failed to retrieve customer notes")
		return
	}

	response.OK(c, "Customer notes retrieved", notes)
}

// GetCustomerActivity handles GET /admin/customers/:id/activity
func (h *AdminCustomerHandler) GetCustomerActivity(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	activity, total, err := h.customerRepo.GetActivity(customerID, page, limit)
	if err != nil {
		h.logger.Error("Failed to get customer activity", zap.Error(err))
		response.InternalServerError(c, "Failed to retrieve customer activity")
		return
	}

	response.Paginated(c, activity, page, limit, total)
}

// GetSegments handles GET /admin/segments
func (h *AdminCustomerHandler) GetSegments(c *gin.Context) {
	segments, err := h.customerRepo.GetSegments()
	if err != nil {
		h.logger.Error("Failed to get segments", zap.Error(err))
		response.InternalServerError(c, "Failed to retrieve customer segments")
		return
	}

	response.OK(c, "Customer segments retrieved", segments)
}

// CreateSegment handles POST /admin/segments
func (h *AdminCustomerHandler) CreateSegment(c *gin.Context) {
	var req struct {
		Name        string      `json:"name" binding:"required"`
		Description string      `json:"description"`
		Conditions  interface{} `json:"conditions"` // JSON conditions for dynamic segments
		Color       string      `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	segment, err := h.customerRepo.CreateSegment(req.Name, req.Description, req.Conditions, req.Color)
	if err != nil {
		h.logger.Error("Failed to create segment", zap.Error(err))
		response.InternalServerError(c, "Failed to create customer segment")
		return
	}

	response.Created(c, "Customer segment created successfully", segment)
}

// UpdateSegment handles PUT /admin/segments/:id
func (h *AdminCustomerHandler) UpdateSegment(c *gin.Context) {
	segmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid segment ID", nil)
		return
	}

	var req struct {
		Name        *string     `json:"name,omitempty"`
		Description *string     `json:"description,omitempty"`
		Conditions  interface{} `json:"conditions,omitempty"`
		Color       *string     `json:"color,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	segment, err := h.customerRepo.UpdateSegment(segmentID, req.Name, req.Description, req.Conditions, req.Color)
	if err != nil {
		h.logger.Error("Failed to update segment", zap.Error(err))
		response.InternalServerError(c, "Failed to update customer segment")
		return
	}

	response.Updated(c, "Customer segment updated successfully", segment)
}

// DeleteSegment handles DELETE /admin/segments/:id
func (h *AdminCustomerHandler) DeleteSegment(c *gin.Context) {
	segmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid segment ID", nil)
		return
	}

	if err := h.customerRepo.DeleteSegment(segmentID); err != nil {
		h.logger.Error("Failed to delete segment", zap.Error(err))
		response.InternalServerError(c, "Failed to delete customer segment")
		return
	}

	response.Deleted(c, "Customer segment deleted successfully")
}

// AssignSegment handles POST /admin/customers/:id/segments
func (h *AdminCustomerHandler) AssignSegment(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid customer ID", nil)
		return
	}

	var req struct {
		SegmentIDs []uuid.UUID `json:"segment_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	if err := h.customerRepo.AssignSegments(customerID, req.SegmentIDs); err != nil {
		h.logger.Error("Failed to assign segments", zap.Error(err))
		response.InternalServerError(c, "Failed to assign customer segments")
		return
	}

	response.OK(c, "Customer segments assigned successfully", nil)
}

// ExportCustomers handles GET /admin/customers/export
func (h *AdminCustomerHandler) ExportCustomers(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")

	filter := CustomerListFilter{
		Status:  c.Query("status"),
		Segment: c.Query("segment"),
		Search:  c.Query("search"),
	}

	data, err := h.customerRepo.Export(filter, format)
	if err != nil {
		h.logger.Error("Failed to export customers", zap.Error(err))
		response.InternalServerError(c, "Failed to export customers")
		return
	}

	response.OK(c, "Customers exported successfully", data)
}

// GetCustomerStats handles GET /admin/customers/stats
func (h *AdminCustomerHandler) GetCustomerStats(c *gin.Context) {
	stats, err := h.customerRepo.GetStats()
	if err != nil {
		h.logger.Error("Failed to get customer stats", zap.Error(err))
		response.InternalServerError(c, "Failed to retrieve customer statistics")
		return
	}

	response.OK(c, "Customer statistics retrieved", stats)
}
