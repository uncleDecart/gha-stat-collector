package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/pkg/token"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.GetHeader("auth")
		if t == "" {
			c.String(http.StatusBadRequest, "auth is not provided")
			c.Abort()
			return
		}
		err := token.TokenValid(t)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
