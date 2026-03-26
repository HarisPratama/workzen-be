package handler

import (
	"context"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ManpowerReqHandler interface {
	GetManpowerReqByTenant(c *fiber.Ctx) error
	GetDetailManpowerRequestByTenant(c *fiber.Ctx) error
	CreateManpowerReq(c *fiber.Ctx) error
}

type manpowerReqHandler struct {
	manpowerService service.ManpowerReqService
}

func (m *manpowerReqHandler) GetDetailManpowerRequestByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)

	if claims.UserID == 0 || claims.TenantID == 0 {
		code := "[HANDLER] GetDetailManpowerRequestByTenant - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	idParam := c.Params("manpowerRequestID")
	manpowerRequestID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetDetailManpowerRequestByTenant - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := m.manpowerService.GetDetailManpowerRequestByTenant(context.Background(), int64(tenantID), manpowerRequestID)
	if err != nil {
		code := "[HANDLER] GetDetailManpowerRequestByTenant - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respManpowerRequest := response.ManpowerReqResponse{
		ID:             result.ID,
		Position:       result.Position,
		RequiredCount:  result.RequiredCount,
		Hired:          result.Hired,
		SalaryMin:      result.SalaryMin,
		SalaryMax:      result.SalaryMax,
		WorkLocation:   result.WorkLocation,
		JobDescription: result.JobDescription,
		DeadlineDate:   result.DeadlineDate.Local().Format("02 January 2006"),
		CreatedAt:      result.CreatedAt.Local().Format("02 January 2006"),
		Status:         result.Status,
		Client: response.ClientResponse{
			ID:          result.Client.ID,
			CompanyName: result.Client.CompanyName,
			Address:     result.Client.Address,
		},
	}

	defaultSuccessResponse.Data = respManpowerRequest
	return c.JSON(&defaultSuccessResponse)
}

func (m *manpowerReqHandler) GetManpowerReqByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetManpowerReqByTenant - 1"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code := "[HANDLER] GetManpowerReqByTenant - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	status := ""
	if c.Query("status") != "" {
		status = c.Query("status")
	}

	reqEntity := entity.ManpowerReqQueryString{
		Limit:     limit,
		Page:      page,
		OrderBy:   orderBy,
		Search:    search,
		OrderType: orderType,
		Status:    status,
	}

	tenantID := claims.TenantID

	if tenantID == 0 {
		code := "[HANDLER] GetManpowerReqByTenant - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Tenant ID is empty"

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	results, totalData, totalPages, err := m.manpowerService.GetManpowerReqByTenant(c.Context(), int64(tenantID), reqEntity)

	if err != nil {
		code := "[HANDLER] GetManpowerReqByTenant - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respManpowerReqs := []response.ManpowerReqResponse{}
	for _, manpowerReq := range results {
		respManpowerReq := response.ManpowerReqResponse{
			ID:             manpowerReq.ID,
			Position:       manpowerReq.Position,
			RequiredCount:  manpowerReq.RequiredCount,
			SalaryMin:      manpowerReq.SalaryMin,
			SalaryMax:      manpowerReq.SalaryMax,
			WorkLocation:   manpowerReq.WorkLocation,
			JobDescription: manpowerReq.JobDescription,
			Status:         manpowerReq.Status,
			DeadlineDate:   manpowerReq.DeadlineDate.Local().Format("02 January 2006"),
			CreatedAt:      manpowerReq.DeadlineDate.Local().Format("02 January 2006"),
			Client: response.ClientResponse{
				ID:          manpowerReq.ClientID,
				CompanyName: manpowerReq.Client.CompanyName,
			},
		}

		respManpowerReqs = append(respManpowerReqs, respManpowerReq)
	}

	defaultSuccessResponse.Data = respManpowerReqs
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func (m *manpowerReqHandler) CreateManpowerReq(c *fiber.Ctx) error {
	var req request.ManPowerRequest
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		code = "[HANDLER] CreateManpowerReq - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] CreateManpowerReq - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateManpowerReq - 3"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.ManpowerReqEntity{
		TenantID:       int64(tenantID),
		ClientID:       req.ClientID,
		Position:       req.Position,
		RequiredCount:  req.RequiredCount,
		SalaryMin:      req.SalaryMin,
		SalaryMax:      req.SalaryMax,
		WorkLocation:   req.WorkLocation,
		JobDescription: req.JobDescription,
		DeadlineDate:   req.DeadlineDate,
	}

	err := m.manpowerService.CreateManpowerReq(c.Context(), reqEntity)

	if err != nil {
		code := "[HANDLER] CreateManpowerReq - 4"
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

func NewManpowerReqHandler(manpowerReqService service.ManpowerReqService) ManpowerReqHandler {
	return &manpowerReqHandler{manpowerService: manpowerReqService}
}
