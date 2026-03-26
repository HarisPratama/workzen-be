package handler

import (
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CandidateHandler interface {
	GetCandidatesByTenant(c *fiber.Ctx) error
	GetCandidateDetailByTenant(c *fiber.Ctx) error
	CreateCandidate(c *fiber.Ctx) error
}

type candidateHandler struct {
	candidateService service.CandidateService
}

func (c2 *candidateHandler) GetCandidateDetailByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetCandidateDetailByTenant - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	id := c.Params("candidateID")
	candidateID, err := conv.StringToInt64(id)
	if err != nil {
		code := "[HANDLER] GetCandidateDetailByTenant - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid candidate ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := c2.candidateService.GetDetailCandidateByTenant(c.Context(), int64(tenantID), candidateID)
	if err != nil {
		code := "[HANDLER] GetCandidateDetailByTenant - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = result

	return c.JSON(defaultSuccessResponse)
}

func (c2 *candidateHandler) GetCandidatesByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetCandidatesByTenant - 1"
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
			code := "[HANDLER] GetCandidatesByTenant - 2"
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

	status := "active"
	if c.Query("status") != "" {
		status = c.Query("status")
	}

	reqEntity := entity.CandidateQueryString{
		Limit:     limit,
		Page:      page,
		OrderBy:   orderBy,
		Search:    search,
		OrderType: orderType,
		Status:    status,
	}

	tenantID := claims.TenantID
	results, totalData, totalPages, err := c2.candidateService.GetCandidatesByTenant(c.Context(), int64(tenantID), reqEntity)

	if err != nil {
		code := "[HANDLER] GetCandidatesByTenant - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respCandidates := []response.CandidateResponse{}
	for _, candidate := range results {
		respCandidate := response.CandidateResponse{
			ID:        candidate.ID,
			FullName:  candidate.FullName,
			Email:     candidate.Email,
			Phone:     candidate.Phone,
			Address:   candidate.Address,
			Source:    candidate.Source,
			Status:    candidate.Status,
			BirthDate: candidate.BirthDate.Local().Format("02 January 2006"),
		}

		respCandidates = append(respCandidates, respCandidate)
	}

	defaultSuccessResponse.Data = respCandidates
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func (c2 *candidateHandler) CreateCandidate(c *fiber.Ctx) error {
	var req request.CandidateRequest
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		code = "[Handler] CreateCandidate - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] CreateCandidate - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateCandidate - 3"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CandidateEntity{
		TenantID:  int64(tenantID),
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		BirthDate: req.BirthDate,
		Address:   req.Address,
		Source:    req.Source,
		Status:    req.Status,
		CitizenID: req.CitizenID,
	}

	err := c2.candidateService.CreateCandidate(c.Context(), reqEntity)

	if err != nil {
		code = "[Handler] CreateCandidate - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func NewCandidateHandler(candidateService service.CandidateService) CandidateHandler {
	return &candidateHandler{candidateService: candidateService}
}
