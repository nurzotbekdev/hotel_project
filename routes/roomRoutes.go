package routes

import (
	"restaurant_manager/controllers"
	"restaurant_manager/middleware"
	"restaurant_manager/services"

	"github.com/gin-gonic/gin"
)

func RoomRoutes(r *gin.Engine) {
	roomService := services.NewRoomServices()
	roomController := controllers.NewRoomController(roomService)

	r.POST("/room", middleware.AuthMiddleware(), middleware.AdminOnly(), roomController.Create)
	r.GET("/room", middleware.AuthMiddleware(), roomController.AllRoom)
	r.GET("/room/:id", middleware.AuthMiddleware(), roomController.RoomDetail)
	r.GET("admin/room", middleware.AuthMiddleware(), middleware.AdminOnly(), roomController.AdminRoomList)
	r.PUT("/room/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), roomController.UpdateRoom)
	r.PUT("/room/:id/status", middleware.AuthMiddleware(), middleware.AdminOnly(), roomController.UpdateStatus)
	r.DELETE("/room/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), roomController.RemoveRoom)
}
