package services

import (
	"errors"
	"math"
	"restaurant_manager/config"
	"restaurant_manager/models"

	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(booking models.Booking) error
	GetBookingMy(userID uint) ([]models.Booking, error)
	GetBookingByID(userID, ID uint) (models.Booking, error)
	GetAllBooking() ([]models.Booking, error)
	GetAdminBookingByID(ID uint) (models.Booking, error)
	EditStatusBooking(ID uint, status string) error
	DeleteBooking(userID, ID uint) error
}

type bookingService struct{}

func NewBookingServices() BookingService {
	return &bookingService{}
}

var (
	BookingNotFound  = errors.New("Booking not found")
	RoomNotAvailable = errors.New("Room is not available")
)

func (s *bookingService) CreateBooking(booking models.Booking) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var room models.Room
		if err := tx.First(&room, booking.RoomID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return RoomNotFound
			}
			return err
		}

		if room.Status != "available" {
			return RoomNotAvailable
		}

		hours := booking.CheckOut.Sub(booking.CheckIn).Hours()
		days := int(math.Ceil(hours / 24))
		if days <= 0 {
			return errors.New("invalid booking dates")
		}

		booking.TotalPrice = float64(days) * room.PricePerNight
		booking.Status = "pending"

		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		if err := tx.Model(&room).Update("status", "booked").Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *bookingService) GetBookingMy(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking

	err := config.DB.
		Preload("Room").
		Preload("Room.RoomImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("Room.Hotel").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error

	return bookings, err
}

func (s *bookingService) GetBookingByID(userID, ID uint) (models.Booking, error) {
	var booking models.Booking

	err := config.DB.
		Preload("Room").
		Preload("Room.RoomImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Preload("Room.Hotel").
		Where("user_id = ? AND id  = ?", userID, ID).
		Order("created_at DESC").
		First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return booking, BookingNotFound
		}
		return booking, err
	}

	return booking, nil
}

func (s *bookingService) GetAllBooking() ([]models.Booking, error) {
	var bookings []models.Booking

	err := config.DB.
		Preload("User").
		Preload("Room").
		Preload("Room.RoomImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Room.Hotel").
		Order("created_at DESC").
		Find(&bookings).Error

	return bookings, err
}

func (s *bookingService) GetAdminBookingByID(ID uint) (models.Booking, error) {
	var booking models.Booking

	err := config.DB.
		Preload("User").
		Preload("Room").
		Preload("Room.RoomImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Room.Hotel").
		Where("id = ?", ID).
		Order("created_at DESC").
		First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return booking, BookingNotFound
		}
		return booking, err
	}

	return booking, nil
}

func (s *bookingService) EditStatusBooking(ID uint, status string) error {
	var booking models.Booking
	if err := config.DB.First(&booking, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BookingNotFound
		}
		return err
	}
	if err := config.DB.Model(&booking).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (s *bookingService) DeleteBooking(userID, bookingID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking

		if err := tx.First(&booking, bookingID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return BookingNotFound
			}
			return err
		}

		if booking.UserID != userID {
			return BookingNotFound
		}

		if err := tx.Model(&models.Room{}).
			Where("id = ?", booking.RoomID).
			Update("status", "available").Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&booking).Error; err != nil {
			return err
		}

		return nil
	})
}
