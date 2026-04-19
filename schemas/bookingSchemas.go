package schemas

import "time"

type BookingSchemas struct {
	RoomID   uint      `json:"room_id"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
}

type UpdateBookingStatus struct {
	Status string `json:"status"`
}
