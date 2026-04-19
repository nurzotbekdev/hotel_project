package services

import (
	"errors"
	"restaurant_manager/config"
	"restaurant_manager/models"

	"gorm.io/gorm"
)

type HotelImageService interface {
	CreateHotelImage(hotelImage models.HotelImage) error
	GetAllHotelImage() ([]models.HotelImage, error)
	DeleteHotelImage(ID uint) error
}

type hotelImageService struct{}

func NewHotelImageService() HotelImageService {
	return &hotelImageService{}
}

var (
	HotelImageNotFound = errors.New("Hotel image not found")
)

func (s *hotelImageService) CreateHotelImage(hotelImage models.HotelImage) error {
	var hotel models.Hotel
	if err := config.DB.First(&hotel, hotelImage.HotelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return HotelNotFound
		}
		return err
	}
	return config.DB.Create(&hotelImage).Error
}

func (s *hotelImageService) GetAllHotelImage() ([]models.HotelImage, error) {
	var hotelImage []models.HotelImage
	err := config.DB.Order("created_at DESC").Find(&hotelImage).Error
	return hotelImage, err
}

func (s *hotelImageService) DeleteHotelImage(ID uint) error {
	var hotelImage models.HotelImage
	if err := config.DB.First(&hotelImage, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return HotelImageNotFound
		}
		return err
	}

	return config.DB.Unscoped().Delete(&hotelImage).Error
}
