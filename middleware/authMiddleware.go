package middleware

import (
	"net/http"
	"os"
	"restaurant_manager/config"
	"restaurant_manager/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("Authorization")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalid",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["type"] != "access" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not an access token",
			})
		}

		userID := uint(claims["sub"].(float64))

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}