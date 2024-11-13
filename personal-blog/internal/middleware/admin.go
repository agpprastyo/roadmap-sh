package middleware

import "github.com/gin-gonic/gin"

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
