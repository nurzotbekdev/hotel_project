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

type RoomController struct {
	RoomService services.RoomService
}

func NewRoomController(room services.RoomService) *RoomController {
	return &RoomController{RoomService: room}
}

func (room *RoomController) Create(ctx *gin.Context) {
	var body schemas.RoomSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateRoom(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newRoom := models.Room{
		HotelID:       body.HotelID,
		RoomNumber:    body.RoomNumber,
		RoomType:      body.RoomType,
		PricePerNight: body.PricePerNight,
		Capacity:      body.Capacity,
		Description:   body.Description,
		Status:        "available",
	}

	if err := room.RoomService.CreateRoom(newRoom); err != nil {
		if errors.Is(err, services.HotelNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.HotelNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add room to database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Room successfully created",
	})
}

func (room *RoomController) AllRoom(ctx *gin.Context) {
	rooms, err := room.RoomService.GetRoomAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rooms",
		})
		return
	}

	response := make([]gin.H, 0)
	for _, row := range rooms {
		images := make([]gin.H, 0)
		for _, img := range row.RoomImages {
			images = append(images, gin.H{
				"id":         img.ID,
				"image":      img.ImageURL,
				"created_at": img.CreatedAt,
			})
		}

		response = append(response, gin.H{
			"id":              row.ID,
			"room_number":     row.RoomNumber,
			"room_type":       row.RoomType,
			"price_per_night": row.PricePerNight,
			"capacity":        row.Capacity,
			"description":     row.Description,
			"status":          row.Status,
			"created_at":      row.CreatedAt,
			"hotel": gin.H{
				"id":   row.Hotel.ID,
				"name": row.Hotel.Name,
			},
			"images": images,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})

}

func (room *RoomController) RoomDetail(ctx *gin.Context) {
	roomID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	row, err := room.RoomService.GetRoomByID(roomID)
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
	for _, img := range row.RoomImages {
		images = append(images, gin.H{
			"id":         img.ID,
			"image":      img.ImageURL,
			"created_at": img.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":              row.ID,
		"room_number":     row.RoomNumber,
		"room_type":       row.RoomType,
		"price_per_night": row.PricePerNight,
		"capacity":        row.Capacity,
		"description":     row.Description,
		"status":          row.Status,
		"created_at":      row.CreatedAt,
		"hotel": gin.H{
			"id":   row.Hotel.ID,
			"name": row.Hotel.Name,
		},
		"images": images,
	})
}

func (room *RoomController) AdminRoomList(ctx *gin.Context) {
	rooms, err := room.RoomService.AdminGetRoomAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rooms",
		})
		return
	}

	response := make([]gin.H, 0)
	for _, row := range rooms {
		images := make([]gin.H, 0)
		for _, img := range row.RoomImages {
			images = append(images, gin.H{
				"id":         img.ID,
				"image":      img.ImageURL,
				"created_at": img.CreatedAt,
			})
		}

		response = append(response, gin.H{
			"id":              row.ID,
			"room_number":     row.RoomNumber,
			"room_type":       row.RoomType,
			"price_per_night": row.PricePerNight,
			"capacity":        row.Capacity,
			"description":     row.Description,
			"status":          row.Status,
			"created_at":      row.CreatedAt,
			"hotel": gin.H{
				"id":   row.Hotel.ID,
				"name": row.Hotel.Name,
			},
			"images": images,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (room *RoomController) UpdateRoom(ctx *gin.Context) {
	roomID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.RoomSchemasUpdate
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.ValidateRoomUpdate(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := room.RoomService.EditRoom(roomID, body); err != nil {
		switch {
		case errors.Is(err, services.RoomNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": services.RoomNotFound.Error()})
		case errors.Is(err, services.HotelNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": services.HotelNotFound.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Room data updated successfully",
	})
}

func (room *RoomController) UpdateStatus(ctx *gin.Context) {
	roomID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body schemas.StatusUpdateRoom
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateRoomStatus(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := room.RoomService.EditStatusRoom(roomID, body.Status); err != nil {
		if errors.Is(err, services.RoomNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.RoomNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Room data updated successfully",
	})
}

func (room *RoomController) RemoveRoom(ctx *gin.Context) {
	roomID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := room.RoomService.DeleteRoom(roomID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.RoomNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Room data deleted successful",
	})
}
