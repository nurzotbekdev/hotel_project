package middleware

import (
	"net/http"
	"restaurant_manager/models"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInterface, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			ctx.Abort()
			return
		}

		user := userInterface.(models.User)

		if user.Role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "You do not have administrator rights",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}