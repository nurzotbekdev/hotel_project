package controllers

import (
	"errors"
	"net/http"
	"restaurant_manager/config"
	"restaurant_manager/helper"
	"restaurant_manager/models"
	"restaurant_manager/schemas"
	"restaurant_manager/services"
	"restaurant_manager/validators"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(user services.UserService) *UserController {
	return &UserController{UserService: user}
}

func (user *UserController) Register(ctx *gin.Context) {
	var body schemas.UserRegister
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateRegisterUser(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newUser := models.User{
		Name:     body.Name,
		Surname:  body.Surname,
		Phone:    body.Phone,
		Email:    body.Email,
		Password: body.Password,
		Role:     "user",
	}

	if err := user.UserService.SignUp(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add user to database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User successfully created",
	})
}

func (user *UserController) Login(ctx *gin.Context) {
	var body schemas.UserLogin
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateLoginUser(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := user.UserService.SignIn(body.Email, body.Password)
	if err != nil {
		if errors.Is(err, services.EmailOrPasswordErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": services.EmailOrPasswordErr,
			})
			return
		}
		if errors.Is(err, services.TokenInvalid) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": services.TokenInvalid,
			})
			return
		}
		if errors.Is(err, services.RedisErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": services.RedisErr,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", accessToken, 1800, "/", "", false, true)
	ctx.SetCookie("RefreshToken", refreshToken, 1209600, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

func (user *UserController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("RefreshToken")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Refresh token not found",
		})
		return
	}

	userID, err := config.RedisClient.Get(config.Ctx, refreshToken).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token not found in Redis",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	newAccessToken, err := user.UserService.RefreshAccessToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie("Authorization", newAccessToken, 1800, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "New token created",
		"user_id": userID,
	})
}

func (user *UserController) MyProfile(ctx *gin.Context) {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	currentUser := userData.(models.User)

	ctx.JSON(http.StatusOK, gin.H{
		"id":         currentUser.ID,
		"name":       currentUser.Name,
		"surname":    currentUser.Surname,
		"phone":      currentUser.Phone,
		"email":      currentUser.Email,
		"role":       currentUser.Role,
		"created_at": currentUser.CreatedAt,
		"updated_at": currentUser.UpdatedAt,
	})
}

func (user *UserController) UpdateRole(ctx *gin.Context) {
	userID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body schemas.EditRole
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateUserRole(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := user.UserService.EditUserRole(userID, body.Role); err != nil {
		if errors.Is(err, services.UserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": services.UserNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role update successful",
	})
}

func (user *UserController) UpdateProfile(ctx *gin.Context) {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}
	currentUser := userData.(models.User)

	var body schemas.EditUserData
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := validators.ValidateUpdateUser(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := user.UserService.EditMyProfile(currentUser.ID, body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User data update successful",
	})
}

func (user *UserController) DeleteAccount(ctx *gin.Context) {
	userData, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}
	currentUser := userData.(models.User)

	if err := user.UserService.DeleteProfile(currentUser.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User data deleted successful",
	})
}
