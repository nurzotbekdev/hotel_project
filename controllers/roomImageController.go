package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"restaurant_manager/helper"
	"restaurant_manager/models"
	"restaurant_manager/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomImageController struct {
	RoomImageService services.RoomImageService
}

func NewRoomImageController(image services.RoomImageService) *RoomImageController {
	return &RoomImageController{RoomImageService: image}
}

func (image *RoomImageController) Create(ctx *gin.Context) {
	roomIDStr := ctx.PostForm("room_id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	file, err := ctx.FormFile("image_url")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Image required",
		})
		return
	}
	os.MkdirAll("uploads/rooms/", os.ModePerm)
	ext := filepath.Ext(file.Filename)
	now := time.Now()
	filename := fmt.Sprintf("%s%s", now.Format("20060102150405"), ext)
	uploadPath := "uploads/rooms/" + filename

	if err := ctx.SaveUploadedFile(file, uploadPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	newRoomImage := models.RoomImage{
		RoomID:   uint(roomID),
		ImageURL: uploadPath,
	}

	if err := image.RoomImageService.CreateRoomImage(newRoomImage); err != nil {
		os.Remove(uploadPath)

		if errors.Is(err, services.RoomNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.RoomNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add room image to database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Room image successfully created",
	})
}

func (image *RoomImageController) AllRoomImage(ctx *gin.Context) {
	roomImages, err := image.RoomImageService.GetAllRoomImage()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if len(roomImages) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Room image not found",
		})
		return
	}

	response := make([]gin.H, 0)
	for _, row := range roomImages {
		response = append(response, gin.H{
			"id":         row.ID,
			"image":      row.ImageURL,
			"created_at": row.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (image *RoomImageController) RemoveRoomImage(ctx *gin.Context) {
	roomImageID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := image.RoomImageService.DeleteRoomImage(roomImageID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.RoomImageNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Room image data deleted successful",
	})
}
