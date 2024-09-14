package middlewares

import (
	"fmt"
	"net/http"

	helper "github.com/datarohit/go-restaurant-management-backend-project/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(clientToken, &helper.SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(helper.JWT_SECRET), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*helper.SignedDetails)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("firstName", claims.FirstName)
		c.Set("lastName", claims.LastName)
		c.Set("uid", claims.UID)

		c.Next()
	}
}
