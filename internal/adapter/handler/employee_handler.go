package handler

import (
	"bwanews/internal/adapter/handler/request"
	"bwanews/internal/adapter/handler/response"
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/service"
	"bwanews/lib/conv"
	validatorLib "bwanews/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type EmployeeHandler interface {
	GetEmployees(c *fiber.Ctx) error
	CreateEmployee(c *fiber.Ctx) error
}

type employeeHandler struct {
	employeeService service.EmployeeService
}

func (e *employeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req request.EmployeeRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	role := claims.Role

	if userID == 0 {
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

	reqEntity := entity.EmployeeEntity{
		Name:        req.Name,
		CitizenID:   req.CitizenID,
		PhoneNumber: req.PhoneNumber,
	}

	err = e.employeeService.CreateEmployee(c.Context(), reqEntity)

	panic("implement me")
}

func (e *employeeHandler) GetEmployees(c *fiber.Ctx) error {
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

	results, totalData, totalPages, err := e.employeeService.GetEmployees(c.Context(), reqEntity)
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
			CreatedAt:   employee.CreatedAt.Local().Format("02 January 2006"),
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
