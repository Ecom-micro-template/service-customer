package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/niaga-platform/service-customer/internal/config"
	"github.com/niaga-platform/service-customer/internal/handlers"
	"github.com/niaga-platform/service-customer/internal/middleware"
	"github.com/niaga-platform/service-customer/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	cfg *config.Config
)

func main() {
	// Load environment
	if os.Getenv("APP_ENV") != "production" {
		godotenv.Load()
	}

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	cfg = config.Load()
	log.Println("✅ Configuration loaded")

	// Initialize database
	var err error
	db, err = gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected")

	// Create schema if it doesn't exist
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS customer").Error; err != nil {
		log.Fatalf("Failed to create customer schema: %v", err)
	}
	log.Println("✅ Customer schema ready")

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.Profile{},
		&models.Address{},
		&models.WishlistItem{},
		&models.CustomerMeasurement{}, // Day 96
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("✅ Database migrations completed")

	// Add unique constraint for wishlist
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_wishlist_user_product 
		ON customer.wishlist_items(user_id, product_id)
	`).Error; err != nil {
		log.Printf("⚠️  Warning: Failed to create unique index on wishlist: %v", err)
	}

	// Initialize handlers
	profileHandler := handlers.NewProfileHandler(db)
	addressHandler := handlers.NewAddressHandler(db)
	wishlistHandler := handlers.NewWishlistHandler(db)
	orderHistoryHandler := handlers.NewOrderHistoryHandler()
	measurementHandler := handlers.NewMeasurementHandler(db) // Day 96

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "customer",
			"time":    time.Now().UTC(),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Customer routes (protected)
		customer := v1.Group("/customer")
		customer.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// Profile
			customer.GET("/profile", profileHandler.GetProfile)
			customer.PUT("/profile", profileHandler.UpdateProfile)

			// Addresses
			customer.GET("/addresses", addressHandler.ListAddresses)
			customer.POST("/addresses", addressHandler.CreateAddress)
			customer.PUT("/addresses/:id", addressHandler.UpdateAddress)
			customer.DELETE("/addresses/:id", addressHandler.DeleteAddress)
			customer.PUT("/addresses/:id/default", addressHandler.SetDefaultAddress)

			// Wishlist
			customer.GET("/wishlist", wishlistHandler.GetWishlist)
			customer.POST("/wishlist", wishlistHandler.AddToWishlist)
			customer.DELETE("/wishlist/:productId", wishlistHandler.RemoveFromWishlist)

			// Order History
			customer.GET("/orders", orderHistoryHandler.GetOrderHistory)

			// Measurements (Day 96)
			customer.GET("/measurements", measurementHandler.List)
			customer.POST("/measurements", measurementHandler.Create)
			customer.GET("/measurements/:id", measurementHandler.GetByID)
			customer.PUT("/measurements/:id", measurementHandler.Update)
			customer.DELETE("/measurements/:id", measurementHandler.Delete)
			customer.PUT("/measurements/:id/set-default", measurementHandler.SetDefault)
		}
	}

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8004"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Customer service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
