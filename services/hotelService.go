package services

import (
	"errors"
	"restaurant_manager/config"
	"restaurant_manager/models"
	"restaurant_manager/schemas"
	"strings"

	"gorm.io/gorm"
)

type HotelService interface {
	CreateHotel(hotel models.Hotel) error
	GetHotelAll() ([]models.Hotel, error)
	GetHotelByID(ID uint) (models.Hotel, error)
	SetHotel(ID uint, data schemas.EditHotel) error
	DeleteHotel(ID uint) error
}

type hotelService struct{}

func NewHotelServices() HotelService {
	return &hotelService{}
}

func (s *hotelService) CreateHotel(hotel models.Hotel) error {
	return config.DB.Create(&hotel).Error
}

func (s *hotelService) GetHotelAll() ([]models.Hotel, error) {
	var hotels []models.Hotel
	err := config.DB.
		Preload("HotelImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Find(&hotels).Error
	return hotels, err
}

func (s *hotelService) GetHotelByID(ID uint) (models.Hotel, error) {
	var hotel models.Hotel
	err := config.DB.
		Preload("HotelImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("id =?", ID).First(&hotel).Error
	return hotel, err
}

func (s *hotelService) SetHotel(ID uint, data schemas.EditHotel) error {
	var hotel models.Hotel

	if err := config.DB.First(&hotel, ID).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}

	if data.Name != nil {
		updates["name"] = strings.TrimSpace(*data.Name)
	}
	if data.Address != nil {
		updates["address"] = strings.TrimSpace(*data.Address)
	}
	if data.Description != nil {
		updates["description"] = strings.TrimSpace(*data.Description)
	}
	if data.Phone != nil {
		updates["phone"] = strings.TrimSpace(*data.Phone)
	}
	if data.Email != nil {
		updates["email"] = strings.TrimSpace(*data.Email)
	}

	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	if err := config.DB.Model(&hotel).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (s *hotelService) DeleteHotel(ID uint) error {
	var hotel models.Hotel
	if err := config.DB.First(&hotel, ID).Error; err != nil {
		return err
	}
	return config.DB.Unscoped().Delete(&hotel).Error
}
