package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RoomID     uint      `json:"room_id"`
	Room       Room      `json:"room" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	TotalPrice float64   `json:"total_price" gorm:"type:decimal(10,2)"`
	Status     string    `json:"status" gorm:"size:30;default:'pending'"`
}
