package models

import "time"

// PlaceSearchRequest represents a search request for places
type PlaceSearchRequest struct {
	Query      string  `json:"query" form:"query"`
	Latitude   float64 `json:"lat" form:"lat" binding:"required"`
	Longitude  float64 `json:"lng" form:"lng" binding:"required"`
	Radius     int     `json:"radius" form:"radius"`
	Categories string  `json:"categories" form:"categories"`
	Limit      int     `json:"limit" form:"limit"`
}

// FoursquarePlace represents basic place data from Foursquare API
type FoursquarePlace struct {
	FSQId    string `json:"fsq_id"`
	Name     string `json:"name"`
	Location struct {
		Address     string  `json:"address"`
		Country     string  `json:"country"`
		CrossStreet string  `json:"cross_street"`
		Locality    string  `json:"locality"`
		Postcode    string  `json:"postcode"`
		Region      string  `json:"region"`
		Latitude    float64 `json:"lat"`
		Longitude   float64 `json:"lng"`
	} `json:"location"`
	Categories []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Icon struct {
			Prefix string `json:"prefix"`
			Suffix string `json:"suffix"`
		} `json:"icon"`
	} `json:"categories"`
	Distance int `json:"distance"`
}

// BusinessIntelligence represents AI-generated business insights
type BusinessIntelligence struct {
	PopularityScore float64  `json:"popularity_score"`
	SentimentScore  float64  `json:"sentiment_score"`
	Specialties     []string `json:"specialties"`
	IdealFor        []string `json:"ideal_for"`
	PriceRange      string   `json:"price_range"`
	Atmosphere      string   `json:"atmosphere"`
	TrendingScore   float64  `json:"trending_score"`
}

// RealTimeContext represents live contextual information
type RealTimeContext struct {
	CurrentStatus        string    `json:"current_status"`
	CrowdLevel          string    `json:"crowd_level"`
	BestVisitTimes      []string  `json:"best_visit_times"`
	LiveEvents          []string  `json:"live_events"`
	EstimatedWaitTime   string    `json:"estimated_wait_time"`
	WeatherImpact       string    `json:"weather_impact"`
	LastUpdated         time.Time `json:"last_updated"`
	ConfidenceScore     float64   `json:"confidence_score"`
}

// AccessibilityIntelligence represents accessibility information
type AccessibilityIntelligence struct {
	WheelchairAccessible bool    `json:"wheelchair_accessible"`
	AccessibilityScore   float64 `json:"accessibility_score"`
	Features             struct {
		RampAccess           bool `json:"ramp_access"`
		Elevator             bool `json:"elevator"`
		AccessibleRestrooms  bool `json:"accessible_restrooms"`
		BrailleSignage       bool `json:"braille_signage"`
		HearingLoop          bool `json:"hearing_loop"`
		WideEntrances        bool `json:"wide_entrances"`
		AccessibleParking    bool `json:"accessible_parking"`
	} `json:"features"`
	InclusiveRecommendations struct {
		MobilityFriendlyAreas   []string `json:"mobility_friendly_areas"`
		SensoryAccommodations   []string `json:"sensory_accommodations"`
		CognitiveSupport        []string `json:"cognitive_support"`
	} `json:"inclusive_recommendations"`
}

// UnifiedRecommendations represents AI-generated unified insights
type UnifiedRecommendations struct {
	ConfidenceScore       float64  `json:"confidence_score"`
	PersonalizedInsights  []string `json:"personalized_insights"`
	AlternativeSuggestions []string `json:"alternative_suggestions"`
	OptimalVisitStrategy   string   `json:"optimal_visit_strategy"`
	AccessibilityNotes     []string `json:"accessibility_notes"`
}

// PlaceIntelligence represents the complete enhanced place data
type PlaceIntelligence struct {
	// Basic place information
	FSQId        string          `json:"fsq_id"`
	Name         string          `json:"name"`
	Location     interface{}     `json:"location"`
	Categories   interface{}     `json:"categories"`
	Distance     int             `json:"distance"`
	
	// Enhanced intelligence
	BusinessIntelligence      BusinessIntelligence      `json:"business_intelligence"`
	RealTimeContext          RealTimeContext           `json:"real_time_context"`
	AccessibilityIntelligence AccessibilityIntelligence `json:"accessibility_intelligence"`
	UnifiedRecommendations   UnifiedRecommendations    `json:"unified_recommendations"`
	
	// Metadata
	ProcessingTime time.Duration `json:"processing_time_ms"`
	DataSources    []string      `json:"data_sources"`
	LastUpdated    time.Time     `json:"last_updated"`
}

// PlaceSearchResponse represents the API response for place search
type PlaceSearchResponse struct {
	Results []PlaceIntelligence `json:"results"`
	Meta    struct {
		Total          int           `json:"total"`
		ProcessingTime time.Duration `json:"processing_time_ms"`
		DataSources    []string      `json:"data_sources"`
	} `json:"meta"`
}

// ErrorResponse represents API error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Services  struct {
		Foursquare   string `json:"foursquare"`
		Intelligence string `json:"intelligence"`
		Cache        string `json:"cache"`
	} `json:"services"`
}