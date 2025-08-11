package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"placeintel-pro/api/models"
)

// CacheService handles Redis caching operations
type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheService creates a new cache service instance
func NewCacheService(host, port, password string) *CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0, // Default DB
	})

	return &CacheService{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Cache key prefixes
const (
	PlaceSearchPrefix     = "search:"
	PlaceIntelligencePrefix = "intel:"
	PlaceDetailsPrefix    = "details:"
	PopularPlacesPrefix   = "popular:"
	TrendsPrefix         = "trends:"
)

// Cache TTL durations
const (
	SearchCacheTTL      = 5 * time.Minute   // Search results change frequently
	IntelligenceCacheTTL = 15 * time.Minute // Intelligence data is more stable
	DetailsCacheTTL     = 30 * time.Minute  // Place details change less frequently
	PopularCacheTTL     = 1 * time.Hour     // Popular places aggregated data
	TrendsCacheTTL      = 2 * time.Hour     // Trends data
)

// CacheSearchResults caches place search results
func (cs *CacheService) CacheSearchResults(key string, results []models.PlaceIntelligence) error {
	cacheKey := PlaceSearchPrefix + key
	
	data, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal search results: %w", err)
	}

	err = cs.client.Set(cs.ctx, cacheKey, data, SearchCacheTTL).Err()
	if err != nil {
		logrus.WithError(err).WithField("key", cacheKey).Error("Failed to cache search results")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"key":     cacheKey,
		"ttl":     SearchCacheTTL,
		"results": len(results),
	}).Debug("Cached search results")

	return nil
}

// GetCachedSearchResults retrieves cached search results
func (cs *CacheService) GetCachedSearchResults(key string) ([]models.PlaceIntelligence, error) {
	cacheKey := PlaceSearchPrefix + key
	
	data, err := cs.client.Get(cs.ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached search results: %w", err)
	}

	var results []models.PlaceIntelligence
	if err := json.Unmarshal([]byte(data), &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached search results: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"key":     cacheKey,
		"results": len(results),
	}).Debug("Retrieved cached search results")

	return results, nil
}

// CachePlaceIntelligence caches individual place intelligence
func (cs *CacheService) CachePlaceIntelligence(placeID string, intelligence *models.PlaceIntelligence) error {
	cacheKey := PlaceIntelligencePrefix + placeID
	
	data, err := json.Marshal(intelligence)
	if err != nil {
		return fmt.Errorf("failed to marshal place intelligence: %w", err)
	}

	err = cs.client.Set(cs.ctx, cacheKey, data, IntelligenceCacheTTL).Err()
	if err != nil {
		logrus.WithError(err).WithField("key", cacheKey).Error("Failed to cache place intelligence")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"key":      cacheKey,
		"place_id": placeID,
		"ttl":      IntelligenceCacheTTL,
	}).Debug("Cached place intelligence")

	return nil
}

// GetCachedPlaceIntelligence retrieves cached place intelligence
func (cs *CacheService) GetCachedPlaceIntelligence(placeID string) (*models.PlaceIntelligence, error) {
	cacheKey := PlaceIntelligencePrefix + placeID
	
	data, err := cs.client.Get(cs.ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached place intelligence: %w", err)
	}

	var intelligence models.PlaceIntelligence
	if err := json.Unmarshal([]byte(data), &intelligence); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached place intelligence: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"key":      cacheKey,
		"place_id": placeID,
	}).Debug("Retrieved cached place intelligence")

	return &intelligence, nil
}

// CachePlaceDetails caches Foursquare place details
func (cs *CacheService) CachePlaceDetails(placeID string, details interface{}) error {
	cacheKey := PlaceDetailsPrefix + placeID
	
	data, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal place details: %w", err)
	}

	err = cs.client.Set(cs.ctx, cacheKey, data, DetailsCacheTTL).Err()
	if err != nil {
		logrus.WithError(err).WithField("key", cacheKey).Error("Failed to cache place details")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"key":      cacheKey,
		"place_id": placeID,
		"ttl":      DetailsCacheTTL,
	}).Debug("Cached place details")

	return nil
}

// GetCachedPlaceDetails retrieves cached place details
func (cs *CacheService) GetCachedPlaceDetails(placeID string) (interface{}, error) {
	cacheKey := PlaceDetailsPrefix + placeID
	
	data, err := cs.client.Get(cs.ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached place details: %w", err)
	}

	var details interface{}
	if err := json.Unmarshal([]byte(data), &details); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached place details: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"key":      cacheKey,
		"place_id": placeID,
	}).Debug("Retrieved cached place details")

	return details, nil
}

// CachePopularPlaces caches popular places data
func (cs *CacheService) CachePopularPlaces(location string, places []models.PlaceIntelligence) error {
	cacheKey := PopularPlacesPrefix + location
	
	data, err := json.Marshal(places)
	if err != nil {
		return fmt.Errorf("failed to marshal popular places: %w", err)
	}

	err = cs.client.Set(cs.ctx, cacheKey, data, PopularCacheTTL).Err()
	if err != nil {
		logrus.WithError(err).WithField("key", cacheKey).Error("Failed to cache popular places")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"key":      cacheKey,
		"location": location,
		"places":   len(places),
		"ttl":      PopularCacheTTL,
	}).Debug("Cached popular places")

	return nil
}

// GetCachedPopularPlaces retrieves cached popular places
func (cs *CacheService) GetCachedPopularPlaces(location string) ([]models.PlaceIntelligence, error) {
	cacheKey := PopularPlacesPrefix + location
	
	data, err := cs.client.Get(cs.ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached popular places: %w", err)
	}

	var places []models.PlaceIntelligence
	if err := json.Unmarshal([]byte(data), &places); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached popular places: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"key":      cacheKey,
		"location": location,
		"places":   len(places),
	}).Debug("Retrieved cached popular places")

	return places, nil
}

// InvalidateCache removes cached data for a specific key pattern
func (cs *CacheService) InvalidateCache(pattern string) error {
	keys, err := cs.client.Keys(cs.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	if len(keys) == 0 {
		return nil // No keys to delete
	}

	err = cs.client.Del(cs.ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"pattern": pattern,
		"deleted": len(keys),
	}).Info("Invalidated cache keys")

	return nil
}

// GetCacheStats returns cache statistics
func (cs *CacheService) GetCacheStats() (map[string]interface{}, error) {
	info, err := cs.client.Info(cs.ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Get key counts for different prefixes
	searchKeys, _ := cs.client.Keys(cs.ctx, PlaceSearchPrefix+"*").Result()
	intelKeys, _ := cs.client.Keys(cs.ctx, PlaceIntelligencePrefix+"*").Result()
	detailKeys, _ := cs.client.Keys(cs.ctx, PlaceDetailsPrefix+"*").Result()
	popularKeys, _ := cs.client.Keys(cs.ctx, PopularPlacesPrefix+"*").Result()

	stats := map[string]interface{}{
		"redis_info":        info,
		"search_keys":       len(searchKeys),
		"intelligence_keys": len(intelKeys),
		"details_keys":      len(detailKeys),
		"popular_keys":      len(popularKeys),
		"total_keys":        len(searchKeys) + len(intelKeys) + len(detailKeys) + len(popularKeys),
	}

	return stats, nil
}

// HealthCheck verifies Redis connectivity
func (cs *CacheService) HealthCheck() error {
	_, err := cs.client.Ping(cs.ctx).Result()
	if err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (cs *CacheService) Close() error {
	return cs.client.Close()
}