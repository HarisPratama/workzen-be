package handler

import (
	"context"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type EmployeeAssignmentHandler interface {
	GetAssignments(c *fiber.Ctx) error
	GetAssignmentByID(c *fiber.Ctx) error
	CreateAssignment(c *fiber.Ctx) error
	UpdateAssignment(c *fiber.Ctx) error
	DeleteAssignment(c *fiber.Ctx) error
}

type employeeAssignmentHandler struct {
	assignmentService service.EmployeeAssignmentService
}

func (h *employeeAssignmentHandler) GetAssignments(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetAssignments - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	query := entity.EmployeeAssignmentQueryString{
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

	results, totalData, totalPages, err := h.assignmentService.GetEmployeeAssignmentsByTenant(context.Background(), int64(tenantID), query)
	if err != nil {
		code := "[HANDLER] GetAssignments - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	var respData []response.EmployeeAssignmentResponse
	for _, item := range results {
		endDate := ""
		if !item.EndDate.IsZero() {
			endDate = item.EndDate.In(jakartaTZ).Format("2006-01-02")
		}
		expectedEndDate := ""
		if !item.ExpectedEndDate.IsZero() {
			expectedEndDate = item.ExpectedEndDate.In(jakartaTZ).Format("2006-01-02")
		}

		respData = append(respData, response.EmployeeAssignmentResponse{
			ID:              item.ID,
			AssignmentType:  item.AssignmentType,
			StartDate:       item.StartDate.In(jakartaTZ).Format("2006-01-02"),
			EndDate:         endDate,
			ExpectedEndDate: expectedEndDate,
			Status:          item.Status,
			Role:            item.Role,
			Position:        item.Position,
			Location:        item.Location,
			RemoteType:      item.RemoteType,
			BillingRate:     item.BillingRate,
			CostRate:        item.CostRate,
			Currency:        item.Currency,
			HoursPerWeek:    item.HoursPerWeek,
			Notes:           item.Notes,
			Employee: response.EmployeeAssignmentEmployeeResp{
				ID:   item.Employee.ID,
				Name: item.Employee.Name,
			},
			Client: response.EmployeeAssignmentClientResp{
				ID:          item.Client.ID,
				CompanyName: item.Client.CompanyName,
			},
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

func (h *employeeAssignmentHandler) GetAssignmentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetAssignmentByID - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("assignmentID")
	assignmentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] GetAssignmentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid assignment ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.assignmentService.GetEmployeeAssignmentByID(context.Background(), assignmentID)
	if err != nil {
		code := "[HANDLER] GetAssignmentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	endDate := ""
	if !result.EndDate.IsZero() {
		endDate = result.EndDate.In(jakartaTZ).Format("2006-01-02")
	}
	expectedEndDate := ""
	if !result.ExpectedEndDate.IsZero() {
		expectedEndDate = result.ExpectedEndDate.In(jakartaTZ).Format("2006-01-02")
	}

	respData := response.EmployeeAssignmentResponse{
		ID:              result.ID,
		AssignmentType:  result.AssignmentType,
		StartDate:       result.StartDate.In(jakartaTZ).Format("2006-01-02"),
		EndDate:         endDate,
		ExpectedEndDate: expectedEndDate,
		Status:          result.Status,
		Role:            result.Role,
		Position:        result.Position,
		Location:        result.Location,
		RemoteType:      result.RemoteType,
		BillingRate:     result.BillingRate,
		CostRate:        result.CostRate,
		Currency:        result.Currency,
		HoursPerWeek:    result.HoursPerWeek,
		Notes:           result.Notes,
		Employee: response.EmployeeAssignmentEmployeeResp{
			ID:   result.Employee.ID,
			Name: result.Employee.Name,
		},
		Client: response.EmployeeAssignmentClientResp{
			ID:          result.Client.ID,
			CompanyName: result.Client.CompanyName,
		},
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = respData

	return c.JSON(defaultSuccessResponse)
}

func (h *employeeAssignmentHandler) CreateAssignment(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] CreateAssignment - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	tenantID := claims.TenantID

	var req entity.EmployeeAssignmentEntityRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateAssignment - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[HANDLER] CreateAssignment - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.assignmentService.CreateEmployeeAssignment(context.Background(), req, int64(tenantID)); err != nil {
		code := "[HANDLER] CreateAssignment - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Assignment created successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (h *employeeAssignmentHandler) UpdateAssignment(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] UpdateAssignment - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("assignmentID")
	assignmentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] UpdateAssignment - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid assignment ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req entity.EmployeeAssignmentUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateAssignment - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.assignmentService.UpdateEmployeeAssignment(context.Background(), assignmentID, req); err != nil {
		code := "[HANDLER] UpdateAssignment - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Assignment updated successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *employeeAssignmentHandler) DeleteAssignment(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] DeleteAssignment - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("assignmentID")
	assignmentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code := "[HANDLER] DeleteAssignment - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid assignment ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.assignmentService.DeleteEmployeeAssignment(context.Background(), assignmentID); err != nil {
		code := "[HANDLER] DeleteAssignment - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Assignment deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func NewEmployeeAssignmentHandler(assignmentService service.EmployeeAssignmentService) EmployeeAssignmentHandler {
	return &employeeAssignmentHandler{assignmentService: assignmentService}
}
