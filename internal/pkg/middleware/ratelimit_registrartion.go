package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	mu      sync.Mutex
	lastReq time.Time
	count   int
}

var visitors = struct {
	sync.Mutex
	m map[string]*visitor
}{m: make(map[string]*visitor)}

func RateLimitMiddleware(maxReq int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		visitors.Lock()
		v, exists := visitors.m[ip]
		if !exists {
			v = &visitor{lastReq: time.Now(), count: 0}
			visitors.m[ip] = v
		}
		visitors.Unlock()

		v.mu.Lock()
		defer v.mu.Unlock()

		// сдвигаем окно
		if time.Since(v.lastReq) > window {
			v.count = 0
			v.lastReq = time.Now()
		}

		v.count++
		if v.count > maxReq {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Слишком много запросов, повторите позже",
			})
			return
		}

		c.Next()
	}
}
