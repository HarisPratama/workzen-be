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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type InterviewHandler interface {
	GetInterviews(c *fiber.Ctx) error
	GetInterviewByID(c *fiber.Ctx) error
	CreateInterview(c *fiber.Ctx) error
	UpdateInterview(c *fiber.Ctx) error
	DeleteInterview(c *fiber.Ctx) error
}

type interviewHandler struct {
	interviewService service.InterviewService
}

func (h *interviewHandler) GetInterviews(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetInterviews - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	query := entity.InterviewQueryString{
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

	results, totalData, totalPages, err := h.interviewService.GetInterviewsByTenant(context.Background(), tenantID, query)
	if err != nil {
		code := "[HANDLER] GetInterviews - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	var respData []response.InterviewResponse
	for _, item := range results {
		respData = append(respData, response.InterviewResponse{
			ID:              item.ID,
			InterviewType:   item.InterviewType,
			ScheduledAt:     item.ScheduledAt.Local().Format("02 January 2006 15:04"),
			DurationMinutes: item.DurationMinutes,
			Location:        item.Location,
			MeetingLink:     item.MeetingLink,
			Status:          item.Status,
			Feedback:        item.Feedback,
			Rating:          item.Rating,
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

func (h *interviewHandler) GetInterviewByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetInterviewByID - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("interviewID")
	interviewID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetInterviewByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid interview ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.interviewService.GetInterviewByID(context.Background(), interviewID)
	if err != nil {
		code := "[HANDLER] GetInterviewByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	respData := response.InterviewResponse{
		ID:              result.ID,
		InterviewType:   result.InterviewType,
		ScheduledAt:     result.ScheduledAt.Local().Format("02 January 2006 15:04"),
		DurationMinutes: result.DurationMinutes,
		Location:        result.Location,
		MeetingLink:     result.MeetingLink,
		Status:          result.Status,
		Feedback:        result.Feedback,
		Rating:          result.Rating,
		CompletedAt:     result.CompletedAt.Local().Format("02 January 2006 15:04"),
		CancelledAt:     result.CancelledAt.Local().Format("02 January 2006 15:04"),
		CancelReason:    result.CancelReason,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData

	return c.JSON(defaultSuccessResponse)
}

func (h *interviewHandler) CreateInterview(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] CreateInterview - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	var req request.InterviewRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateInterview - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[HANDLER] CreateInterview - 3"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	// Parse scheduled time
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		code := "[HANDLER] CreateInterview - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid scheduled time format"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.InterviewEntityRequest{
		CandidateID:       req.CandidateID,
		ManpowerRequestID: req.ManpowerRequestID,
		InterviewType:     req.InterviewType,
		ScheduledAt:       req.ScheduledAt,
		DurationMinutes:   req.DurationMinutes,
		Location:          req.Location,
		MeetingLink:       req.MeetingLink,
	}

	if err := h.interviewService.CreateInterview(context.Background(), reqEntity, tenantID); err != nil {
		code := "[HANDLER] CreateInterview - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Interview created successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (h *interviewHandler) UpdateInterview(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] UpdateInterview - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("interviewID")
	interviewID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] UpdateInterview - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid interview ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.InterviewUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateInterview - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[HANDLER] UpdateInterview - 4"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.InterviewUpdateRequest{
		Status:    req.Status,
		Feedback:  req.Feedback,
		Rating:    req.Rating,
		Location:  req.Location,
		ScheduledAt: req.ScheduledAt,
	}

	if err := h.interviewService.UpdateInterview(context.Background(), interviewID, reqEntity); err != nil {
		code := "[HANDLER] UpdateInterview - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Interview updated successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *interviewHandler) DeleteInterview(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] DeleteInterview - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("interviewID")
	interviewID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] DeleteInterview - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid interview ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.interviewService.DeleteInterview(context.Background(), interviewID); err != nil {
		code := "[HANDLER] DeleteInterview - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Interview deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func NewInterviewHandler(interviewService service.InterviewService) InterviewHandler {
	return &interviewHandler{interviewService: interviewService}
}