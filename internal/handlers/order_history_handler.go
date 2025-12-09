package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niaga-platform/service-customer/internal/middleware"
)

// OrderHistoryHandler handles order history requests
type OrderHistoryHandler struct {
	orderServiceURL string
	httpClient      *http.Client
}

// NewOrderHistoryHandler creates a new order history handler
func NewOrderHistoryHandler() *OrderHistoryHandler {
	orderURL := os.Getenv("ORDER_SERVICE_URL")
	if orderURL == "" {
		orderURL = "http://kilang-order:8005"
	}

	return &OrderHistoryHandler{
		orderServiceURL: orderURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// OrderResponse represents the response from service-order
type OrderResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Orders []Order `json:"orders"`
		Total  int64   `json:"total"`
		Page   int     `json:"page"`
		Limit  int     `json:"limit"`
	} `json:"data"`
}

// Order represents an order from service-order
type Order struct {
	ID              string  `json:"id"`
	OrderNumber     string  `json:"orderNumber"`
	Status          string  `json:"status"`
	PaymentStatus   string  `json:"paymentStatus"`
	Total           float64 `json:"total"`
	ShippingAddress struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		City    string `json:"city"`
	} `json:"shippingAddress"`
	CreatedAt string `json:"createdAt"`
}

// GetOrderHistory retrieves the customer's order history
// GET /api/v1/customer/orders
func (h *OrderHistoryHandler) GetOrderHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Get pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Build request URL
	url := fmt.Sprintf("%s/api/v1/orders?page=%d&limit=%d", h.orderServiceURL, page, limit)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Forward authorization header
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("X-User-ID", userID.String())

	// Make request to service-order
	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Order service unavailable",
			"orders":  []gin.H{},
			"total":   0,
			"user_id": userID.String(),
		})
		return
	}
	defer resp.Body.Close()

	// Parse response
	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse order response"})
		return
	}

	if !orderResp.Success {
		c.JSON(http.StatusOK, gin.H{
			"orders":  []gin.H{},
			"total":   0,
			"page":    page,
			"limit":   limit,
			"user_id": userID.String(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"orders":  orderResp.Data.Orders,
		"total":   orderResp.Data.Total,
		"page":    orderResp.Data.Page,
		"limit":   orderResp.Data.Limit,
		"user_id": userID.String(),
	})
}

// GetOrder retrieves a single order by ID
// GET /api/v1/customer/orders/:id
func (h *OrderHistoryHandler) GetOrder(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID required"})
		return
	}

	// Build request URL
	url := fmt.Sprintf("%s/api/v1/orders/%s", h.orderServiceURL, orderID)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Forward authorization header
	req.Header.Set("Authorization", c.GetHeader("Authorization"))
	req.Header.Set("X-User-ID", userID.String())

	// Make request to service-order
	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Order service unavailable"})
		return
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse order response"})
		return
	}

	c.JSON(resp.StatusCode, result)
}
