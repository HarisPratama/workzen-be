package handler

import (
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OverviewHandler interface {
	GetOverview(c *fiber.Ctx) error
}

type overviewHandler struct {
	overviewService service.OverviewService
}

func (h *overviewHandler) GetOverview(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := int64(claims.TenantID)

	result, err := h.overviewService.GetOverviewByTenant(c.Context(), tenantID)
	if err != nil {
		code := "[HANDLER] GetOverview - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta = response.Meta{Status: true, Message: "Success"}
	defaultSuccessResponse.Data = result
	defaultSuccessResponse.Pagination = nil

	return c.JSON(defaultSuccessResponse)
}

func NewOverviewHandler(overviewService service.OverviewService) OverviewHandler {
	return &overviewHandler{overviewService: overviewService}
}
