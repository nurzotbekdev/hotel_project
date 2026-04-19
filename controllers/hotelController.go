package controllers

import (
	"net/http"
	"restaurant_manager/helper"
	"restaurant_manager/models"
	"restaurant_manager/schemas"
	"restaurant_manager/services"
	"restaurant_manager/validators"

	"github.com/gin-gonic/gin"
)

type HotelController struct {
	HotelService services.HotelService
}

func NewHotelController(hotel services.HotelService) *HotelController {
	return &HotelController{HotelService: hotel}
}

func (hotel *HotelController) Create(ctx *gin.Context) {
	var body schemas.HotelSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateCreateHotel(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newHotel := models.Hotel{
		Name:        body.Name,
		Address:     body.Address,
		Description: body.Description,
		Phone:       body.Phone,
		Email:       body.Email,
	}

	if err := hotel.HotelService.CreateHotel(newHotel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add hotel to database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hotel successfully created",
	})
}

func (hotel *HotelController) AllHotel(ctx *gin.Context) {
	hotels, err := hotel.HotelService.GetHotelAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch hotels",
		})
		return
	}

	response := make([]gin.H, 0)
	for _, row := range hotels {

		images := make([]gin.H, 0)
		for _, img := range row.HotelImages {
			images = append(images, gin.H{
				"id":         img.ID,
				"image":      img.HotelImage,
				"created_at": img.CreatedAt,
			})
		}

		response = append(response, gin.H{
			"id":          row.ID,
			"name":        row.Name,
			"address":     row.Address,
			"description": row.Description,
			"phone":       row.Phone,
			"email":       row.Email,
			"created_at":  row.CreatedAt,
			"images":      images,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (hotel *HotelController) HotelDetail(ctx *gin.Context) {
	hotelID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	row, err := hotel.HotelService.GetHotelByID(hotelID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Hotel not found",
		})
		return
	}

	images := make([]gin.H, 0)
	for _, img := range row.HotelImages {
		images = append(images, gin.H{
			"id":         img.ID,
			"image":      img.HotelImage,
			"created_at": img.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          row.ID,
		"name":        row.Name,
		"address":     row.Address,
		"description": row.Description,
		"phone":       row.Phone,
		"email":       row.Email,
		"created_at":  row.CreatedAt,
		"images":      images,
	})
}

func (hotel *HotelController) UpdateHotel(ctx *gin.Context) {
	hotelID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.EditHotel
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := validators.ValidateUpdateHotel(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := hotel.HotelService.SetHotel(hotelID, body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hotel data updated successfully",
	})
}

func (hotel *HotelController) DeleteHotel(ctx *gin.Context) {
	hotelID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := hotel.HotelService.DeleteHotel(hotelID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hotel data deleted successful",
	})
}
