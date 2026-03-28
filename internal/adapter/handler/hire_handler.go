package handler

import (
	"errors"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type HireHandler interface {
	HireCandidate(c *fiber.Ctx) error
}

type hireHandler struct {
	hireService service.HireService
}

func (h *hireHandler) HireCandidate(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] HireCandidate - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	idParam := c.Params("applicationID")
	applicationID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] HireCandidate - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid application ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req service.HireRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] HireCandidate - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.hireService.HireCandidate(c.Context(), int64(tenantID), applicationID, req); err != nil {
		code := "[HANDLER] HireCandidate - 4"
		log.Errorw(code, err)

		if errors.Is(err, service.ErrApplicationNotFound) {
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()
			return c.Status(fiber.StatusNotFound).JSON(errorResp)
		}

		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Candidate hired successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func NewHireHandler(hireService service.HireService) HireHandler {
	return &hireHandler{hireService: hireService}
}
