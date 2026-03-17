package handler

import (
	"bwanews/internal/adapter/handler/request"
	"bwanews/internal/adapter/handler/response"
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/service"
	"bwanews/lib/conv"
	validatorLib "bwanews/lib/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type PayrollHandler interface {
	CreatePayroll(c *fiber.Ctx) error
	UpdatePayroll(c *fiber.Ctx) error
	DeletePayroll(c *fiber.Ctx) error
	GetPayrollByID(c *fiber.Ctx) error
	GetPayrollsByTenant(c *fiber.Ctx) error
	GetPayrollsByEmployee(c *fiber.Ctx) error
	GetPayrollsByPeriod(c *fiber.Ctx) error
	ProcessPayroll(c *fiber.Ctx) error
	MarkAsPaid(c *fiber.Ctx) error
	AddPayrollItem(c *fiber.Ctx) error
	RemovePayrollItem(c *fiber.Ctx) error
	GetPayrollSummary(c *fiber.Ctx) error
}

type payrollHandler struct {
	payrollService service.PayrollService
}

func NewPayrollHandler(payrollService service.PayrollService) PayrollHandler {
	return &payrollHandler{
		payrollService: payrollService,
	}
}

func (h *payrollHandler) CreatePayroll(c *fiber.Ctx) error {
	var req request.PayrollRequest
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized access",
			},
		})
	}

	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] CreatePayroll - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	payroll := entity.Payroll{
		TenantID:      uuid.UUID{},
		EmployeeID:    req.EmployeeID,
		PeriodStart:   req.PeriodStart,
		PeriodEnd:     req.PeriodEnd,
		BasicSalary:   req.BasicSalary,
		Allowances:    req.Allowances,
		Deductions:    req.Deductions,
		Tax:           req.Tax,
		NetSalary:     req.NetSalary,
		Status:        entity.PayrollStatusDraft,
		Notes:         req.Notes,
	}

	createdPayroll, err := h.payrollService.CreatePayroll(c.Context(), payroll)
	if err != nil {
		log.Errorw("[Handler] CreatePayroll - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Payroll created successfully",
		},
		Data: createdPayroll,
	})
}

func (h *payrollHandler) UpdatePayroll(c *fiber.Ctx) error {
	id := c.Params("id")
	payrollID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid payroll ID",
			},
		})
	}

	var req request.PayrollUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] UpdatePayroll - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	payroll := entity.Payroll{
		ID:          payrollID,
		BasicSalary: req.BasicSalary,
		Allowances:  req.Allowances,
		Deductions:  req.Deductions,
		Tax:         req.Tax,
		NetSalary:   req.NetSalary,
		Notes:       req.Notes,
	}

	updatedPayroll, err := h.payrollService.UpdatePayroll(c.Context(), payrollID, payroll)
	if err != nil {
		log.Errorw("[Handler] UpdatePayroll - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Payroll updated successfully",
		},
		Data: updatedPayroll,
	})
}

func (h *payrollHandler) DeletePayroll(c *fiber.Ctx) error {
	id := c.Params("id")
	payrollID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid payroll ID",
			},
		})
	}

	if err := h.payrollService.DeletePayroll(c.Context(), payrollID); err != nil {
		log.Errorw("[Handler] DeletePayroll - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Payroll deleted successfully",
		},
	})
}

func (h *payrollHandler) GetPayrollByID(c *fiber.Ctx) error {
	id := c.Params("id")
	payrollID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid payroll ID",
			},
		})
	}

	payroll, err := h.payrollService.GetPayrollByID(c.Context(), payrollID)
	if err != nil {
		log.Errorw("[Handler] GetPayrollByID - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Success",
		},
		Data: payroll,
	})
}

func (h *payrollHandler) GetPayrollsByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	tenantID := claims.TenantID

	if tenantID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized access",
			},
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	payrolls, total, err := h.payrollService.GetPayrollsByTenant(c.Context(), uuid.UUID{}, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetPayrollsByTenant - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Success",
		},
		Data: response.PaginationData{
			Data:  payrolls,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *payrollHandler) GetPayrollsByEmployee(c *fiber.Ctx) error {
	employeeID := c.Params("employeeId")
	employeeUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid employee ID",
			},
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	payrolls, total, err := h.payrollService.GetPayrollsByEmployee(c.Context(), employeeUUID, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetPayrollsByEmployee - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response