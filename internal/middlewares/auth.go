package middlewares

import (
	"go-auth-postgres/internal/auth"
	"go-auth-postgres/internal/repositories"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userRepository repositories.UsersRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the access token from the cookie
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
			return
		}

		secret := os.Getenv("ACCESS_TOKEN_SECRET")

		_, claims, err := auth.ValidateToken(accessToken, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			return
		}

		// Verify that the user exists
		user, err := userRepository.FindByID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("userID", user)
		c.Next()
	}
}
