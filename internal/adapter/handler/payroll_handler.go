package handler

import (
	"time"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	"workzen-be/lib/conv"
	validatorLib "workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type PayrollHandler interface {
	CreatePayroll(c *fiber.Ctx) error
	UpdatePayroll(c *fiber.Ctx) error
	DeletePayroll(c *fiber.Ctx) error
	GetPayrollByID(c *fiber.Ctx) error
	GetPayrollsByTenant(c *fiber.Ctx) error
	GetPayrollsByEmployee(c *fiber.Ctx) error
	ProcessPayroll(c *fiber.Ctx) error
	MarkAsPaid(c *fiber.Ctx) error
	CalculatePayrollSummary(c *fiber.Ctx) error
}

type payrollHandler struct {
	payrollService service.PayrollService
}

func NewPayrollHandler(payrollService service.PayrollService) PayrollHandler {
	return &payrollHandler{payrollService: payrollService}
}

func (h *payrollHandler) CreatePayroll(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] CreatePayroll - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.PayrollRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreatePayroll - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[HANDLER] CreatePayroll - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	periodStart, _ := time.ParseInLocation("2006-01-02", req.PeriodStart, jakartaTZ)
	periodEnd, _ := time.ParseInLocation("2006-01-02", req.PeriodEnd, jakartaTZ)

	payroll := entity.Payroll{
		TenantID:    int64(claims.TenantID),
		EmployeeID:  req.EmployeeID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		BasicSalary: req.BasicSalary,
		Allowances:  req.Allowances,
		Deductions:  req.Deductions,
		Tax:         req.Tax,
		Notes:       req.Notes,
	}

	result, err := h.payrollService.CreatePayroll(c.Context(), payroll)
	if err != nil {
		code := "[HANDLER] CreatePayroll - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Payroll created successfully"
	defaultSuccessResponse.Data = result

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (h *payrollHandler) UpdatePayroll(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] UpdatePayroll - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("payrollID")
	payrollID, err := conv.StringToInt64(id)
	if err != nil {
		code := "[HANDLER] UpdatePayroll - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid payroll ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.PayrollUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdatePayroll - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	payroll := entity.Payroll{
		BasicSalary: req.BasicSalary,
		Allowances:  req.Allowances,
		Deductions:  req.Deductions,
		Tax:         req.Tax,
		Notes:       req.Notes,
	}

	result, err := h.payrollService.UpdatePayroll(c.Context(), payrollID, payroll)
	if err != nil {
		code := "[HANDLER] UpdatePayroll - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Payroll updated successfully"
	defaultSuccessResponse.Data = result

	return c.JSON(defaultSuccessResponse)
}

func (h *payrollHandler) DeletePayroll(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] DeletePayroll - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("payrollID")
	payrollID, err := conv.StringToInt64(id)
	if err != nil {
		code := "[HANDLER] DeletePayroll - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid payroll ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.payrollService.DeletePayroll(c.Context(), payrollID); err != nil {
		code := "[HANDLER] DeletePayroll - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Payroll deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *payrollHandler) GetPayrollByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetPayrollByID - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("payrollID")
	payrollID, err := conv.StringToInt64(id)
	if err != nil {
		code := "[HANDLER] GetPayrollByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid payroll ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.payrollService.GetPayrollByID(c.Context(), payrollID)
	if err != nil {
		code := "[HANDLER] GetPayrollByID - 3"
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

func (h *payrollHandler) GetPayrollsByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetPayrollsByTenant - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	results, total, err := h.payrollService.GetPayrollsByTenant(c.Context(), int64(claims.TenantID), page, limit)
	if err != nil {
		code := "[HANDLER] GetPayrollsByTenant - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = results
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(total),
		Page:         page,
		PerPage:      limit,
	}

	return c.JSON(defaultSuccessResponse)
}

func (h *payrollHandler) GetPayrollsByEmployee(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetPayrollsByEmployee - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	employeeID := c.Params("employeeID")
	employeeIntID, err := conv.StringToInt64(employeeID)
	if err != nil {
		code := "[HANDLER] GetPayrollsByEmployee - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	results, total, err := h.payrollService.GetPayrollsByEmployee(c.Context(), employeeIntID, page, limit)
	if err != nil {
		code := "[HANDLER] GetPayrollsByEmployee - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = results
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(total),
		Page:         page,
		PerPage:      limit,
	}

	return c.JSON(defaultSuccessResponse)
}

func (h *payrollHandler) ProcessPayroll(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] ProcessPayroll - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("payrollID")
	payrollID, err := conv.StringToInt64(id)
	if err != nil {
		code := "[HANDLER] ProcessPayroll - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid payroll ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.payrollService.ProcessPayroll(c.Context(), payrollID); err != nil {
		code := "[HANDLER] ProcessPayroll - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Payroll processed successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *payrollHandler) MarkAsPaid(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] MarkAsPaid - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("payrollID")
	payrollID, err := conv.StringToInt64(id)
	if err != nil {
		code := "[HANDLER] MarkAsPaid - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid payroll ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.MarkAsPaidRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] MarkAsPaid - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	paidAt, _ := time.ParseInLocation("2006-01-02", req.PaymentDate, jakartaTZ)

	if err := h.payrollService.MarkAsPaid(c.Context(), payrollID, paidAt); err != nil {
		code := "[HANDLER] MarkAsPaid - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Payroll marked as paid"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *payrollHandler) CalculatePayrollSummary(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] CalculatePayrollSummary - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, _ := time.ParseInLocation("2006-01-02", startDateStr, jakartaTZ)
	endDate, _ := time.ParseInLocation("2006-01-02", endDateStr, jakartaTZ)

	result, err := h.payrollService.CalculatePayrollSummary(c.Context(), int64(claims.TenantID), startDate, endDate)
	if err != nil {
		code := "[HANDLER] CalculatePayrollSummary - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Success"
	defaultSuccessResponse.Data = result

	return c.JSON(defaultSuccessResponse)
}
