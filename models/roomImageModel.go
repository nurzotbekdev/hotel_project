package models

import "gorm.io/gorm"

type RoomImage struct {
	gorm.Model
	RoomID   uint   `json:"room_id"`
	Room     Room   `json:"room" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageURL string `json:"image_url" gorm:"size:255"`
}
