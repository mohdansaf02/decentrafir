package middleware

import (
	"net/http"
	"strings"

	"github.com/decntrafir/backend/internal/models"
	"github.com/decntrafir/backend/pkg/jwtutil"
	"github.com/gin-gonic/gin"
)

const ContextUserKey = "userClaims"

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwtutil.ParseAccessToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(ContextUserKey, claims)
		c.Next()
	}
}

func RoleMiddleware(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get(ContextUserKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userClaims := claims.(*jwtutil.Claims)
		for _, r := range roles {
			if userClaims.Role == string(r) {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
