package services

import (
	"errors"
	"restaurant_manager/config"
	"restaurant_manager/models"

	"gorm.io/gorm"
)

type RoomImageService interface {
	CreateRoomImage(image models.RoomImage) error
	GetAllRoomImage() ([]models.RoomImage, error)
	DeleteRoomImage(ID uint) error
}

type roomImageService struct{}

func NewRoomImageService() RoomImageService {
	return &roomImageService{}
}

var (
	RoomImageNotFound = errors.New("Room image not fonund")
)

func (s *roomImageService) CreateRoomImage(image models.RoomImage) error {
	var room models.Room
	if err := config.DB.First(&room, image.RoomID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoomNotFound
		}
		return err
	}

	return config.DB.Create(&image).Error
}

func (s *roomImageService) GetAllRoomImage() ([]models.RoomImage, error) {
	var image []models.RoomImage
	err := config.DB.Order("created_at DESC").Find(&image).Error
	return image, err
}

func (s *roomImageService) DeleteRoomImage(ID uint) error {
	var roomImage models.RoomImage
	if err := config.DB.First(&roomImage, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoomImageNotFound
		}
		return err
	}

	return config.DB.Unscoped().Delete(&roomImage).Error
}
