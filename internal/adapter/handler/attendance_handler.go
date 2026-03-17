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

type AttendanceHandler interface {
	CreateAttendance(c *fiber.Ctx) error
	UpdateAttendance(c *fiber.Ctx) error
	DeleteAttendance(c *fiber.Ctx) error
	GetAttendanceByID(c *fiber.Ctx) error
	GetAttendancesByEmployee(c *fiber.Ctx) error
	GetAttendancesByTenant(c *fiber.Ctx) error
	GetAttendancesByPeriod(c *fiber.Ctx) error
	CheckIn(c *fiber.Ctx) error
	CheckOut(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	GetAttendanceSummary(c *fiber.Ctx) error
	GetTodayAttendance(c *fiber.Ctx) error
}

type attendanceHandler struct {
	attendanceService service.AttendanceService
}

func NewAttendanceHandler(attendanceService service.AttendanceService) AttendanceHandler {
	return &attendanceHandler{
		attendanceService: attendanceService,
	}
}

func (h *attendanceHandler) CreateAttendance(c *fiber.Ctx) error {
	var req request.AttendanceRequest
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
		log.Errorw("[Handler] CreateAttendance - BodyParser", err)
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

	date, _ := time.Parse("2006-01-02", req.Date)
	attendance := entity.Attendance{
		TenantID:   uuid.UUID{},
		EmployeeID: req.EmployeeID,
		Date:       date,
		Status:     entity.AttendanceStatus(req.Status),
		Type:       entity.AttendanceType(req.Type),
		Notes:      req.Notes,
	}

	if req.CheckIn != "" {
		checkIn, _ := time.Parse("15:04:05", req.CheckIn)
		attendance.CheckIn = &checkIn
	}

	if req.CheckOut != "" {
		checkOut, _ := time.Parse("15:04:05", req.CheckOut)
		attendance.CheckOut = &checkOut
	}

	if attendance.CheckIn != nil && attendance.CheckOut != nil {
		attendance.CalculateWorkHours()
	}

	createdAttendance, err := h.attendanceService.CreateAttendance(c.Context(), attendance)
	if err != nil {
		log.Errorw("[Handler] CreateAttendance - Service", err)
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
			Message: "Attendance created successfully",
		},
		Data: createdAttendance,
	})
}

func (h *attendanceHandler) UpdateAttendance(c *fiber.Ctx) error {
	id := c.Params("id")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid attendance ID",
			},
		})
	}

	var req request.AttendanceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] UpdateAttendance - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	attendance := entity.Attendance{
		ID:     attendanceID,
		Status: entity.AttendanceStatus(req.Status),
		Notes:  req.Notes,
	}

	if req.CheckIn != "" {
		checkIn, _ := time.Parse("15:04:05", req.CheckIn)
		attendance.CheckIn = &checkIn
	}

	if req.CheckOut != "" {
		checkOut, _ := time.Parse("15:04:05", req.CheckOut)
		attendance.CheckOut = &checkOut
	}

	updatedAttendance, err := h.attendanceService.UpdateAttendance(c.Context(), attendanceID, attendance)
	if err != nil {
		log.Errorw("[Handler] UpdateAttendance - Service", err)
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
			Message: "Attendance updated successfully",
		},
		Data: updatedAttendance,
	})
}

func (h *attendanceHandler) DeleteAttendance(c *fiber.Ctx) error {
	id := c.Params("id")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid attendance ID",
			},
		})
	}

	if err := h.attendanceService.DeleteAttendance(c.Context(), attendanceID); err != nil {
		log.Errorw("[Handler] DeleteAttendance - Service", err)
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
			Message: "Attendance deleted successfully",
		},
	})
}

func (h *attendanceHandler) GetAttendanceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid attendance ID",
			},
		})
	}

	attendance, err := h.attendanceService.GetAttendanceByID(c.Context(), attendanceID)
	if err != nil {
		log.Errorw("[Handler] GetAttendanceByID - Service", err)
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
		Data: attendance,
	})
}

func (h *attendanceHandler) GetAttendancesByEmployee(c *fiber.Ctx) error {
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

	attendances, total, err := h.attendanceService.GetAttendancesByEmployee(c.Context(), employeeUUID, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetAttendancesByEmployee - Service", err)
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
			Data:  attendances,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *attendanceHandler) GetAttendancesByTenant(c *fiber.Ctx) error {
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

	attendances, total, err := h.attendanceService.GetAttendancesByTenant(c.Context(), uuid.UUID{}, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetAttendancesByTenant - Service", err)
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
			Data:  attendances,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *attendanceHandler) GetAttendancesByPeriod(c *fiber.Ctx) error {
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

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid start date format",
			},
		})
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid end date format",
			},
		})
	}

	attendances, total, err := h.attendanceService.GetAttendancesByPeriod(c.Context(), uuid.UUID{}, start, end, page, limit)
	if err != nil {
		log.Errorw("[Handler] GetAttendancesByPeriod - Service", err)
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
			Data:  attendances,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

func (h *attendanceHandler) CheckIn(c *fiber.Ctx) error {
	id := c.Params("id")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid attendance ID",
			},
		})
	}

	var req request.CheckInRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] CheckIn - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	checkInTime := time.Now()
	if req.CheckInTime != "" {
		checkInTime, _ = time.Parse("15:04:05", req.CheckInTime)
	}

	if err := h.attendanceService.CheckIn(c.Context(), attendanceID, checkInTime, req.Location); err != nil {
		log.Errorw("[Handler] CheckIn - Service", err)
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
			Message: "Check-in successful",
		},
	})
}

func (h *attendanceHandler) CheckOut(c *fiber.Ctx) error {
	id := c.Params("id")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid attendance ID",
			},
		})
	}

	var req request.CheckOutRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorw("[Handler] CheckOut - BodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		})
	}

	checkOutTime := time.Now()
	if req.CheckOutTime != "" {
		checkOutTime, _ = time.Parse("15:04:05", req.CheckOutTime)
	}

	if err := h.attendanceService.CheckOut(c.Context(), attendanceID, checkOutTime); err != nil {
	