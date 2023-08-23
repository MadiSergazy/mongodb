package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	token "github.com/MadiSergazy/mongodb/GolangEcommerceProject/tokens"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authentication header"})
			c.Abort()
			return
		}

		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
