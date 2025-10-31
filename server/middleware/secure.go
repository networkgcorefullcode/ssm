package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
)

// RateLimiter manages rate limiting for clients
type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*ClientLimiter
	config  *factory.RateLimit
}

// ClientLimiter tracks individual client rate limiting
type ClientLimiter struct {
	tokens    int
	lastReset time.Time
}

var rateLimiter *RateLimiter

// InitRateLimiter initializes the rate limiter with configuration
func InitRateLimiter(config *factory.RateLimit) {
	if config == nil || !config.Enabled {
		logger.AppLog.Info("Rate limiting is disabled")
		return
	}

	rateLimiter = &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		config:  config,
	}

	// Start cleanup routine
	go rateLimiter.cleanup()

	logger.AppLog.Infof("Rate limiter initialized: %d requests/min, burst: %d",
		config.RequestsPerMin, config.BurstSize)
}

// cleanup removes old client entries periodically
func (rl *RateLimiter) cleanup() {
	cleanupInterval := time.Duration(rl.config.CleanupInterval) * time.Minute
	if cleanupInterval <= 0 {
		cleanupInterval = 10 * time.Minute // default 10 minutes
	}

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for clientIP, client := range rl.clients {
			// Remove clients inactive for more than 2 cleanup intervals
			if now.Sub(client.lastReset) > cleanupInterval*2 {
				delete(rl.clients, clientIP)
			}
		}
		rl.mu.Unlock()
		logger.AppLog.Debug("Rate limiter cleanup completed")
	}
}

// isAllowed checks if a client is allowed to make a request
func (rl *RateLimiter) isAllowed(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[clientIP]

	if !exists {
		// New client
		rl.clients[clientIP] = &ClientLimiter{
			tokens:    rl.config.RequestsPerMin - 1,
			lastReset: now,
		}
		return true
	}

	// Check if we need to reset the token bucket (1 minute window)
	if now.Sub(client.lastReset) >= time.Minute {
		client.tokens = rl.config.RequestsPerMin
		client.lastReset = now
	}

	// Check if client has tokens available
	if client.tokens > 0 {
		client.tokens--
		return true
	}

	// Check burst allowance
	if rl.config.BurstSize > 0 && client.tokens > -rl.config.BurstSize {
		client.tokens--
		return true
	}

	return false
}

// getRateLimitHeaders returns rate limit headers
func (rl *RateLimiter) getRateLimitHeaders(clientIP string) (remaining int, resetTime int64) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	client, exists := rl.clients[clientIP]
	if !exists {
		return rl.config.RequestsPerMin, time.Now().Add(time.Minute).Unix()
	}

	remaining = client.tokens
	if remaining < 0 {
		remaining = 0
	}

	resetTime = client.lastReset.Add(time.Minute).Unix()
	return remaining, resetTime
}

// secureRequest adds security features to the request, for example DoS protection, rate limiting, etc.
func SecureRequest(c *gin.Context) {
	// Apply rate limiting if enabled
	if rateLimiter != nil && rateLimiter.config.Enabled {
		clientIP := c.ClientIP()

		// Check rate limit
		if !rateLimiter.isAllowed(clientIP) {
			logger.AppLog.Warnf("Rate limit exceeded for IP: %s", clientIP)

			// Add rate limit headers
			remaining, resetTime := rateLimiter.getRateLimitHeaders(clientIP)
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rateLimiter.config.RequestsPerMin))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     "Too many requests, please try again later",
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		// Add rate limit headers for successful requests
		remaining, resetTime := rateLimiter.getRateLimitHeaders(clientIP)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rateLimiter.config.RequestsPerMin))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime))
	}

	// Call the next middleware or endpoint handler
	c.Next()
}
