package routes

import (
	"restaurant_manager/controllers"
	"restaurant_manager/middleware"
	"restaurant_manager/services"

	"github.com/gin-gonic/gin"
)

func RoomImageRoutes(r *gin.Engine) {
	roomImageService := services.NewRoomImageService()
	roomImageController := controllers.NewRoomImageController(roomImageService)

	r.POST("/room/image", middleware.AuthMiddleware(), middleware.AdminOnly(), roomImageController.Create)
	r.GET("/room/image", middleware.AuthMiddleware(), roomImageController.AllRoomImage)
	r.DELETE("/room/image/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), roomImageController.RemoveRoomImage)
}
