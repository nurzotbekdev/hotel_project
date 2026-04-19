package schemas

type RoomSchemas struct {
	HotelID       uint    `json:"hotel_id"`
	RoomNumber    string  `json:"room_number"`
	RoomType      string  `json:"room_type"`
	PricePerNight float64 `json:"price_per_night"`
	Capacity      int     `json:"capacity"`
	Description   string  `json:"description"`
}

type StatusUpdateRoom struct {
	Status string `json:"status"`
}

type RoomSchemasUpdate struct {
	HotelID       *uint    `json:"hotel_id"`
	RoomNumber    *string  `json:"room_number"`
	RoomType      *string  `json:"room_type"`
	PricePerNight *float64 `json:"price_per_night"`
	Capacity      *int     `json:"capacity"`
	Description   *string  `json:"description"`
}
