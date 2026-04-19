package validators

import (
	"errors"
	"regexp"
	"restaurant_manager/config"
	"restaurant_manager/helper"
	"restaurant_manager/schemas"
	"strings"
)

var phoneRegex = regexp.MustCompile(`^(\+998)?[0-9]{9}$`)

func ValidateRegisterUser(req *schemas.UserRegister) error {
	if req == nil {
		return errors.New("request is empty")
	}

	switch {
	case strings.TrimSpace(req.Name) == "":
		return errors.New("Name is required")
	case len(req.Name) < 3:
		return errors.New("Name must be at least 3 characters")

	case strings.TrimSpace(req.Surname) == "":
		return errors.New("Surname is required")
	case len(req.Surname) < 3:
		return errors.New("Surname must be at least 3 characters")

	case strings.TrimSpace(req.Phone) == "":
		return errors.New("Phone is required")
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

	case strings.TrimSpace(req.Password) == "":
		return errors.New("Password is required")
	case len(req.Password) < 6:
		return errors.New("Password must be at least 6 characters")
	case !helper.HasNumber(req.Password):
		return errors.New("password must contain at least 1 number")
	case !helper.HasUpper(req.Password):
		return errors.New("password must contain at least 1 uppercase letter")
	}

	return nil
}

func ValidateLoginUser(req *schemas.UserLogin) error {
	switch {
	case strings.TrimSpace(req.Email) == "":
		return errors.New("Email is required")
	case !strings.Contains(req.Email, "@"):
		return errors.New("Invalid email format")
	case len(req.Email) < 4:
		return errors.New("Email must be at least 6 characters")

	case strings.TrimSpace(req.Password) == "":
		return errors.New("Password is required")
	case len(req.Password) < 6:
		return errors.New("Password must be at least 6 characters")
	case !helper.HasNumber(req.Password):
		return errors.New("password must contain at least 1 number")
	case !helper.HasUpper(req.Password):
		return errors.New("password must contain at least 1 uppercase letter")
	}

	return nil
}

func ValidateUserRole(req *schemas.EditRole) error {
	role := strings.TrimSpace(req.Role)

	if role == "" {
		return errors.New("Role is required")
	}

	if role != "admin" && role != "user" {
		return errors.New("Invalid role. Only 'admin' or 'user' allowed")
	}

	return nil
}

func ValidateUpdateUser(req *schemas.EditUserData) error {
	if req == nil {
		return errors.New("request is empty")
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 {
			return errors.New("Name must be at least 3 characters")
		}
	}

	if req.Surname != nil {
		surname := strings.TrimSpace(*req.Surname)
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

	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if password == "" {
			return errors.New("Password cannot be empty")
		}
		if len(password) < 6 {
			return errors.New("Password must be at least 6 characters")
		}
		if !helper.HasNumber(password) {
			return errors.New("Password must contain at least 1 number")
		}
		if !helper.HasUpper(password) {
			return errors.New("Password must contain at least 1 uppercase letter")
		}
	}

	return nil
}
