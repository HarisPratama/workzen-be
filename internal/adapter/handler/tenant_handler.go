package handler

import (
	"bwanews/internal/adapter/handler/request"
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/service"
	validatorLib "bwanews/lib/validator"

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
		errorResp.Meta.Message = "Invalid request body"

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
		Plan:        "TRIAL",
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

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func NewTenantHandler(tenantService service.TenantService) TenantHandler {
	return &tenantHandler{tenantService: tenantService}
}
