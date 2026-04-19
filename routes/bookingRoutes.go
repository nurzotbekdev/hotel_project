package routes

import (
	"restaurant_manager/controllers"
	"restaurant_manager/middleware"
	"restaurant_manager/services"

	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.Engine) {
	bookingService := services.NewBookingServices()
	bookingController := controllers.NewBookingController(bookingService)

	r.POST("/booking", middleware.AuthMiddleware(), bookingController.Create)
	r.GET("/my/booking", middleware.AuthMiddleware(), bookingController.MyBookings)
	r.GET("/booking/:id", middleware.AuthMiddleware(), bookingController.MyBookingDetail)
	r.GET("admin/booking", middleware.AuthMiddleware(), middleware.AdminOnly(), bookingController.AdminBookingList)
	r.GET("admin/booking/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), bookingController.AdminBookingByID)
	r.PUT("/booking/:id/status", middleware.AuthMiddleware(), middleware.AdminOnly(), bookingController.UpdateStatus)
	r.DELETE("/booking/:id", middleware.AuthMiddleware(), bookingController.RemoveBooking)
}
