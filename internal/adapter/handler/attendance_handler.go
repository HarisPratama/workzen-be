package handler

import (
	"time"
	"workzen-be/internal/adapter/handler/request"
	"workzen-be/internal/adapter/handler/response"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/service"
	validatorLib "workzen-be/lib/validator"

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
	return &attendanceHandler{attendanceService: attendanceService}
}

func (h *attendanceHandler) CreateAttendance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] CreateAttendance - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.AttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateAttendance - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[HANDLER] CreateAttendance - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	employeeUUID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		code := "[HANDLER] CreateAttendance - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	date, _ := time.Parse("2006-01-02", req.Date)
	attendance := entity.Attendance{
		TenantID:   uuid.UUID{},
		EmployeeID: employeeUUID,
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

	result, err := h.attendanceService.CreateAttendance(c.Context(), attendance)
	if err != nil {
		code := "[HANDLER] CreateAttendance - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Attendance created successfully"
	defaultSuccessResponse.Data = result

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (h *attendanceHandler) UpdateAttendance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] UpdateAttendance - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("attendanceID")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		code := "[HANDLER] UpdateAttendance - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid attendance ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.AttendanceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateAttendance - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
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

	result, err := h.attendanceService.UpdateAttendance(c.Context(), attendanceID, attendance)
	if err != nil {
		code := "[HANDLER] UpdateAttendance - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Attendance updated successfully"
	defaultSuccessResponse.Data = result

	return c.JSON(defaultSuccessResponse)
}

func (h *attendanceHandler) DeleteAttendance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] DeleteAttendance - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("attendanceID")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		code := "[HANDLER] DeleteAttendance - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid attendance ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.attendanceService.DeleteAttendance(c.Context(), attendanceID); err != nil {
		code := "[HANDLER] DeleteAttendance - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Attendance deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *attendanceHandler) GetAttendanceByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetAttendanceByID - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("attendanceID")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		code := "[HANDLER] GetAttendanceByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid attendance ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.attendanceService.GetAttendanceByID(c.Context(), attendanceID)
	if err != nil {
		code := "[HANDLER] GetAttendanceByID - 3"
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

func (h *attendanceHandler) GetAttendancesByEmployee(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetAttendancesByEmployee - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	employeeID := c.Params("employeeID")
	employeeUUID, err := uuid.Parse(employeeID)
	if err != nil {
		code := "[HANDLER] GetAttendancesByEmployee - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	results, total, err := h.attendanceService.GetAttendancesByEmployee(c.Context(), employeeUUID, page, limit)
	if err != nil {
		code := "[HANDLER] GetAttendancesByEmployee - 3"
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

func (h *attendanceHandler) GetAttendancesByTenant(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetAttendancesByTenant - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	results, total, err := h.attendanceService.GetAttendancesByTenant(c.Context(), uuid.UUID{}, page, limit)
	if err != nil {
		code := "[HANDLER] GetAttendancesByTenant - 2"
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

func (h *attendanceHandler) GetAttendancesByPeriod(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetAttendancesByPeriod - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		code := "[HANDLER] GetAttendancesByPeriod - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid start date format"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		code := "[HANDLER] GetAttendancesByPeriod - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid end date format"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	results, total, err := h.attendanceService.GetAttendancesByPeriod(c.Context(), uuid.UUID{}, start, end, page, limit)
	if err != nil {
		code := "[HANDLER] GetAttendancesByPeriod - 4"
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

func (h *attendanceHandler) CheckIn(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] CheckIn - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("attendanceID")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		code := "[HANDLER] CheckIn - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid attendance ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.CheckInRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CheckIn - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	checkInTime := time.Now()
	if req.CheckInTime != "" {
		checkInTime, _ = time.Parse("15:04:05", req.CheckInTime)
	}

	var location *string
	if req.Location != "" {
		location = &req.Location
	}

	if err := h.attendanceService.CheckIn(c.Context(), attendanceID, checkInTime, location); err != nil {
		code := "[HANDLER] CheckIn - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Check-in successful"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *attendanceHandler) CheckOut(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] CheckOut - 1"
		log.Errorw(code, "unauthorized")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("attendanceID")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		code := "[HANDLER] CheckOut - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid attendance ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.CheckOutRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] CheckOut - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	checkOutTime := time.Now()
	if req.CheckOutTime != "" {
		checkOutTime, _ = time.Parse("15:04:05", req.CheckOutTime)
	}

	if err := h.attendanceService.CheckOut(c.Context(), attendanceID, checkOutTime); err != nil {
		code := "[HANDLER] CheckOut - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Check-out successful"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *attendanceHandler) UpdateStatus(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] UpdateStatus - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	id := c.Params("attendanceID")
	attendanceID, err := uuid.Parse(id)
	if err != nil {
		code := "[HANDLER] UpdateStatus - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid attendance ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var req request.UpdateAttendanceStatusRequest
	if err := c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateStatus - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := h.attendanceService.UpdateStatus(c.Context(), attendanceID, entity.AttendanceStatus(req.Status)); err != nil {
		code := "[HANDLER] UpdateStatus - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Attendance status updated successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (h *attendanceHandler) GetAttendanceSummary(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetAttendanceSummary - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	employeeID := c.Query("employee_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if employeeID == "" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "employee_id is required"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	employeeUUID, err := uuid.Parse(employeeID)
	if err != nil {
		code := "[HANDLER] GetAttendanceSummary - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	result, err := h.attendanceService.GetAttendanceSummary(c.Context(), employeeUUID, start, end)
	if err != nil {
		code := "[HANDLER] GetAttendanceSummary - 3"
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

func (h *attendanceHandler) GetTodayAttendance(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.TenantID == 0 {
		code := "[HANDLER] GetTodayAttendance - 1"
		log.Errorw(code, "tenant id is empty")
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	employeeID := c.Params("employeeID")
	if employeeID == "" {
		employeeID = c.Query("employee_id")
	}

	if employeeID == "" {
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "employee_id is required"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	employeeUUID, err := uuid.Parse(employeeID)
	if err != nil {
		code := "[HANDLER] GetTodayAttendance - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid employee ID"
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := h.attendanceService.GetTodayAttendance(c.Context(), employeeUUID)
	if err != nil {
		code := "[HANDLER] GetTodayAttendance - 3"
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
