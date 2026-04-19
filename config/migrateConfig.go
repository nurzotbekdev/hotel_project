package config

import "restaurant_manager/models"

func MigrateConfig() {
	DB.AutoMigrate(
		&models.User{},
		&models.Hotel{},
		&models.HotelImage{},
		&models.Room{},
		&models.RoomImage{},
		&models.Booking{},
	)
}
