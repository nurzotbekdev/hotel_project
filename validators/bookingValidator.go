package validators

import (
	"errors"
	"restaurant_manager/schemas"
	"strings"
)

func ValidateBookingStatus(req *schemas.UpdateBookingStatus) error {
	if req == nil {
		return errors.New("request is empty")
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		return errors.New("status is required")
	}

	switch status {
	case "pending", "approved", "cancelled":
		return nil
	default:
		return errors.New("invalid status value")
	}

}
