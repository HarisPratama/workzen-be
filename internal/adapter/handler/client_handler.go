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

type ClientHandler interface {
	GetClientByTenant(c *fiber.Ctx) error
	GetClientDetailByTenant(c *fiber.Ctx) error
	CreateClient(c *fiber.Ctx) error
	UpdateClient(c *fiber.Ctx) error
	DeleteClient(c *fiber.Ctx) error
}

type clientHandler struct {
	clientService service.ClientService
}

func (c2 *clientHandler) GetClientByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetClientByTenant - 1"
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
			code := "[HANDLER] GetClientByTenant - 2"
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

	reqEntity := entity.ClientQueryString{
		Limit:     limit,
		Page:      page,
		OrderBy:   orderBy,
		Search:    search,
		OrderType: orderType,
		Status:    status,
	}

	tenantID := claims.TenantID
	results, totalData, totalPages, err := c2.clientService.GetClientByTenant(c.Context(), reqEntity, int64(tenantID))

	if err != nil {
		code := "[HANDLER] GetClientByTenant - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respClients := []response.ClientResponse{}
	for _, client := range results {
		respClient := response.ClientResponse{
			ID:          client.ID,
			CompanyName: client.CompanyName,
			Address:     client.Address,
			CreatedAt:   client.CreatedAt.Local().Format("02 January 2006"),
		}
		respClients = append(respClients, respClient)
	}

	defaultSuccessResponse.Data = respClients
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func (c2 *clientHandler) GetClientDetailByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("clientID")
	clientID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid client ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := c2.clientService.GetDetailClientByTenant(c.Context(), int64(tenantID), clientID)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = result

	return c.JSON(defaultSuccessResponse)
}

func (c2 *clientHandler) CreateClient(c *fiber.Ctx) error {
	var req request.ClientRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	tenantID := claims.TenantID

	if userID == 0 || tenantID == 0 {
		code = "[Handler] CreateClient - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] CreateClient - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateClient - 3"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	client := &entity.ClientEntity{
		CompanyName: req.CompanyName,
		Address:     req.Address,
		TenantID:    int64(tenantID),
	}

	err := c2.clientService.CreateClient(c.Context(), client)
	if err != nil {
		code := "[HANDLER] CreateClient - 4"
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

func (c2 *clientHandler) UpdateClient(c *fiber.Ctx) error {
	var req request.ClientRequest
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("clientID")
	clientID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid client ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

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

	client := &entity.ClientEntity{
		CompanyName: req.CompanyName,
		Address:     req.Address,
	}

	err = c2.clientService.UpdateClient(c.Context(), int64(tenantID), clientID, client)
	if err != nil {
		code := "[HANDLER] UpdateClient - 1"
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

func (c2 *clientHandler) DeleteClient(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("clientID")
	clientID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid client ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = c2.clientService.DeleteClient(c.Context(), int64(tenantID), clientID)
	if err != nil {
		code := "[HANDLER] DeleteClient - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Client deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func NewClientHandler(clientService service.ClientService) ClientHandler {
	return &clientHandler{clientService: clientService}
}
