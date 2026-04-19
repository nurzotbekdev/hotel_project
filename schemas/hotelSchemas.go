package schemas

type HotelSchemas struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}

type EditHotel struct {
	Name        *string `json:"name"`
	Address     *string `json:"address"`
	Description *string `json:"description"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email"`
}
