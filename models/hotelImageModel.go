package models

import "gorm.io/gorm"

type HotelImage struct {
	gorm.Model
	HotelID    uint   `json:"hotel_id"`
	Hotel      Hotel  `json:"hotel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	HotelImage string `json:"hotel_image" gorm:"size:255"`
}
