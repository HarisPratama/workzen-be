package handler

import (
	"errors"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type TenantHandler interface {
	RegisterTenant(c *fiber.Ctx) error
}

type tenantHandler struct {
	tenantService service.TenantService
}

func (t *tenantHandler) RegisterTenant(c *fiber.Ctx) error {
	var req request.RegisterTenantRequest

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] RegisterTenant - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body, please check your input format"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] RegisterTenant - 2"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.RegisterTenantEntity{
		CompanyName: req.CompanyName,
		Address:     req.CompanyAddress,
		Plan:        "FREE",
		Name:        req.AdminName,
		Email:       req.AdminEmail,
		Password:    req.Password,
	}

	err = t.tenantService.RegisterTenant(c.Context(), reqEntity)

	if err != nil {
		code := "[HANDLER] RegisterTenant - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(errorResp)
		default:
			errorResp.Meta.Message = "Registration failed, please try again later"
			return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
		}
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Registration successful"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func NewTenantHandler(tenantService service.TenantService) TenantHandler {
	return &tenantHandler{tenantService: tenantService}
}
