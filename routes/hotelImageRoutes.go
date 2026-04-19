package routes

import (
	"restaurant_manager/controllers"
	"restaurant_manager/middleware"
	"restaurant_manager/services"

	"github.com/gin-gonic/gin"
)

func HotelImageRoutes(r *gin.Engine) {
	hotelImageService := services.NewHotelImageService()
	hotelImageController := controllers.NewHotelImageController(hotelImageService)

	r.POST("/hotel/image", middleware.AuthMiddleware(), middleware.AdminOnly(), hotelImageController.Create)
	r.GET("/hotel/image", middleware.AuthMiddleware(), hotelImageController.AllHotelImage)
	r.DELETE("/hotel/image/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), hotelImageController.RemoveHotelImage)
}
