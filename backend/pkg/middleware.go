package pkg

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			JSONUnauthorized(c, ErrMissingAuth)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			JSONUnauthorized(c, ErrInvalidToken)
			c.Abort()
			return
		}

		token, err := VerifyJWT(tokenString)
		if err != nil || !token.Valid {
			JSONUnauthorized(c, ErrInvalidToken)
			c.Abort()
			return
		}

		claims, err := ExtractClaims(token)
		if err != nil {
			JSONUnauthorized(c, ErrInvalidToken)
			c.Abort()
			return
		}

		userIDFloat, _ := claims["user_id"].(float64)
		email, _ := claims["email"].(string)

		c.Set("userID", uint(userIDFloat))
		c.Set("email", email)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetEmail extracts email from context
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	emailStr, ok := email.(string)
	return emailStr, ok
}
