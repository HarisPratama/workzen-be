package handler

import (
	"errors"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHandler interface {
	UpdatePassword(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] UpdatePassword - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if req.ConfirmPassword != req.NewPassword {
		code := "[HANDLER] UpdatePassword - 4"
		err = errors.New("password not match")
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = u.userService.UpdatePassword(c.Context(), req.NewPassword, int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] UpdatePassword - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (u *userHandler) GetUserByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetUserByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserByID(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] GetUserByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	defaultSuccessResponse.Data = resp

	return c.JSON(defaultSuccessResponse)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}
