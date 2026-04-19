package routes

import (
	"restaurant_manager/controllers"
	"restaurant_manager/middleware"
	"restaurant_manager/services"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userService := services.NewUserServices()
	userController := controllers.NewUserController(userService)

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.POST("/refresh", userController.RefreshToken)
	r.GET("/profile", middleware.AuthMiddleware(), userController.MyProfile)
	r.PATCH("/update/:id/role", middleware.AuthMiddleware(), middleware.AdminOnly(), userController.UpdateRole)
	r.PUT("/profile", middleware.AuthMiddleware(), userController.UpdateProfile)
	r.DELETE("/profile", middleware.AuthMiddleware(), userController.DeleteAccount)
}
