package service

import (
	"context"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
	UpdatePassword(ctx context.Context, newPass string, id int64) error
}

type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	result, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		code := "[SERVICE] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (u *userService) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	password, err := conv.HashPassword(newPass)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	err = u.userRepo.UpdatePassword(ctx, password, id)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
