package controllers

import (
	"errors"
	"net/http"
	"restaurant_manager/helper"
	"restaurant_manager/models"
	"restaurant_manager/schemas"
	"restaurant_manager/services"
	"restaurant_manager/validators"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingController struct {
	BookingService services.BookingService
}

func NewBookingController(booking services.BookingService) *BookingController {
	return &BookingController{BookingService: booking}
}

func (booking *BookingController) Create(ctx *gin.Context) {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}
	currentUser := userData.(models.User)

	var body schemas.BookingSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	newBooking := models.Booking{
		UserID:   currentUser.ID,
		RoomID:   body.RoomID,
		CheckIn:  body.CheckIn,
		CheckOut: body.CheckOut,
	}

	if err := booking.BookingService.CreateBooking(newBooking); err != nil {

		if errors.Is(err, services.RoomNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.RoomNotAvailable) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add booking to database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Booking successfully created",
	})
}

func (booking *BookingController) MyBookings(ctx *gin.Context) {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}
	currentUser := userData.(models.User)

	bookings, err := booking.BookingService.GetBookingMy(currentUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bookings",
		})
		return
	}

	response := make([]gin.H, 0)
	for _, row := range bookings {
		images := make([]gin.H, 0)
		for _, img := range row.Room.RoomImages {
			images = append(images, gin.H{
				"id":         img.ID,
				"image":      img.ImageURL,
				"created_at": img.CreatedAt,
			})
		}

		response = append(response, gin.H{
			"id":          row.ID,
			"check_in":    row.CheckIn,
			"check_out":   row.CheckOut,
			"total_price": row.TotalPrice,
			"status":      row.Status,
			"room": gin.H{
				"id":          row.Room.ID,
				"room_number": row.Room.RoomNumber,
				"room_type":   row.Room.RoomType,
			},
			"hotel": gin.H{
				"id":   row.Room.Hotel.ID,
				"name": row.Room.Hotel.Name,
			},
			"images": images,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})

}

func (booking *BookingController) MyBookingDetail(ctx *gin.Context) {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}
	currentUser := userData.(models.User)

	bookingID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	row, err := booking.BookingService.GetBookingByID(currentUser.ID, bookingID)
	if err != nil {
		if errors.Is(err, services.RoomNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.RoomNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rooms",
		})
		return
	}

	images := make([]gin.H, 0)
	for _, img := range row.Room.RoomImages {
		images = append(images, gin.H{
			"id":         img.ID,
			"image":      img.ImageURL,
			"created_at": img.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          row.ID,
		"check_in":    row.CheckIn,
		"check_out":   row.CheckOut,
		"total_price": row.TotalPrice,
		"status":      row.Status,
		"room": gin.H{
			"id":          row.Room.ID,
			"room_number": row.Room.RoomNumber,
			"room_type":   row.Room.RoomType,
		},
		"hotel": gin.H{
			"id":   row.Room.Hotel.ID,
			"name": row.Room.Hotel.Name,
		},
		"images": images,
	})
}

func (booking *BookingController) AdminBookingList(ctx *gin.Context) {
	bookings, err := booking.BookingService.GetAllBooking()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch bookings",
		})
		return
	}

	response := make([]gin.H, 0)
	for _, row := range bookings {
		images := make([]gin.H, 0)
		for _, img := range row.Room.RoomImages {
			images = append(images, gin.H{
				"id":         img.ID,
				"image":      img.ImageURL,
				"created_at": img.CreatedAt,
			})
		}

		response = append(response, gin.H{
			"id":          row.ID,
			"check_in":    row.CheckIn,
			"check_out":   row.CheckOut,
			"total_price": row.TotalPrice,
			"status":      row.Status,
			"room": gin.H{
				"id":          row.Room.ID,
				"room_number": row.Room.RoomNumber,
				"room_type":   row.Room.RoomType,
			},
			"hotel": gin.H{
				"id":   row.Room.Hotel.ID,
				"name": row.Room.Hotel.Name,
			},
			"user": gin.H{
				"id":      row.User.ID,
				"name":    row.User.Name,
				"surname": row.User.Surname,
				"phone":   row.User.Phone,
			},
			"images": images,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (booking *BookingController) AdminBookingByID(ctx *gin.Context) {
	bookingID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	row, err := booking.BookingService.GetAdminBookingByID(bookingID)
	if err != nil {
		if errors.Is(err, services.RoomNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.RoomNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rooms",
		})
		return
	}

	images := make([]gin.H, 0)
	for _, img := range row.Room.RoomImages {
		images = append(images, gin.H{
			"id":         img.ID,
			"image":      img.ImageURL,
			"created_at": img.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          row.ID,
		"check_in":    row.CheckIn,
		"check_out":   row.CheckOut,
		"total_price": row.TotalPrice,
		"status":      row.Status,
		"room": gin.H{
			"id":          row.Room.ID,
			"room_number": row.Room.RoomNumber,
			"room_type":   row.Room.RoomType,
		},
		"hotel": gin.H{
			"id":   row.Room.Hotel.ID,
			"name": row.Room.Hotel.Name,
		},
		"user": gin.H{
			"id":      row.User.ID,
			"name":    row.User.Name,
			"surname": row.User.Surname,
			"phone":   row.User.Phone,
		},
		"images": images,
	})
}

func (booking *BookingController) UpdateStatus(ctx *gin.Context) {
	bookingID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body schemas.UpdateBookingStatus
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateBookingStatus(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := booking.BookingService.EditStatusBooking(bookingID, body.Status); err != nil {
		if errors.Is(err, services.BookingNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.BookingNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Booking data updated successfully",
	})
}

func (booking *BookingController) RemoveBooking(ctx *gin.Context) {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}
	currentUser := userData.(models.User)

	bookingID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := booking.BookingService.DeleteBooking(currentUser.ID, bookingID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.BookingNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Booking data deleted successful",
	})
}
