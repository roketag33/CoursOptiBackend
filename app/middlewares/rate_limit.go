package middlewares

import (
	"dungeons/app/models"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Count     int
	ExpiresAt time.Time
}

var (
	clients = sync.Map{}
	limit   = 100 // max 100 req par minute
)

func RateLimitMiddleware() gin.HandlerFunc {
	// Cleanup background routine
	go func() {
		for {
			time.Sleep(time.Minute)
			now := time.Now()
			clients.Range(func(key, value interface{}) bool {
				client := value.(*Client)
				if now.After(client.ExpiresAt) {
					clients.Delete(key)
				}
				return true
			})
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		val, exists := clients.Load(ip)
		if !exists {
			clients.Store(ip, &Client{
				Count:     1,
				ExpiresAt: now.Add(time.Minute),
			})
			c.Next()
			return
		}

		client := val.(*Client)

		if now.After(client.ExpiresAt) {
			client.Count = 1
			client.ExpiresAt = now.Add(time.Minute)
			clients.Store(ip, client)
			c.Next()
			return
		}

		if client.Count >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, models.Success(http.StatusTooManyRequests, "middleware.RateLimit.Exceeded", "Trop de requêtes, calme-toi !"))
			return
		}

		client.Count++
		clients.Store(ip, client)
		c.Next()
	}
}
