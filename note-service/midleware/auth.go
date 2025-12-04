package midleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"myproject/auth" // <-- ПОДСТАВЬ свой module path из note-service/go.mod
)

// ключ, под которым будем класть user_id в контекст
const userIDContextKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}
		tokenStr := parts[1]

		userID, err := auth.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// кладём userID в контекст
		c.Set(userIDContextKey, userID)

		c.Next()
	}
}

// Хелпер для хендлеров
func GetUserID(c *gin.Context) (int, bool) {
	v, ok := c.Get(userIDContextKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}
