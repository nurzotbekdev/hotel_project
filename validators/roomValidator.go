package validators

import (
	"errors"
	"restaurant_manager/schemas"
	"strings"
)

func ValidateRoom(req *schemas.RoomSchemas) error {
	if req == nil {
		return errors.New("request is empty")
	}

	if strings.TrimSpace(req.RoomNumber) == "" {
		return errors.New("Room number is required")
	}

	if len(req.RoomNumber) > 50 {
		return errors.New("Room number is too long")
	}

	if strings.TrimSpace(req.RoomType) == "" {
		return errors.New("Room type is required")
	}

	switch req.RoomType {
	case "single", "double", "suite":
	default:
		return errors.New("Invalid room type (allowed: single, double, suite)")
	}

	if req.PricePerNight <= 0 {
		return errors.New("Price per night must be greater than 0")
	}

	if req.Capacity <= 0 {
		return errors.New("Capacity must be at least 1")
	}

	if len(req.Description) > 700 {
		return errors.New("Description is too long (max 500 characters)")
	}

	return nil
}

func ValidateRoomStatus(req *schemas.StatusUpdateRoom) error {
	if req == nil {
		return errors.New("request is empty")
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		return errors.New("status is required")
	}

	switch status {
	case "available", "booked", "maintenance":
		return nil
	default:
		return errors.New("invalid status (allowed: available, booked, maintenance)")
	}
}

func ValidateRoomUpdate(req *schemas.RoomSchemasUpdate) error {
	if req == nil {
		return errors.New("request is empty")
	}

	if req.RoomNumber != nil {
		roomNumber := strings.TrimSpace(*req.RoomNumber)
		if roomNumber == "" {
			return errors.New("Room number cannot be empty")
		}
		if len(roomNumber) > 50 {
			return errors.New("Room number is too long")
		}
	}

	if req.RoomType != nil {
		switch *req.RoomType {
		case "single", "double", "suite":
		default:
			return errors.New("Invalid room type (allowed: single, double, suite)")
		}
	}

	if req.PricePerNight != nil && *req.PricePerNight <= 0 {
		return errors.New("Price per night must be greater than 0")
	}

	if req.Capacity != nil && *req.Capacity <= 0 {
		return errors.New("Capacity must be at least 1")
	}

	if req.Description != nil && len(*req.Description) > 700 {
		return errors.New("Description is too long (max 700 characters)")
	}

	return nil
}
