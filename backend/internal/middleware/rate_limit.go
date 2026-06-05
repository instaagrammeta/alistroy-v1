package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit is a simple per-IP token-bucket rate limiter.
func RateLimit(rps int, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	last := make(map[string]time.Time)
	var mu sync.Mutex

	get := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		l, ok := limiters[ip]
		if !ok {
			l = rate.NewLimiter(rate.Limit(rps), burst)
			limiters[ip] = l
		}
		last[ip] = time.Now()
		// occasional GC of old entries
		if len(limiters) > 5000 {
			cutoff := time.Now().Add(-30 * time.Minute)
			for k, t := range last {
				if t.Before(cutoff) {
					delete(limiters, k)
					delete(last, k)
				}
			}
		}
		return l
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !get(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{"code": "rate_limited", "message": "too many requests"},
			})
			return
		}
		c.Next()
	}
}
