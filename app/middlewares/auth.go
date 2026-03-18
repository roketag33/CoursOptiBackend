package middlewares

import (
	"dungeons/app/models"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Success(http.StatusUnauthorized, "middleware.Auth.Missing", "Token manquant"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Success(http.StatusUnauthorized, "middleware.Auth.InvalidFormat", "Format Bearer <token> attendu"))
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("TOKEN_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Success(http.StatusUnauthorized, "middleware.Auth.Invalid", "Token invalide ou expiré"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Success(http.StatusUnauthorized, "middleware.Auth.InvalidClaims", "Claims invalides"))
			return
		}

		c.Set("playerID", claims["customID"])
		c.Next()
	}
}
