package services

import (
	"errors"
	"restaurant_manager/config"
	"restaurant_manager/models"
	"restaurant_manager/schemas"
	"restaurant_manager/security"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	SignUp(user models.User) error
	SignIn(email, password string) (string, string, error)
	RefreshAccessToken(refreshToken string) (string, error)
	EditUserRole(userID uint, role string) error
	EditMyProfile(userID uint, data schemas.EditUserData) error
	DeleteProfile(userID uint) error
}

type userService struct{}

func NewUserServices() UserService {
	return &userService{}
}

var (
	EmailOrPasswordErr = errors.New("Email or password error.")
	TokenInvalid       = errors.New("Token was not created")
	RedisErr           = errors.New("Refresh token not saved to Redis")
	UserNotFound       = errors.New("User not found")
)

func (s *userService) SignUp(user models.User) error {
	hash, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return config.DB.Create(&user).Error
}

func (s *userService) SignIn(email, password string) (string, string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", EmailOrPasswordErr
	}
	if !security.CheckPassword(user.Password, password) {
		return "", "", EmailOrPasswordErr
	}
	accessToken, err := security.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", TokenInvalid
	}
	refreshToken, err := security.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", TokenInvalid
	}
	err = config.RedisClient.Set(config.Ctx, refreshToken, user.ID, 14*24*time.Hour).Err()
	if err != nil {
		return "", "", RedisErr
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RefreshAccessToken(refreshToken string) (string, error) {
	userID, err := security.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", TokenInvalid
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return "", EmailOrPasswordErr
	}

	accessToken, err := security.GenerateAccessToken(user.ID)
	if err != nil {
		return "", TokenInvalid
	}

	return accessToken, nil
}

func (s *userService) EditUserRole(userID uint, role string) error {
	result := config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", role)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return UserNotFound
	}

	return nil
}

func (s *userService) EditMyProfile(userID uint, data schemas.EditUserData) error {
	updates := map[string]interface{}{}

	if data.Name != nil {
		updates["name"] = *data.Name
	}
	if data.Surname != nil {
		updates["surname"] = *data.Surname
	}
	if data.Phone != nil {
		updates["phone"] = *data.Phone
	}
	if data.Email != nil {
		updates["email"] = *data.Email
	}

	if data.Password != nil && *data.Password != "" {
		hash, err := security.HashPassword(*data.Password)
		if err != nil {
			return err
		}
		updates["password"] = hash
	}

	if len(updates) == 0 {
		return errors.New("No fields to update")
	}

	result := config.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *userService) DeleteProfile(userID uint) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}
	return config.DB.Unscoped().Delete(&user).Error
}
