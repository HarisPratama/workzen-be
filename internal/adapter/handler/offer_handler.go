package handler

import (
	"bwanews/internal/adapter/handler/request"
	"bwanews/internal/adapter/handler/response"
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/service"
	"bwanews/lib/conv"
	validatorLib "bwanews/lib/validator"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OfferHandler interface {
	GetOffers(c *fiber.Ctx) error
	GetOfferByID(c *fiber.Ctx) error
	CreateOffer(c *fiber.Ctx) error
	UpdateOffer(c *fiber.Ctx) error
	DeleteOffer(c *fiber.Ctx) error
}

type offerHandler struct {
	offerService service.OfferService
}

func (h *offerHandler) GetOffers(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetOffers - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	query := entity.OfferQueryString{
		Limit:  10,
		Page:   1,
		Search: c.Query("search"),
		Status: c.Query("status"),
	}

	if c.Query("limit") != "" {
		if val, err := conv.StringToInt(c.Query("limit")); err == nil {
			query.Limit = val
		}
	}

	if c.Query("page") != "" {
		if val, err := conv.StringToInt(c.Query("page")); err == nil {
			query.Page = val
		}
	}

	if c.Query("order_by") != "" {
		query.OrderBy = c.Query("order_by")
	}

	if c.Query("order_type") != "" {
		query.OrderType = c.Query("order_type")
	}

	results, totalData, totalPages, err := h.offerService.GetOffersByTenant(context.Background(), tenantID, query)
	if err != nil {
		code := "[HANDLER] GetOffers - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	var respData []response.OfferResponse
	for _, item := range results {
		respData = append(respData, response.OfferResponse{
			ID:                     item.ID,
			Position:               item.Position,
			Department:             item.Department,
			EmploymentType:         item.EmploymentType,
			BaseSalary:             item.BaseSalary,
			Currency:               item.Currency,
			Status:                 item.Status,
			StartDate:              item.StartDate,
			ExpiryDate:             item.ExpiryDate,
			SentAt:                 item.SentAt,
			RespondedAt:            item.RespondedAt,
		})
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: totalData,
		Page:         query.Page,
		PerPage:      query.Limit,
		TotalPages:   totalPages,
	}

	return c.JSON(defaultSuccessResponse)
}

func (h *offerHandler) GetOfferByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetOfferByID - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("offerID")
	offerID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetOfferByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid offer ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.offerService.GetOfferByID(context.Background(), offerID)
	if err != nil {
		code := "[HANDLER] GetOfferByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	respData := response.OfferResponse{
		ID:                     result.ID,
		Position:               result.Position,
		Department:             result.Department,
		EmploymentType:         result.EmploymentType,
		BaseSalary:             result.BaseSalary,
		Currency:               result.Currency,
		Bonus:                  result.Bonus,
		Benefits:               result.Benefits,
		ProbationPeriodMonths:  result.ProbationPeriodMonths,
		NoticePeriodDays:       result.NoticePeriodDays,
		StartDate:              result.StartDate,
		ExpiryDate:             result.ExpiryDate,
		Status:                 result.Status,
		SentAt:                 result.SentAt,
		RespondedAt:          result.RespondedAt,
		Notes:                  result.Notes,
		Terms:                  result.Terms,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData

	return c.JSON(defaultSuccessResponse)
}

func (h *offerHandler) CreateOffer(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] CreateOffer - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	var req request.OfferRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateOffer - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[HANDLER] CreateOffer - 3"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.offerService.CreateOffer(context.Background(), req, tenantID); err != nil {
		code := "[HANDLER] CreateOffer - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Offer created successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (h *offerHandler) UpdateOffer(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] UpdateOffer - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("offerID")
	offerID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] UpdateOffer - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid offer ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.OfferUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateOffer - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[HANDLER] UpdateOffer - 4"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.offerService.UpdateOffer(context.Background(), offerID, req); err != nil {
		code := "[HANDLER] UpdateOffer - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Offer updated successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *offerHandler) DeleteOffer(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] DeleteOffer - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("offerID")
	offerID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] DeleteOffer - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid offer ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.offerService.DeleteOffer(context.Background(), offerID); err != nil {
		code := "[HANDLER] DeleteOffer - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Offer deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func NewOfferHandler(offerService service.OfferService) OfferHandler {
	return &offerHandler{offerService: offerService}
}