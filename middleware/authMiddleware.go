package middleware

import (
	"github.com/eaaslan/go-jwt/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")

		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no authorization"})
			c.Abort()
			return
		}
		claims, err := helpers.ValidateToken(clientToken)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.UID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
