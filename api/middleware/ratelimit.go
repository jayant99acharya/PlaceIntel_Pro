package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds the rate limiter for each client
type RateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ClientRateLimiters holds rate limiters for all clients
type ClientRateLimiters struct {
	clients map[string]*RateLimiter
	mu      sync.RWMutex
}

var clientLimiters = &ClientRateLimiters{
	clients: make(map[string]*RateLimiter),
}

// RateLimit middleware for API rate limiting
func RateLimit() gin.HandlerFunc {
	// Get rate limit configuration from environment
	requestsStr := os.Getenv("RATE_LIMIT_REQUESTS")
	if requestsStr == "" {
		requestsStr = "100"
	}
	requests, _ := strconv.Atoi(requestsStr)

	windowStr := os.Getenv("RATE_LIMIT_WINDOW")
	if windowStr == "" {
		windowStr = "3600"
	}
	window, _ := strconv.Atoi(windowStr)

	// Calculate rate limit (requests per second)
	rps := rate.Limit(float64(requests) / float64(window))
	burst := requests / 10 // Allow burst of 10% of total requests

	return gin.HandlerFunc(func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		// Get or create rate limiter for this client
		limiter := getRateLimiter(clientIP, rps, burst)
		
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// getRateLimiter gets or creates a rate limiter for a client
func getRateLimiter(clientIP string, rps rate.Limit, burst int) *rate.Limiter {
	clientLimiters.mu.Lock()
	defer clientLimiters.mu.Unlock()

	// Clean up old entries (older than 1 hour)
	now := time.Now()
	for ip, rl := range clientLimiters.clients {
		if now.Sub(rl.lastSeen) > time.Hour {
			delete(clientLimiters.clients, ip)
		}
	}

	// Get or create rate limiter for this client
	if rl, exists := clientLimiters.clients[clientIP]; exists {
		rl.lastSeen = now
		return rl.limiter
	}

	// Create new rate limiter
	limiter := rate.NewLimiter(rps, burst)
	clientLimiters.clients[clientIP] = &RateLimiter{
		limiter:  limiter,
		lastSeen: now,
	}

	return limiter
}