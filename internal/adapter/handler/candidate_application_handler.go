package handler

import (
	"strings"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CandidateApplicationHandler interface {
	GetCandidateApplicationByTenantMR(c *fiber.Ctx) error
	GetCandidateApplicationDetail(c *fiber.Ctx) error
	CreateCandidateApplication(c *fiber.Ctx) error
	UpdateCandidateApplication(c *fiber.Ctx) error
	DeleteCandidateApplication(c *fiber.Ctx) error
}

type candidateApplicationHandler struct {
	candidateApplicationService service.CandidateApplicationService
}

func (c2 *candidateApplicationHandler) CreateCandidateApplication(c *fiber.Ctx) error {
	var req request.CandidateApplicationRequest
	claims, ok := c.Locals("user").(*entity.JwtData)

	if !ok || claims.TenantID == 0 {
		code = "[HANDLER] CreateCandidateApplication - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] CreateCandidateApplication - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateCandidateApplication - 3"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	application := &entity.CandidateApplicationEntity{
		TenantID:          int64(tenantID),
		CandidateID:       req.CandidateID,
		ManpowerRequestID: req.ManpowerRequestID,
	}

	err := c2.candidateApplicationService.CreateCandidateApplication(c.Context(), application)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  false,
				"message": "Candidate already applied",
			})
		}

		code := "[HANDLER] CreateCandidateApplication - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Candidate applied successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (c2 *candidateApplicationHandler) GetCandidateApplicationDetail(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("applicationID")
	applicationID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid application ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := c2.candidateApplicationService.GetDetailCandidateApplication(c.Context(), int64(tenantID), applicationID)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	resp := response.CandidateApplicationResponse{
		ID:        result.ID,
		Status:    result.Status,
		AppliedAt: result.AppliedAt,
		Candidate: response.CandidateResponse{
			ID:       result.Candidate.ID,
			FullName: result.Candidate.FullName,
			Email:    result.Candidate.Email,
			Phone:    result.Candidate.Phone,
			Status:   result.Candidate.Status,
		},
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = resp

	return c.JSON(defaultSuccessResponse)
}

func (c2 *candidateApplicationHandler) UpdateCandidateApplication(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("applicationID")
	applicationID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid application ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.CandidateApplicationUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	application := &entity.CandidateApplicationEntity{
		Status: req.Status,
	}

	err = c2.candidateApplicationService.UpdateCandidateApplication(c.Context(), int64(tenantID), applicationID, application)
	if err != nil {
		code := "[HANDLER] UpdateCandidateApplication - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (c2 *candidateApplicationHandler) DeleteCandidateApplication(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("applicationID")
	applicationID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid application ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = c2.candidateApplicationService.DeleteCandidateApplication(c.Context(), int64(tenantID), applicationID)
	if err != nil {
		code := "[HANDLER] DeleteCandidateApplication - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Candidate application deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (c2 *candidateApplicationHandler) GetCandidateApplicationByTenantMR(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetCandidateApplicationByTenantMR - 1"
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
			code := "[HANDLER] GetCandidateApplicationByTenantMR - 2"
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

	reqEntity := entity.CandidateApplicationQueryString{
		Limit:     limit,
		Page:      page,
		OrderBy:   orderBy,
		Search:    search,
		OrderType: orderType,
		Status:    status,
	}

	tenantID := claims.TenantID
	idParam := c.Params("manpowerRequestID")
	manpowerRequestID, err := conv.StringToInt64(idParam)

	if tenantID == 0 {
		code := "[HANDLER] GetCandidateApplicationByTenantMR - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Tenant ID is empty"
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	results, totalData, totalPages, err := c2.candidateApplicationService.GetCandidateApplicationByTenantMR(c.Context(), int64(tenantID), manpowerRequestID, reqEntity)

	if err != nil {
		code := "[HANDLER] GetCandidateApplicationByTenantMR - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respCandidateApplications := []response.CandidateApplicationResponse{}
	for _, result := range results {
		respCandidateApplication := response.CandidateApplicationResponse{
			ID:        result.ID,
			Status:    result.Status,
			AppliedAt: result.AppliedAt,
			Candidate: response.CandidateResponse{
				ID:       result.Candidate.ID,
				FullName: result.Candidate.FullName,
				Email:    result.Candidate.Email,
				Status:   result.Candidate.Status,
			},
		}
		respCandidateApplications = append(respCandidateApplications, respCandidateApplication)
	}

	defaultSuccessResponse.Data = respCandidateApplications
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func NewCandidateApplicationHandler(candidateApplicationService service.CandidateApplicationService) CandidateApplicationHandler {
	return &candidateApplicationHandler{candidateApplicationService: candidateApplicationService}
}
