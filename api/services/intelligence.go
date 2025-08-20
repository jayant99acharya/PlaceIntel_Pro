package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"placeintel-pro/api/models"
)

// IntelligenceService handles communication with Python ML/AI engine
type IntelligenceService struct {
	baseURL    string
	httpClient *http.Client
}

// NewIntelligenceService creates a new intelligence service instance
func NewIntelligenceService(baseURL string) *IntelligenceService {
	return &IntelligenceService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// IntelligenceRequest represents the request to Python intelligence service
type IntelligenceRequest struct {
	Place    models.FoursquarePlace `json:"place"`
	Context  map[string]interface{} `json:"context"`
	Features []string               `json:"features"`
}

// IntelligenceResponse represents the response from Python intelligence service
type IntelligenceResponse struct {
	BusinessIntelligence       models.BusinessIntelligence       `json:"business_intelligence"`
	RealTimeContext            models.RealTimeContext            `json:"real_time_context"`
	AccessibilityIntelligence  models.AccessibilityIntelligence  `json:"accessibility_intelligence"`
	UnifiedRecommendations     models.UnifiedRecommendations     `json:"unified_recommendations"`
	ProcessingTime             float64                           `json:"processing_time_ms"`
	DataSources                []string                          `json:"data_sources"`
}

// EnhancePlaceWithIntelligence enriches place data with AI-generated intelligence
func (is *IntelligenceService) EnhancePlaceWithIntelligence(place models.FoursquarePlace) (*models.PlaceIntelligence, error) {
	startTime := time.Now()

	// Prepare request payload
	req := IntelligenceRequest{
		Place: place,
		Context: map[string]interface{}{
			"timestamp": time.Now().UTC(),
			"location": map[string]interface{}{
				"lat": place.Location.Latitude,
				"lng": place.Location.Longitude,
			},
		},
		Features: []string{
			"business_intelligence",
			"real_time_context",
			"accessibility_intelligence",
			"unified_recommendations",
		},
	}

	// Make request to Python intelligence service
	intelligenceResp, err := is.processIntelligence(req)
	if err != nil {
		// If intelligence service fails, return basic place data with empty intelligence
		logrus.WithError(err).Warn("Intelligence service failed, returning basic place data")
		return is.createBasicPlaceIntelligence(place, time.Since(startTime)), nil
	}

	// Create enhanced place intelligence
	placeIntel := &models.PlaceIntelligence{
		FSQId:      place.FSQPlaceId,
		Name:       place.Name,
		Location:   place.Location,
		Categories: place.Categories,
		Distance:   place.Distance,

		BusinessIntelligence:      intelligenceResp.BusinessIntelligence,
		RealTimeContext:           intelligenceResp.RealTimeContext,
		AccessibilityIntelligence: intelligenceResp.AccessibilityIntelligence,
		UnifiedRecommendations:    intelligenceResp.UnifiedRecommendations,

		ProcessingTime: time.Since(startTime),
		DataSources:    intelligenceResp.DataSources,
		LastUpdated:    time.Now().UTC(),
	}

	logrus.WithFields(logrus.Fields{
		"place_id":        place.FSQPlaceId,
		"place_name":      place.Name,
		"processing_time": time.Since(startTime),
	}).Info("Place intelligence enhancement completed")

	return placeIntel, nil
}

// EnhancePlacesWithIntelligence processes multiple places concurrently
func (is *IntelligenceService) EnhancePlacesWithIntelligence(places []models.FoursquarePlace) ([]models.PlaceIntelligence, error) {
	if len(places) == 0 {
		return []models.PlaceIntelligence{}, nil
	}

	// Process places concurrently for better performance
	results := make([]models.PlaceIntelligence, len(places))
	errors := make([]error, len(places))
	
	// Use a semaphore to limit concurrent requests
	semaphore := make(chan struct{}, 5) // Max 5 concurrent requests
	done := make(chan int, len(places))

	for i, place := range places {
		go func(index int, p models.FoursquarePlace) {
			semaphore <- struct{}{} // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			enhanced, err := is.EnhancePlaceWithIntelligence(p)
			if err != nil {
				errors[index] = err
				// Create basic place intelligence on error
				enhanced = is.createBasicPlaceIntelligence(p, 0)
			}
			results[index] = *enhanced
			done <- index
		}(i, place)
	}

	// Wait for all goroutines to complete
	for i := 0; i < len(places); i++ {
		<-done
	}

	// Log any errors but don't fail the entire request
	for i, err := range errors {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"place_index": i,
				"place_name":  places[i].Name,
				"error":       err,
			}).Warn("Failed to enhance place with intelligence")
		}
	}

	return results, nil
}

// processIntelligence makes a request to the Python intelligence service
func (is *IntelligenceService) processIntelligence(req IntelligenceRequest) (*IntelligenceResponse, error) {
	// Serialize request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal intelligence request: %w", err)
	}

	// Make HTTP request
	url := fmt.Sprintf("%s/api/v1/intelligence/enhance", is.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create intelligence request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := is.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call intelligence service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read intelligence response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("intelligence service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var intelligenceResp IntelligenceResponse
	if err := json.Unmarshal(body, &intelligenceResp); err != nil {
		return nil, fmt.Errorf("failed to parse intelligence response: %w", err)
	}

	return &intelligenceResp, nil
}

// createBasicPlaceIntelligence creates a basic place intelligence when AI service fails
func (is *IntelligenceService) createBasicPlaceIntelligence(place models.FoursquarePlace, processingTime time.Duration) *models.PlaceIntelligence {
	// Create basic intelligence with default values
	businessIntel := models.BusinessIntelligence{
		PopularityScore: 5.0, // Default neutral score
		SentimentScore:  3.5, // Default neutral sentiment
		Specialties:     []string{},
		IdealFor:        []string{},
		PriceRange:      "unknown",
		Atmosphere:      "unknown",
		TrendingScore:   0.0,
	}

	realTimeContext := models.RealTimeContext{
		CurrentStatus:     "unknown",
		CrowdLevel:        "unknown",
		BestVisitTimes:    []string{},
		LiveEvents:        []string{},
		EstimatedWaitTime: "unknown",
		WeatherImpact:     "none",
		LastUpdated:       time.Now().UTC(),
		ConfidenceScore:   0.0,
	}

	accessibilityIntel := models.AccessibilityIntelligence{
		WheelchairAccessible: false, // Conservative default
		AccessibilityScore:   0.0,   // Unknown
		Features: struct {
			RampAccess           bool `json:"ramp_access"`
			Elevator             bool `json:"elevator"`
			AccessibleRestrooms  bool `json:"accessible_restrooms"`
			BrailleSignage       bool `json:"braille_signage"`
			HearingLoop          bool `json:"hearing_loop"`
			WideEntrances        bool `json:"wide_entrances"`
			AccessibleParking    bool `json:"accessible_parking"`
		}{},
		InclusiveRecommendations: struct {
			MobilityFriendlyAreas []string `json:"mobility_friendly_areas"`
			SensoryAccommodations []string `json:"sensory_accommodations"`
			CognitiveSupport      []string `json:"cognitive_support"`
		}{
			MobilityFriendlyAreas: []string{},
			SensoryAccommodations: []string{},
			CognitiveSupport:      []string{},
		},
	}

	unifiedRecommendations := models.UnifiedRecommendations{
		ConfidenceScore:        0.0,
		PersonalizedInsights:   []string{"Basic place information available"},
		AlternativeSuggestions: []string{},
		OptimalVisitStrategy:   "Contact venue for current information",
		AccessibilityNotes:     []string{"Accessibility information not available - please contact venue"},
	}

	return &models.PlaceIntelligence{
		FSQId:      place.FSQPlaceId,
		Name:       place.Name,
		Location:   place.Location,
		Categories: place.Categories,
		Distance:   place.Distance,

		BusinessIntelligence:      businessIntel,
		RealTimeContext:          realTimeContext,
		AccessibilityIntelligence: accessibilityIntel,
		UnifiedRecommendations:   unifiedRecommendations,

		ProcessingTime: processingTime,
		DataSources:    []string{"foursquare"},
		LastUpdated:    time.Now().UTC(),
	}
}

// HealthCheck verifies connectivity to intelligence service
func (is *IntelligenceService) HealthCheck() error {
	url := fmt.Sprintf("%s/health", is.baseURL)
	
	resp, err := is.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("intelligence service health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("intelligence service returned status %d", resp.StatusCode)
	}

	return nil
}