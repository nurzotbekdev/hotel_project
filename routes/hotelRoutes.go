package routes

import (
	"restaurant_manager/controllers"
	"restaurant_manager/middleware"
	"restaurant_manager/services"

	"github.com/gin-gonic/gin"
)

func HotelRoutes(r *gin.Engine) {
	hotelService := services.NewHotelServices()
	hotelController := controllers.NewHotelController(hotelService)

	r.POST("/hotel", middleware.AuthMiddleware(), middleware.AdminOnly(), hotelController.Create)
	r.GET("/hotel", middleware.AuthMiddleware(), hotelController.AllHotel)
	r.GET("/hotel/:id", middleware.AuthMiddleware(), hotelController.HotelDetail)
	r.PUT("/hotel", middleware.AuthMiddleware(), middleware.AdminOnly(), hotelController.UpdateHotel)
	r.DELETE("/hotel", middleware.AuthMiddleware(), middleware.AdminOnly(), hotelController.DeleteHotel)
}
