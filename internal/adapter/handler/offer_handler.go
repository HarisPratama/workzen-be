package handler

import (
	"context"
	"strings"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

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
	status := strings.ToUpper(c.Query("status"))

	query := entity.OfferQueryString{
		Limit:  10,
		Page:   1,
		Search: c.Query("search"),
		Status: status,
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

	if c.Query("candidate_application_id") != "" {
		if val, err := conv.StringToInt64(c.Query("candidate_application_id")); err == nil {
			query.CandidateApplicationID = val
		}
	}

	results, totalData, totalPages, err := h.offerService.GetOffersByTenant(context.Background(), int64(tenantID), query)
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
			CandidateApplicationID: item.CandidateApplicationID,
			Position:               item.Position,
			Department:             item.Department,
			EmploymentType:         item.EmploymentType,
			BaseSalary:             item.BaseSalary,
			Currency:               item.Currency,
			Status:                 item.Status,
			StartDate:              item.StartDate.In(jakartaTZ).Format("2006-01-02"),
			ExpiryDate:             item.ExpiryDate.In(jakartaTZ).Format("2006-01-02"),
			SentAt:                 item.SentAt,
			RespondedAt:            item.RespondedAt,
		})
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         query.Page,
		PerPage:      query.Limit,
		TotalPages:   int(totalPages),
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
		ID:                    result.ID,
		Position:              result.Position,
		Department:            result.Department,
		EmploymentType:        result.EmploymentType,
		BaseSalary:            result.BaseSalary,
		Currency:              result.Currency,
		Bonus:                 result.Bonus,
		Benefits:              result.Benefits,
		ProbationPeriodMonths: result.ProbationPeriodMonths,
		NoticePeriodDays:      result.NoticePeriodDays,
		StartDate:             result.StartDate.In(jakartaTZ).Format("2006-01-02"),
		ExpiryDate:            result.ExpiryDate.In(jakartaTZ).Format("2006-01-02"),
		Status:                result.Status,
		SentAt:                result.SentAt,
		RespondedAt:           result.RespondedAt,
		Notes:                 result.Notes,
		Terms:                 result.Terms,
		Feedback:              result.Feedback,
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

	var req entity.OfferEntityRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateOffer - 2"
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

	if err := h.offerService.CreateOffer(context.Background(), req, int64(tenantID)); err != nil {
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

	var req entity.OfferUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateOffer - 3"
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
