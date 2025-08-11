package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"placeintel-pro/api/handlers"
	"placeintel-pro/api/middleware"
	"placeintel-pro/api/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	// Initialize logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if level, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info")); err == nil {
		logrus.SetLevel(level)
	}

	// Initialize services
	foursquareService := services.NewFoursquareService(getEnv("FOURSQUARE_API_KEY", ""))
	intelligenceService := services.NewIntelligenceService(getEnv("PYTHON_SERVICE_URL", "http://localhost:5000"))
	cacheService := services.NewCacheService(
		getEnv("REDIS_HOST", "localhost"),
		getEnv("REDIS_PORT", "6379"),
		getEnv("REDIS_PASSWORD", ""),
	)

	// Initialize Gin router
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.New()
	
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit())

	// Initialize handlers
	placeHandler := handlers.NewPlaceHandler(foursquareService, intelligenceService, cacheService)

	// API Routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
				"version":   "1.0.0",
			})
		})

		// Place intelligence endpoints
		places := v1.Group("/places")
		{
			places.GET("/search", placeHandler.SearchPlaces)
			places.GET("/intelligence", placeHandler.GetPlaceIntelligence)
			places.GET("/:place_id/details", placeHandler.GetPlaceDetails)
			places.GET("/:place_id/intelligence", placeHandler.GetPlaceIntelligenceByID)
		}

		// Analytics endpoints
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/popular", placeHandler.GetPopularPlaces)
			analytics.GET("/trends", placeHandler.GetTrends)
		}
	}

	// API Documentation
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "PlaceIntel Pro API",
			"version":     "1.0.0",
			"description": "Universal Location Intelligence Platform",
			"endpoints": gin.H{
				"health":                    "GET /api/v1/health",
				"search_places":             "GET /api/v1/places/search",
				"place_intelligence":        "GET /api/v1/places/intelligence",
				"place_details":             "GET /api/v1/places/:place_id/details",
				"place_intelligence_by_id":  "GET /api/v1/places/:place_id/intelligence",
				"popular_places":            "GET /api/v1/analytics/popular",
				"trends":                    "GET /api/v1/analytics/trends",
			},
		})
	})

	// Start server
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		logrus.Infof("PlaceIntel Pro API server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}

	logrus.Info("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}