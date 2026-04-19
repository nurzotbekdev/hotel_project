package main

import (
	"restaurant_manager/config"
	"restaurant_manager/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.EnvConfig()
	config.ConnectingDatabase()
	config.MigrateConfig()
	config.ConnectingRedis()
}

func main() {
	router := gin.Default()

	routes.UserRoutes(router)
	routes.HotelRoutes(router)
	routes.HotelImageRoutes(router)
	routes.RoomRoutes(router)
	routes.RoomImageRoutes(router)
	routes.BookingRoutes(router)

	router.Run()
}
