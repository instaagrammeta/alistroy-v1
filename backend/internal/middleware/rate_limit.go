package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit is a per-IP token bucket limiter with periodic GC of stale entries.
func RateLimit(rps, burst int) gin.HandlerFunc {
	var mu sync.Mutex
	limiters := make(map[string]*rate.Limiter)
	seen := make(map[string]time.Time)

	get := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()
		l, ok := limiters[ip]
		if !ok {
			l = rate.NewLimiter(rate.Limit(rps), burst)
			limiters[ip] = l
		}
		seen[ip] = time.Now()
		if len(limiters) > 10000 {
			cut := time.Now().Add(-30 * time.Minute)
			for k, t := range seen {
				if t.Before(cut) {
					delete(limiters, k)
					delete(seen, k)
				}
			}
		}
		return l
	}

	return func(c *gin.Context) {
		if !get(c.ClientIP()).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{"code": "rate_limited", "message": "too many requests"},
			})
			return
		}
		c.Next()
	}
}
