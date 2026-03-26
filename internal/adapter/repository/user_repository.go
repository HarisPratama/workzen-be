package repository

import (
	"context"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
	UpdatePassword(ctx context.Context, newPass string, id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var modelUser model.User
	err := u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code := "[REPOSITORY] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:    modelUser.ID,
		Name:  modelUser.Name,
		Email: modelUser.Email,
	}, nil
}

func (u *userRepository) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	err = u.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPass).Error
	if err != nil {
		code := "[REPOSITORY] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
