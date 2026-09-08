package auth

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth(apiToken string) gin.HandlerFunc {
	return func(c *gin.Context)  {
		auth := c.GetHeader("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		if subtle.ConstantTimeCompare([]byte(token), []byte(apiToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"})
			return 
		}
		c.Next()
	}
}