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

type HotelImageController struct {
	HotelImageService services.HotelImageService
}

func NewHotelImageController(hotelImage services.HotelImageService) *HotelImageController {
	return &HotelImageController{HotelImageService: hotelImage}
}

func (hotelImage *HotelImageController) Create(ctx *gin.Context) {
	hotelIDStr := ctx.PostForm("hotel_id")
	hotelID, err := strconv.ParseUint(hotelIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid hotel ID",
		})
		return
	}

	file, err := ctx.FormFile("hotel_image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Image required",
		})
		return
	}
	os.MkdirAll("uploads/hotels/", os.ModePerm)
	ext := filepath.Ext(file.Filename)
	now := time.Now()
	filename := fmt.Sprintf("%s%s", now.Format("20060102150405"), ext)
	uploadPath := "uploads/hotels/" + filename

	if err := ctx.SaveUploadedFile(file, uploadPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	newHotelImage := models.HotelImage{
		HotelID:    uint(hotelID),
		HotelImage: uploadPath,
	}

	if err := hotelImage.HotelImageService.CreateHotelImage(newHotelImage); err != nil {
		os.Remove(uploadPath)

		if errors.Is(err, services.HotelNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.HotelNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add hotel image to database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hotel image successfully created",
	})
}

func (hotelImage *HotelImageController) AllHotelImage(ctx *gin.Context) {
	hotelImages, err := hotelImage.HotelImageService.GetAllHotelImage()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if len(hotelImages) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Hotel image not found",
		})
		return
	}

	response := make([]gin.H, 0)
	for _, row := range hotelImages {
		response = append(response, gin.H{
			"id":         row.ID,
			"image":      row.HotelImage,
			"created_at": row.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (hotelImage *HotelImageController) RemoveHotelImage(ctx *gin.Context) {
	hotelImageID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := hotelImage.HotelImageService.DeleteHotelImage(hotelImageID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.HotelImageNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hotel image data deleted successful",
	})
}
