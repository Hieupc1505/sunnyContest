package middleware

import (
	"fmt"
	app "go-rest-api-boilerplate/internal"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter defines a per-IP rate limiter with whitelist/blacklist support
type RateLimiter struct {
	limiters   map[string]*rate.Limiter
	lastAccess map[string]time.Time
	whitelist  map[string]bool
	blacklist  map[string]bool
	mu         sync.Mutex
	rps        rate.Limit
	burst      int
}

// NewRateLimiter initializes the rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters:   make(map[string]*rate.Limiter),
		lastAccess: make(map[string]time.Time),
		whitelist:  make(map[string]bool),
		blacklist:  make(map[string]bool),
		rps:        rate.Limit(rps),
		burst:      burst,
	}

	// Start periodic cleanup
	go rl.cleanupInactiveLimiters()
	return rl
}

// getLimiter returns the rate limiter for a specific IP
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rps, rl.burst)
		rl.limiters[ip] = limiter
	}
	// Update last access timestamp
	rl.lastAccess[ip] = time.Now()
	return limiter
}

// cleanupInactiveLimiters removes limiters that haven't been used for 1 minute
func (rl *RateLimiter) cleanupInactiveLimiters() {
	for {
		time.Sleep(1 * time.Minute)
		rl.mu.Lock()
		cutoff := time.Now().Add(-1 * time.Minute)
		for ip, lastSeen := range rl.lastAccess {
			if lastSeen.Before(cutoff) {
				delete(rl.limiters, ip)
				delete(rl.lastAccess, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// extractIP safely retrieves the real IP from headers or RemoteAddr
func extractIP(c *gin.Context) string {
	// Check X-Forwarded-For header (comma-separated list)
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0]) // Return first IP
	}

	// Check X-Real-IP (alternative)
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Default: use ClientIP from Gin (handles proxy and direct requests)
	return c.ClientIP()
}

// RateLimitMiddleware applies rate limiting per IP
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := extractIP(c)
		fmt.Println("ip: ", ip)
		// Handle blacklist
		if rl.IsBlacklisted(ip) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"e": app.ErrForbidden})
			return
		}

		// Allow whitelisted IPs
		if rl.IsWhitelisted(ip) {
			c.Next()
			return
		}

		// Enforce rate limit
		limiter := rl.getLimiter(ip)
		slog.Info("limiter: ", limiter, !limiter.Allow())
		if !limiter.Allow() {
			log.Printf("Rate limit exceeded for IP: %s", ip)
			c.AbortWithStatusJSON(429, gin.H{"e": app.ErrTooManyRequests})
			return
		}

		c.Next()
	}
}

// IsWhitelisted checks if an IP is whitelisted
func (rl *RateLimiter) IsWhitelisted(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	_, exists := rl.whitelist[ip]
	return exists
}

// IsBlacklisted checks if an IP is blacklisted
func (rl *RateLimiter) IsBlacklisted(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	_, exists := rl.blacklist[ip]
	return exists
}

// AddWhitelist adds an IP to the whitelist
func (rl *RateLimiter) AddWhitelist(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.whitelist[ip] = true
}

// AddBlacklist adds an IP to the blacklist
func (rl *RateLimiter) AddBlacklist(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.blacklist[ip] = true
}
