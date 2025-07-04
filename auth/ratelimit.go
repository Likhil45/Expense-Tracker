package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

// Extract user_id from JWT context (set by your JWT middleware)
func RateLimitMiddleware() gin.HandlerFunc {
	// 10 requests per minute
	rate, _ := limiter.NewRateFromFormatted("10-M")
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		identifier := ""
		if exists {
			identifier = userID.(string)
		} else {
			// fallback to IP if no user_id (should not happen for protected routes)
			identifier = strings.Split(c.ClientIP(), ":")[0]
		}

		context, err := instance.Get(c, identifier)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%v", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%v", context.Reset))

		if context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
