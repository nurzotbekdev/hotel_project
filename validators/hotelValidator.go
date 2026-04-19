package validators

import (
	"errors"
	"regexp"
	"restaurant_manager/config"
	"restaurant_manager/helper"
	"restaurant_manager/schemas"
	"strings"
)

func ValidateCreateHotel(req *schemas.HotelSchemas) error {
	if req == nil {
		return errors.New("Request is empty.")
	}

	switch {
	case strings.TrimSpace(req.Name) == "":
		return errors.New("Hotel name is required")
	case len(req.Name) < 3:
		return errors.New("Hotel name must be at least 3 characters")

	case strings.TrimSpace(req.Address) == "":
		return errors.New("Hotel address is required")
	case len(req.Address) < 4:
		return errors.New("Hotel address must be at least 4 characters")

	case strings.TrimSpace(req.Phone) == "":
		return errors.New("Hotel phone is required")
	case !regexp.MustCompile(`^(\+998)?[0-9]{9}$`).MatchString(req.Phone):
		return errors.New("Invalid phone number format")

	case strings.TrimSpace(req.Email) == "":
		return errors.New("Email is required")
	case !strings.Contains(req.Email, "@"):
		return errors.New("Invalid email format")
	case len(req.Email) < 6:
		return errors.New("Email must be at least 6 characters")
	case !helper.IsEmailUnique(config.DB, req.Email):
		return errors.New("Email already exists")
	}

	return nil
}

func ValidateUpdateHotel(req *schemas.EditHotel) error {
	if req == nil {
		return errors.New("request is empty")
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 {
			return errors.New("Name must be at least 3 characters")
		}
	}

	if req.Address != nil {
		surname := strings.TrimSpace(*req.Address)
		if len(surname) < 3 {
			return errors.New("Surname must be at least 3 characters")
		}
	}

	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		if !phoneRegex.MatchString(phone) {
			return errors.New("Invalid phone number format")
		}
	}

	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if len(email) < 6 || !strings.Contains(email, "@") {
			return errors.New("Invalid email format")
		}
	}

	return nil
}
