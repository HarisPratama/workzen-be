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

type EmployeeHandler interface {
	GetEmployees(c *fiber.Ctx) error
	GetEmployeeDetail(c *fiber.Ctx) error
	CreateEmployee(c *fiber.Ctx) error
	UpdateEmployee(c *fiber.Ctx) error
	DeleteEmployee(c *fiber.Ctx) error
}

type employeeHandler struct {
	employeeService service.EmployeeService
}

func (e *employeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req request.EmployeeRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	tenantID := claims.TenantID

	if userID == 0 || tenantID == 0 {
		code = "[Handler] CreateEmployee - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] CreateEmployee - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateEmployee - 3"
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	employee := &entity.EmployeeEntity{
		Name:        req.Name,
		CitizenID:   req.CitizenID,
		PhoneNumber: req.PhoneNumber,
		TenantID:    int64(tenantID),
	}

	err = e.employeeService.CreateEmployee(c.Context(), employee)
	if err != nil {
		code := "[HANDLER] CreateEmployee - 4"
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

func (e *employeeHandler) GetEmployeeDetail(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("employeeID")
	employeeID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := e.employeeService.GetDetailEmployeeByTenant(c.Context(), int64(tenantID), employeeID)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = response.EmployeeResponse{
		ID:          result.ID,
		Name:        result.Name,
		CitizenID:   result.CitizenID,
		PhoneNumber: result.PhoneNumber,
		Status:      result.Status,
		CreatedAt:   result.CreatedAt.In(jakartaTZ).Format("02 January 2006"),
	}

	return c.JSON(defaultSuccessResponse)
}

func (e *employeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	var req request.EmployeeRequest
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("employeeID")
	employeeID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
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

	employee := &entity.EmployeeEntity{
		Name:        req.Name,
		CitizenID:   req.CitizenID,
		PhoneNumber: req.PhoneNumber,
	}

	err = e.employeeService.UpdateEmployee(c.Context(), int64(tenantID), employeeID, employee)
	if err != nil {
		code := "[HANDLER] UpdateEmployee - 1"
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

func (e *employeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("employeeID")
	employeeID, err := conv.StringToInt64(id)
	if err != nil {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = e.employeeService.DeleteEmployee(c.Context(), int64(tenantID), employeeID)
	if err != nil {
		code := "[HANDLER] DeleteEmployee - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Employee deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (e *employeeHandler) GetEmployees(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code := "[HANDLER] GetEmployees - 1"
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
			code := "[HANDLER] GetEmployees - 2"
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

	reqEntity := entity.EmployeeQueryString{
		Limit:     limit,
		Page:      page,
		OrderBy:   orderBy,
		Search:    search,
		OrderType: orderType,
		Status:    status,
	}

	role := claims.Role
	tenantID := claims.TenantID
	results, totalData, totalPages, err := e.employeeService.GetEmployees(c.Context(), reqEntity, role, int64(tenantID))

	if err != nil {
		code := "[HANDLER] GetEmployees - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"

	respEmployees := []response.EmployeeResponse{}
	for _, employee := range results {
		respEmployee := response.EmployeeResponse{
			ID:          employee.ID,
			Name:        employee.Name,
			CitizenID:   employee.CitizenID,
			PhoneNumber: employee.PhoneNumber,
			Status:      employee.Status,
			CreatedAt:   employee.CreatedAt.In(jakartaTZ).Format("02 January 2006"),
		}
		respEmployees = append(respEmployees, respEmployee)
	}

	defaultSuccessResponse.Data = respEmployees
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}

	return c.JSON(defaultSuccessResponse)
}

func NewEmployeeHandler(employeeService service.EmployeeService) EmployeeHandler {
	return &employeeHandler{employeeService: employeeService}
}
