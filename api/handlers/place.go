package handlers

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"placeintel-pro/api/models"
	"placeintel-pro/api/services"
)

// PlaceHandler handles place-related API endpoints
type PlaceHandler struct {
	foursquareService   *services.FoursquareService
	intelligenceService *services.IntelligenceService
	cacheService        *services.CacheService
}

// NewPlaceHandler creates a new place handler instance
func NewPlaceHandler(
	foursquareService *services.FoursquareService,
	intelligenceService *services.IntelligenceService,
	cacheService *services.CacheService,
) *PlaceHandler {
	return &PlaceHandler{
		foursquareService:   foursquareService,
		intelligenceService: intelligenceService,
		cacheService:        cacheService,
	}
}

// SearchPlaces handles place search requests
func (ph *PlaceHandler) SearchPlaces(c *gin.Context) {
	startTime := time.Now()

	// Parse and validate request parameters
	var req models.PlaceSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate required parameters
	if req.Latitude == 0 || req.Longitude == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing required parameters",
			Message: "Latitude and longitude are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Generate cache key
	cacheKey := ph.generateSearchCacheKey(req)

	// Try to get results from cache first
	if cachedResults, err := ph.cacheService.GetCachedSearchResults(cacheKey); err == nil && cachedResults != nil {
		logrus.WithField("cache_key", cacheKey).Info("Returning cached search results")
		
		response := models.PlaceSearchResponse{
			Results: cachedResults,
			Meta: struct {
				Total          int           `json:"total"`
				ProcessingTime time.Duration `json:"processing_time_ms"`
				DataSources    []string      `json:"data_sources"`
			}{
				Total:          len(cachedResults),
				ProcessingTime: time.Since(startTime),
				DataSources:    []string{"cache", "foursquare", "intelligence"},
			},
		}
		
		c.JSON(http.StatusOK, response)
		return
	}

	// Search places using Foursquare API
	places, err := ph.foursquareService.SearchPlaces(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to search places")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to search places",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Enhance places with intelligence
	enhancedPlaces, err := ph.intelligenceService.EnhancePlacesWithIntelligence(places)
	if err != nil {
		logrus.WithError(err).Error("Failed to enhance places with intelligence")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to enhance places with intelligence",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Cache the results
	if err := ph.cacheService.CacheSearchResults(cacheKey, enhancedPlaces); err != nil {
		logrus.WithError(err).Warn("Failed to cache search results")
	}

	// Prepare response
	response := models.PlaceSearchResponse{
		Results: enhancedPlaces,
		Meta: struct {
			Total          int           `json:"total"`
			ProcessingTime time.Duration `json:"processing_time_ms"`
			DataSources    []string      `json:"data_sources"`
		}{
			Total:          len(enhancedPlaces),
			ProcessingTime: time.Since(startTime),
			DataSources:    []string{"foursquare", "intelligence"},
		},
	}

	logrus.WithFields(logrus.Fields{
		"results":         len(enhancedPlaces),
		"processing_time": time.Since(startTime),
		"query":          req.Query,
		"location":       fmt.Sprintf("%.6f,%.6f", req.Latitude, req.Longitude),
	}).Info("Place search completed successfully")

	c.JSON(http.StatusOK, response)
}

// GetPlaceIntelligence handles requests for place intelligence by search parameters
func (ph *PlaceHandler) GetPlaceIntelligence(c *gin.Context) {
	// This endpoint is similar to SearchPlaces but focuses on intelligence
	ph.SearchPlaces(c)
}

// GetPlaceDetails handles requests for detailed place information
func (ph *PlaceHandler) GetPlaceDetails(c *gin.Context) {
	placeID := c.Param("place_id")
	if placeID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing place ID",
			Message: "Place ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	startTime := time.Now()

	// Try cache first
	if cachedDetails, err := ph.cacheService.GetCachedPlaceDetails(placeID); err == nil && cachedDetails != nil {
		logrus.WithField("place_id", placeID).Info("Returning cached place details")
		c.JSON(http.StatusOK, cachedDetails)
		return
	}

	// Get place details from Foursquare
	details, err := ph.foursquareService.GetPlaceDetails(placeID)
	if err != nil {
		logrus.WithError(err).WithField("place_id", placeID).Error("Failed to get place details")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get place details",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Cache the details
	if err := ph.cacheService.CachePlaceDetails(placeID, details); err != nil {
		logrus.WithError(err).Warn("Failed to cache place details")
	}

	logrus.WithFields(logrus.Fields{
		"place_id":        placeID,
		"place_name":      details.Name,
		"processing_time": time.Since(startTime),
	}).Info("Place details retrieved successfully")

	c.JSON(http.StatusOK, details)
}

// GetPlaceIntelligenceByID handles requests for place intelligence by place ID
func (ph *PlaceHandler) GetPlaceIntelligenceByID(c *gin.Context) {
	placeID := c.Param("place_id")
	if placeID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing place ID",
			Message: "Place ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	startTime := time.Now()

	// Try cache first
	if cachedIntelligence, err := ph.cacheService.GetCachedPlaceIntelligence(placeID); err == nil && cachedIntelligence != nil {
		logrus.WithField("place_id", placeID).Info("Returning cached place intelligence")
		c.JSON(http.StatusOK, cachedIntelligence)
		return
	}

	// Get place details from Foursquare first
	details, err := ph.foursquareService.GetPlaceDetails(placeID)
	if err != nil {
		logrus.WithError(err).WithField("place_id", placeID).Error("Failed to get place details for intelligence")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get place details",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert details to FoursquarePlace format for intelligence processing
	place := models.FoursquarePlace{
		FSQId: details.FSQId,
		Name:  details.Name,
		Location: struct {
			Address     string  `json:"address"`
			Country     string  `json:"country"`
			CrossStreet string  `json:"cross_street"`
			Locality    string  `json:"locality"`
			Postcode    string  `json:"postcode"`
			Region      string  `json:"region"`
			Latitude    float64 `json:"lat"`
			Longitude   float64 `json:"lng"`
		}{
			Address:     details.Location.Address,
			Country:     details.Location.Country,
			CrossStreet: details.Location.CrossStreet,
			Locality:    details.Location.Locality,
			Postcode:    details.Location.Postcode,
			Region:      details.Location.Region,
			Latitude:    details.Location.Latitude,
			Longitude:   details.Location.Longitude,
		},
		Categories: details.Categories,
		Distance:   0, // Not applicable for direct place lookup
	}

	// Enhance with intelligence
	intelligence, err := ph.intelligenceService.EnhancePlaceWithIntelligence(place)
	if err != nil {
		logrus.WithError(err).WithField("place_id", placeID).Error("Failed to enhance place with intelligence")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to enhance place with intelligence",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Cache the intelligence
	if err := ph.cacheService.CachePlaceIntelligence(placeID, intelligence); err != nil {
		logrus.WithError(err).Warn("Failed to cache place intelligence")
	}

	logrus.WithFields(logrus.Fields{
		"place_id":        placeID,
		"place_name":      intelligence.Name,
		"processing_time": time.Since(startTime),
	}).Info("Place intelligence retrieved successfully")

	c.JSON(http.StatusOK, intelligence)
}

// GetPopularPlaces handles requests for popular places in an area
func (ph *PlaceHandler) GetPopularPlaces(c *gin.Context) {
	// Parse location parameters
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	
	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing location parameters",
			Message: "Latitude and longitude are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid latitude",
			Message: "Latitude must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid longitude",
			Message: "Longitude must be a valid number",
			Code:    http.StatusBadRequest,
		})
		return
	}

	locationKey := fmt.Sprintf("%.4f,%.4f", lat, lng)

	// Try cache first
	if cachedPopular, err := ph.cacheService.GetCachedPopularPlaces(locationKey); err == nil && cachedPopular != nil {
		logrus.WithField("location", locationKey).Info("Returning cached popular places")
		c.JSON(http.StatusOK, gin.H{
			"popular_places": cachedPopular,
			"location":       locationKey,
			"cached":         true,
		})
		return
	}

	// Search for popular places (high-rated, trending)
	req := models.PlaceSearchRequest{
		Latitude:  lat,
		Longitude: lng,
		Radius:    2000, // 2km radius for popular places
		Limit:     20,   // Top 20 popular places
	}

	places, err := ph.foursquareService.SearchPlaces(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to search for popular places")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to search for popular places",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Enhance with intelligence and filter for popular ones
	enhancedPlaces, err := ph.intelligenceService.EnhancePlacesWithIntelligence(places)
	if err != nil {
		logrus.WithError(err).Error("Failed to enhance popular places with intelligence")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to enhance places with intelligence",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Filter and sort by popularity score
	var popularPlaces []models.PlaceIntelligence
	for _, place := range enhancedPlaces {
		if place.BusinessIntelligence.PopularityScore >= 7.0 { // High popularity threshold
			popularPlaces = append(popularPlaces, place)
		}
	}

	// Cache the results
	if err := ph.cacheService.CachePopularPlaces(locationKey, popularPlaces); err != nil {
		logrus.WithError(err).Warn("Failed to cache popular places")
	}

	c.JSON(http.StatusOK, gin.H{
		"popular_places": popularPlaces,
		"location":       locationKey,
		"total":          len(popularPlaces),
		"cached":         false,
	})
}

// GetTrends handles requests for trending places and insights
func (ph *PlaceHandler) GetTrends(c *gin.Context) {
	// This is a placeholder for trends analysis
	// In a real implementation, this would analyze historical data
	c.JSON(http.StatusOK, gin.H{
		"trends": gin.H{
			"trending_categories": []string{"coffee", "restaurants", "fitness"},
			"peak_hours":         []string{"12:00-14:00", "18:00-20:00"},
			"popular_areas":      []string{"downtown", "business district"},
		},
		"message": "Trends analysis coming soon",
	})
}

// generateSearchCacheKey creates a unique cache key for search requests
func (ph *PlaceHandler) generateSearchCacheKey(req models.PlaceSearchRequest) string {
	key := fmt.Sprintf("%.6f,%.6f,%s,%s,%d,%d",
		req.Latitude, req.Longitude, req.Query, req.Categories, req.Radius, req.Limit)
	
	// Create MD5 hash for shorter, consistent keys
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}