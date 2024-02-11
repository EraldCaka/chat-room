package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/signup" || c.Request.URL.Path == "/login" {
			c.Next()
			return
		}
		_, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		//if !isValidToken(cookie) {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		//	c.Abort()
		//	return
		//}

		c.Next()
	}
}
