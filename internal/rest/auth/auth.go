package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer valid_token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
			c.Abort()
			return
		}
		c.Next()
	}
}
