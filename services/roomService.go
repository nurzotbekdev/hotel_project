package services

import (
	"errors"
	"restaurant_manager/config"
	"restaurant_manager/models"
	"restaurant_manager/schemas"

	"gorm.io/gorm"
)

type RoomService interface {
	CreateRoom(room models.Room) error
	GetRoomAll() ([]models.Room, error)
	GetRoomByID(ID uint) (models.Room, error)
	AdminGetRoomAll() ([]models.Room, error)
	EditRoom(ID uint, data schemas.RoomSchemasUpdate) error
	EditStatusRoom(ID uint, status string) error
	DeleteRoom(ID uint) error
}

type roomService struct{}

func NewRoomServices() RoomService {
	return &roomService{}
}

var (
	HotelNotFound = errors.New("Hotel not found")
	RoomNotFound  = errors.New("Room not found")
)

func (s *roomService) CreateRoom(room models.Room) error {
	var hotel models.Hotel
	if err := config.DB.First(&hotel, room.HotelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return HotelNotFound
		}
		return err
	}

	return config.DB.Create(&room).Error
}

func (s *roomService) GetRoomAll() ([]models.Room, error) {
	var rooms []models.Room

	err := config.DB.
		Preload("Hotel").
		Preload("RoomImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("status =?", "available").
		Find(&rooms).Error

	return rooms, err
}

func (s *roomService) GetRoomByID(ID uint) (models.Room, error) {
	var room models.Room

	err := config.DB.
		Preload("Hotel").
		Preload("RoomImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("id = ? AND status = ?", ID, "available").
		First(&room).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return room, RoomNotFound
		}
		return room, err
	}

	return room, nil
}

func (s *roomService) AdminGetRoomAll() ([]models.Room, error) {
	var rooms []models.Room

	err := config.DB.
		Order("rooms.created_at DESC").
		Preload("Hotel").
		Preload("RoomImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Find(&rooms).Error

	return rooms, err
}

func (s *roomService) EditRoom(ID uint, data schemas.RoomSchemasUpdate) error {
	var room models.Room
	if err := config.DB.First(&room, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoomNotFound
		}
		return err
	}

	updates := map[string]interface{}{}

	if data.HotelID != nil && *data.HotelID != room.HotelID {
		var hotel models.Hotel
		if err := config.DB.First(&hotel, *data.HotelID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return HotelNotFound
			}
			return err
		}
		updates["hotel_id"] = *data.HotelID
	}

	if data.RoomNumber != nil {
		updates["room_number"] = *data.RoomNumber
	}
	if data.RoomType != nil {
		updates["room_type"] = *data.RoomType
	}
	if data.PricePerNight != nil {
		updates["price_per_night"] = *data.PricePerNight
	}
	if data.Capacity != nil {
		updates["capacity"] = *data.Capacity
	}
	if data.Description != nil {
		updates["description"] = *data.Description
	}

	if len(updates) == 0 {
		return errors.New("No fields to update")
	}

	return config.DB.Model(&room).Updates(updates).Error
}

func (s *roomService) EditStatusRoom(ID uint, status string) error {
	var room models.Room
	if err := config.DB.First(&room, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoomNotFound
		}
		return err
	}

	if err := config.DB.Model(&room).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func (s *roomService) DeleteRoom(ID uint) error {
	var room models.Room
	if err := config.DB.First(&room, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoomNotFound
		}
		return err
	}

	return config.DB.Unscoped().Delete(&room).Error
}
