package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"placeintel-pro/api/models"
)

const (
	FoursquareAPIBase = "https://api.foursquare.com/v3"
	DefaultRadius     = 1000
	DefaultLimit      = 20
)

// FoursquareService handles interactions with Foursquare Places API
type FoursquareService struct {
	apiKey     string
	httpClient *http.Client
}

// NewFoursquareService creates a new Foursquare service instance
func NewFoursquareService(apiKey string) *FoursquareService {
	return &FoursquareService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FoursquareSearchResponse represents the response from Foursquare search API
type FoursquareSearchResponse struct {
	Results []models.FoursquarePlace `json:"results"`
	Context struct {
		GeoBounds struct {
			Circle struct {
				Center struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"center"`
				Radius int `json:"radius"`
			} `json:"circle"`
		} `json:"geo_bounds"`
	} `json:"context"`
}

// FoursquarePlaceDetails represents detailed place information from Foursquare
type FoursquarePlaceDetails struct {
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
	Tel         string `json:"tel"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Hours       struct {
		Display   string `json:"display"`
		IsLocalHoliday bool `json:"is_local_holiday"`
		OpenNow   bool   `json:"open_now"`
		Regular   []struct {
			Close string `json:"close"`
			Day   int    `json:"day"`
			Open  string `json:"open"`
		} `json:"regular"`
	} `json:"hours"`
	Rating  float64 `json:"rating"`
	Stats   struct {
		TotalPhotos   int `json:"total_photos"`
		TotalRatings  int `json:"total_ratings"`
		TotalTips     int `json:"total_tips"`
	} `json:"stats"`
	Price int `json:"price"`
}

// SearchPlaces searches for places using Foursquare Places API
func (fs *FoursquareService) SearchPlaces(req models.PlaceSearchRequest) ([]models.FoursquarePlace, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("ll", fmt.Sprintf("%.6f,%.6f", req.Latitude, req.Longitude))
	
	if req.Query != "" {
		params.Add("query", req.Query)
	}
	
	if req.Categories != "" {
		params.Add("categories", req.Categories)
	}
	
	radius := req.Radius
	if radius == 0 {
		radius = DefaultRadius
	}
	params.Add("radius", strconv.Itoa(radius))
	
	limit := req.Limit
	if limit == 0 {
		limit = DefaultLimit
	}
	params.Add("limit", strconv.Itoa(limit))
	
	// Add additional useful fields
	params.Add("fields", "fsq_id,name,location,categories,distance,tel,website,rating,price,hours")

	// Make API request
	apiURL := fmt.Sprintf("%s/places/search?%s", FoursquareAPIBase, params.Encode())
	
	logrus.WithFields(logrus.Fields{
		"url":    apiURL,
		"query":  req.Query,
		"lat":    req.Latitude,
		"lng":    req.Longitude,
		"radius": radius,
	}).Info("Searching places via Foursquare API")

	resp, err := fs.makeRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search places: %w", err)
	}

	var searchResp FoursquareSearchResponse
	if err := json.Unmarshal(resp, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	logrus.WithField("results_count", len(searchResp.Results)).Info("Places search completed")
	return searchResp.Results, nil
}

// GetPlaceDetails retrieves detailed information for a specific place
func (fs *FoursquareService) GetPlaceDetails(placeID string) (*FoursquarePlaceDetails, error) {
	// Build API URL with comprehensive fields
	fields := "fsq_id,name,location,categories,tel,website,email,description,hours,rating,stats,price,photos"
	apiURL := fmt.Sprintf("%s/places/%s?fields=%s", FoursquareAPIBase, placeID, fields)

	logrus.WithFields(logrus.Fields{
		"place_id": placeID,
		"url":      apiURL,
	}).Info("Fetching place details via Foursquare API")

	resp, err := fs.makeRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get place details: %w", err)
	}

	var placeDetails FoursquarePlaceDetails
	if err := json.Unmarshal(resp, &placeDetails); err != nil {
		return nil, fmt.Errorf("failed to parse place details response: %w", err)
	}

	logrus.WithField("place_name", placeDetails.Name).Info("Place details retrieved successfully")
	return &placeDetails, nil
}

// makeRequest makes an HTTP request to Foursquare API with proper authentication
func (fs *FoursquareService) makeRequest(method, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	req.Header.Set("Authorization", fs.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := fs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("Foursquare API request failed")
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// HealthCheck verifies connectivity to Foursquare API
func (fs *FoursquareService) HealthCheck() error {
	// Make a simple search request to verify API connectivity
	testURL := fmt.Sprintf("%s/places/search?ll=40.7128,-74.0060&limit=1", FoursquareAPIBase)
	
	_, err := fs.makeRequest("GET", testURL)
	if err != nil {
		return fmt.Errorf("foursquare API health check failed: %w", err)
	}
	
	return nil
}