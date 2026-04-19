package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	HotelID       uint    `json:"hotel_id"`
	Hotel         Hotel   `json:"hotel" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoomNumber    string  `json:"room_number" gorm:"size:255"`
	RoomType      string  `json:"room_type" gorm:"size:150"`
	PricePerNight float64 `json:"price_per_night" gorm:"type:decimal(10,2);not null"`
	Capacity      int     `json:"capacity" gorm:"not null"`
	Description   string  `json:"description" gorm:"type:text"`
	Status        string  `json:"status" gorm:"size:100;default:'available'"`

	RoomImages []RoomImage `json:"room_image" gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
