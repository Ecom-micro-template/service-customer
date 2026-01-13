package repository

import (
	"github.com/google/uuid"
	"github.com/Ecom-micro-template/service-customer/internal/domain"
	"gorm.io/gorm"
)

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	// CRUD operations
	ListAdmin(filter models.CustomerListFilter) ([]models.Customer, int64, error)
	GetByID(id uuid.UUID) (*models.Customer, error)
	Create(req *models.CreateCustomerRequest, createdBy *uuid.UUID) (*models.Customer, error)
	Update(id uuid.UUID, req *models.UpdateCustomerRequest) (*models.Customer, error)
	Delete(id uuid.UUID) error

	// Order-related
	GetCustomerOrders(customerID uuid.UUID, page, limit int) ([]CustomerOrderSummary, int64, error)

	// Notes
	AddNote(customerID uuid.UUID, note string, isPrivate bool, createdBy uuid.UUID) (*models.CustomerNote, error)
	GetNotes(customerID uuid.UUID) ([]models.CustomerNote, error)

	// Activity
	GetActivity(customerID uuid.UUID, page, limit int) ([]models.CustomerActivity, int64, error)

	// Segments
	GetSegments() ([]models.CustomerSegment, error)
	CreateSegment(name, description string, conditions interface{}, color string) (*models.CustomerSegment, error)
	UpdateSegment(id uuid.UUID, name, description *string, conditions interface{}, color *string) (*models.CustomerSegment, error)
	DeleteSegment(id uuid.UUID) error
	AssignSegments(customerID uuid.UUID, segmentIDs []uuid.UUID) error

	// Export and stats
	Export(filter models.CustomerListFilter, format string) (interface{}, error)
	GetStats() (*CustomerStats, error)
}

// CustomerOrderItem represents an item in a customer order
type CustomerOrderItem struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	SKU         string    `json:"sku"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Total       float64   `json:"total"`
	ImageURL    string    `json:"image_url"`
}

// CustomerOrderSummary represents a summarized order for a customer
type CustomerOrderSummary struct {
	ID            uuid.UUID           `json:"id"`
	OrderNum      string              `json:"order_number"`
	Total         float64             `json:"total"`
	Subtotal      float64             `json:"subtotal"`
	Status        string              `json:"status"`
	PaymentStatus string              `json:"payment_status"`
	Items         []CustomerOrderItem `json:"items"`
	CreatedAt     string              `json:"created_at"`
}

// CustomerStats represents customer statistics
type CustomerStats struct {
	TotalCustomers    int64   `json:"total_customers"`
	ActiveCustomers   int64   `json:"active_customers"`
	NewCustomersToday int64   `json:"new_customers_today"`
	NewCustomersMonth int64   `json:"new_customers_month"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// customerRepository is the concrete implementation
type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) ListAdmin(filter models.CustomerListFilter) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	query := r.db.Model(&models.Customer{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", search, search, search)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.Limit
	query = query.Order(filter.SortBy + " " + filter.SortOrder).Offset(offset).Limit(filter.Limit)

	if err := query.Find(&customers).Error; err != nil {
		return nil, 0, err
	}
	return customers, total, nil
}

func (r *customerRepository) GetByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Create(req *models.CreateCustomerRequest, createdBy *uuid.UUID) (*models.Customer, error) {
	customer := &models.Customer{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Status:    "active",
	}
	if err := r.db.Create(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) Update(id uuid.UUID, req *models.UpdateCustomerRequest) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := r.db.Model(&customer).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) GetCustomerOrders(customerID uuid.UUID, page, limit int) ([]CustomerOrderSummary, int64, error) {
	var total int64

	offset := (page - 1) * limit

	// Count total orders
	if err := r.db.Table("public.orders").
		Where("customer_id = ? AND deleted_at IS NULL", customerID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Struct for raw order data
	type rawOrder struct {
		ID            uuid.UUID `gorm:"column:id"`
		OrderNumber   string    `gorm:"column:order_number"`
		Total         float64   `gorm:"column:total"`
		Subtotal      float64   `gorm:"column:subtotal"`
		Status        string    `gorm:"column:status"`
		PaymentStatus string    `gorm:"column:payment_status"`
		CreatedAt     string    `gorm:"column:created_at"`
	}

	var rawOrders []rawOrder

	// Fetch orders
	if err := r.db.Table("public.orders").
		Select("id, order_number, total, subtotal, status, payment_status, created_at").
		Where("customer_id = ? AND deleted_at IS NULL", customerID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&rawOrders).Error; err != nil {
		return nil, 0, err
	}

	// Convert to CustomerOrderSummary and fetch items for each order
	orders := make([]CustomerOrderSummary, len(rawOrders))
	for i, ro := range rawOrders {
		orders[i] = CustomerOrderSummary{
			ID:            ro.ID,
			OrderNum:      ro.OrderNumber,
			Total:         ro.Total,
			Subtotal:      ro.Subtotal,
			Status:        ro.Status,
			PaymentStatus: ro.PaymentStatus,
			CreatedAt:     ro.CreatedAt,
			Items:         []CustomerOrderItem{},
		}

		// Fetch order items
		var items []CustomerOrderItem
		if err := r.db.Table("public.order_items").
			Select("id, product_id, product_name, sku, quantity, unit_price, (quantity * unit_price) as total, image_url").
			Where("order_id = ?", ro.ID).
			Scan(&items).Error; err == nil {
			orders[i].Items = items
		}
	}

	return orders, total, nil
}

func (r *customerRepository) AddNote(customerID uuid.UUID, note string, isPrivate bool, createdBy uuid.UUID) (*models.CustomerNote, error) {
	n := &models.CustomerNote{
		CustomerID: customerID,
		Note:       note,
		IsPrivate:  isPrivate,
		CreatedBy:  &createdBy,
	}
	if err := r.db.Create(n).Error; err != nil {
		return nil, err
	}
	return n, nil
}

func (r *customerRepository) GetNotes(customerID uuid.UUID) ([]models.CustomerNote, error) {
	var notes []models.CustomerNote
	if err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *customerRepository) GetActivity(customerID uuid.UUID, page, limit int) ([]models.CustomerActivity, int64, error) {
	var activities []models.CustomerActivity
	var total int64

	query := r.db.Model(&models.CustomerActivity{}).Where("customer_id = ?", customerID)
	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&activities).Error; err != nil {
		return nil, 0, err
	}
	return activities, total, nil
}

func (r *customerRepository) GetSegments() ([]models.CustomerSegment, error) {
	var segments []models.CustomerSegment
	if err := r.db.Find(&segments).Error; err != nil {
		return nil, err
	}
	return segments, nil
}

func (r *customerRepository) CreateSegment(name, description string, conditions interface{}, color string) (*models.CustomerSegment, error) {
	segment := &models.CustomerSegment{
		Name:        name,
		Description: description,
		Color:       color,
	}
	if err := r.db.Create(segment).Error; err != nil {
		return nil, err
	}
	return segment, nil
}

func (r *customerRepository) UpdateSegment(id uuid.UUID, name, description *string, conditions interface{}, color *string) (*models.CustomerSegment, error) {
	var segment models.CustomerSegment
	if err := r.db.First(&segment, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if name != nil {
		updates["name"] = *name
	}
	if description != nil {
		updates["description"] = *description
	}
	if color != nil {
		updates["color"] = *color
	}

	if err := r.db.Model(&segment).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &segment, nil
}

func (r *customerRepository) DeleteSegment(id uuid.UUID) error {
	return r.db.Delete(&models.CustomerSegment{}, "id = ?", id).Error
}

func (r *customerRepository) AssignSegments(customerID uuid.UUID, segmentIDs []uuid.UUID) error {
	// Clear existing assignments
	r.db.Where("customer_id = ?", customerID).Delete(&models.CustomerSegmentAssignment{})

	// Create new assignments
	for _, segmentID := range segmentIDs {
		assignment := &models.CustomerSegmentAssignment{
			CustomerID: customerID,
			SegmentID:  segmentID,
		}
		if err := r.db.Create(assignment).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *customerRepository) Export(filter models.CustomerListFilter, format string) (interface{}, error) {
	customers, _, err := r.ListAdmin(filter)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *customerRepository) GetStats() (*CustomerStats, error) {
	stats := &CustomerStats{}

	r.db.Model(&models.Customer{}).Count(&stats.TotalCustomers)
	r.db.Model(&models.Customer{}).Where("status = ?", "active").Count(&stats.ActiveCustomers)
	r.db.Model(&models.Customer{}).Where("created_at >= CURRENT_DATE").Count(&stats.NewCustomersToday)
	r.db.Model(&models.Customer{}).Where("created_at >= date_trunc('month', CURRENT_DATE)").Count(&stats.NewCustomersMonth)

	return stats, nil
}
