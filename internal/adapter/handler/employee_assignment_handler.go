package handler

import (
	"bwanews/internal/adapter/handler/request"
	"bwanews/internal/adapter/handler/response"
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/service"
	"bwanews/lib/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type EmployeeAssignmentHandler interface {
	CreateAssignment(c *fiber.Ctx) error
	UpdateAssignment(c *fiber.Ctx) error
	DeleteAssignment(c *fiber.Ctx) error
	GetAssignmentByID(c *fiber.Ctx) error
	GetAssignmentsByEmployee(c *fiber.Ctx) error
	GetAssignmentsByProject(c *fiber.Ctx) error
	GetAssignmentsByTenant(c *fiber.Ctx) error
	GetActiveAssignmentsByEmployee(c *fiber.Ctx) error
	StartAssignment(c *fiber.Ctx) error
	EndAssignment(c *fiber.Ctx) error
	UpdateAssignmentStatus(c *fiber.Ctx) error
	GetAssignmentUtilization(c *fiber.Ctx) error
}

type employeeAssignmentHandler struct {
	assignmentService service.EmployeeAssignmentService
}

func NewEmployeeAssignmentHandler(assignmentService service.EmployeeAssignmentService) EmployeeAssignmentHandler {
	return &employeeAssignmentHandler{
		assignmentService: assignmentService,
	}
}

func (h *employeeAssignmentHandler) CreateAssignment(c *fiber.Ctx) error {
	var req request.EmployeeAssignmentRequest
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
		log.Errorw("[Handler] CreateAssignment - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)

	assignment := entity.EmployeeAssignment{
		TenantID:              uuid.UUID{},
		EmployeeID:            req.EmployeeID,
		ProjectID:             req.ProjectID,
		Role:                  req.Role,
		StartDate:             startDate,
		AllocationPercentage:  req.AllocationPercentage,
		Billable:              req.Billable,
		Status:                entity.AssignmentStatusPending,
		Notes:                 req.Notes,
	}

	if req.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		assignment.EndDate = &endDate
	}

	createdAssignment, err := h.assignmentService.CreateAssignment(c.Context(), assignment)
	if err != nil {
		log.Errorw("[Handler] CreateAssignment - Service", err)
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
			Message: "Assignment created successfully",
		},
		Data: createdAssignment,
	})
}

func (h *employeeAssignmentHandler) UpdateAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid assignment ID",
			},
		})
	}

	var req request.EmployeeAssignmentUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] UpdateAssignment - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

	assignment := entity.EmployeeAssignment{
		ID:                   assignmentID,
		Role:                 req.Role,
		AllocationPercentage: req.AllocationPercentage,
		Billable:             req.Billable,
		Notes:                req.Notes,
	}

	if req.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		assignment.EndDate = &endDate
	}

	updatedAssignment, err := h.assignmentService.UpdateAssignment(c.Context(), assignmentID, assignment)
	if err != nil {
		log.Errorw("[Handler] UpdateAssignment - Service", err)
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
			Message: "Assignment updated successfully",
		},
		Data: updatedAssignment,
	})
}

func (h *employeeAssignmentHandler) DeleteAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid assignment ID",
			},
		})
	}

	if err := h.assignmentService.DeleteAssignment(c.Context(), assignmentID); err != nil {
		log.Errorw("[Handler] DeleteAssignment - Service", err)
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
			Message: "Assignment deleted successfully",
		},
	})
}

func (h *employeeAssignmentHandler) GetAssignmentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid assignment ID",
			},
		})
	}

	assignment, err := h.assignmentService.GetAssignmentByID(c.Context(), assignmentID)
	if err != nil {
		log.Errorw("[Handler] GetAssignmentByID - Service", err)
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
		Data: assignment,
	})
}

func (h *employeeAssignmentHandler) GetAssignmentsByEmployee(c *fiber.Ctx) error {
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

	assignments, total, err := h.assignmentService.GetAssignmentsByEmployee(c.Context(), employeeUUID, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetAssignmentsByEmployee - Service", err)
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
			Data:  assignments,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *employeeAssignmentHandler) GetAssignmentsByProject(c *fiber.Ctx) error {
	projectID := c.Params("projectId")
	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid project ID",
			},
		})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	assignments, total, err := h.assignmentService.GetAssignmentsByProject(c.Context(), projectUUID, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetAssignmentsByProject - Service", err)
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
			Data:  assignments,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *employeeAssignmentHandler) GetAssignmentsByTenant(c *fiber.Ctx) error {
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

	assignments, total, err := h.assignmentService.GetAssignmentsByTenant(c.Context(), uuid.UUID{}, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetAssignmentsByTenant - Service", err)
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
			Data:  assignments,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *employeeAssignmentHandler) GetActiveAssignmentsByEmployee(c *fiber.Ctx) error {
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

	assignments, err := h.assignmentService.GetActiveAssignmentsByEmployee(c.Context(), employeeUUID)
	if err != nil {
		log.Errorw("[Handler] GetActiveAssignmentsByEmployee - Service", err)
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
		Data: assignments,
	})
}

func (h *employeeAssignmentHandler) StartAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid assignment ID",
			},
		})
	}

	if err := h.assignmentService.StartAssignment(c.Context(), assignmentID); err != nil {
		log.Errorw("[Handler] StartAssignment - Service", err)
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
			Message: "Assignment started successfully",
		},
	})
}

func (h *employeeAssignmentHandler) EndAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid assignment ID",
			},
		})
	}

	var req request.EndAssignmentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] EndAssignment - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	if err := h.assignmentService.EndAssignment(c.Context(), assignmentID, req.EndDate, req.Reason); err != nil {
		log.Errorw("[Handler] EndAssignment - Service", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		})
	}

