package service

import (
	"context"
	"fmt"
	"time"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type AttendanceService interface {
	CreateAttendance(ctx context.Context, req entity.Attendance) (*entity.Attendance, error)
	UpdateAttendance(ctx context.Context, id uuid.UUID, req entity.Attendance) (*entity.Attendance, error)
	DeleteAttendance(ctx context.Context, id uuid.UUID) error
	GetAttendanceByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error)
	GetAttendancesByEmployee(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error)
	GetAttendancesByTenant(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error)
	GetAttendancesByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Attendance, int64, error)
	CheckIn(ctx context.Context, attendanceID uuid.UUID, checkInTime time.Time, location *string) error
	CheckOut(ctx context.Context, attendanceID uuid.UUID, checkOutTime time.Time) error
	UpdateStatus(ctx context.Context, attendanceID uuid.UUID, status entity.AttendanceStatus) error
	GetAttendanceSummary(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*entity.AttendanceSummary, error)
	GetTodayAttendance(ctx context.Context, employeeID uuid.UUID) (*entity.Attendance, error)
}

type attendanceService struct {
	attendanceRepo repository.AttendanceRepository
}

func NewAttendanceService(attendanceRepo repository.AttendanceRepository) AttendanceService {
	return &attendanceService{
		attendanceRepo: attendanceRepo,
	}
}

func (s *attendanceService) CreateAttendance(ctx context.Context, req entity.Attendance) (*entity.Attendance, error) {
	// Validate required fields
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Set default status
	if req.Status == "" {
		req.Status = entity.AttendanceStatusPresent
	}

	// Calculate work hours if check-in and check-out are provided
	if req.CheckIn != nil && req.CheckOut != nil {
		req.CalculateWorkHours()
	}

	if err := s.attendanceRepo.Create(ctx, &req); err != nil {
		log.Errorw("failed to create attendance", "error", err)
		return nil, fmt.Errorf("failed to create attendance: %w", err)
	}

	return &req, nil
}

func (s *attendanceService) UpdateAttendance(ctx context.Context, id uuid.UUID, req entity.Attendance) (*entity.Attendance, error) {
	// Check if attendance exists
	existing, err := s.attendanceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("attendance not found: %w", err)
	}

	// Update fields
	existing.Status = req.Status
	existing.Notes = req.Notes

	if req.CheckIn != nil {
		existing.CheckIn = req.CheckIn
	}
	if req.CheckOut != nil {
		existing.CheckOut = req.CheckOut
	}

	// Recalculate work hours
	if existing.CheckIn != nil && existing.CheckOut != nil {
		existing.CalculateWorkHours()
	}

	if err := s.attendanceRepo.Update(ctx, existing); err != nil {
		log.Errorw("failed to update attendance", "error", err)
		return nil, fmt.Errorf("failed to update attendance: %w", err)
	}

	return existing, nil
}

func (s *attendanceService) DeleteAttendance(ctx context.Context, id uuid.UUID) error {
	// Check if attendance exists
	_, err := s.attendanceRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("attendance not found: %w", err)
	}

	if err := s.attendanceRepo.Delete(ctx, id); err != nil {
		log.Errorw("failed to delete attendance", "error", err)
		return fmt.Errorf("failed to delete attendance: %w", err)
	}

	return nil
}

func (s *attendanceService) GetAttendanceByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error) {
	attendance, err := s.attendanceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("attendance not found: %w", err)
	}
	return attendance, nil
}

func (s *attendanceService) GetAttendancesByEmployee(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error) {
	return s.attendanceRepo.FindByEmployeeID(ctx, employeeID, page, limit)
}

func (s *attendanceService) GetAttendancesByTenant(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error) {
	return s.attendanceRepo.FindByTenantID(ctx, tenantID, page, limit)
}

func (s *attendanceService) GetAttendancesByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Attendance, int64, error) {
	return s.attendanceRepo.FindByPeriod(ctx, tenantID, startDate, endDate, page, limit)
}

func (s *attendanceService) CheckIn(ctx context.Context, attendanceID uuid.UUID, checkInTime time.Time, location *string) error {
	return s.attendanceRepo.CheckIn(ctx, attendanceID, checkInTime, location)
}

func (s *attendanceService) CheckOut(ctx context.Context, attendanceID uuid.UUID, checkOutTime time.Time) error {
	return s.attendanceRepo.CheckOut(ctx, attendanceID, checkOutTime)
}

func (s *attendanceService) UpdateStatus(ctx context.Context, attendanceID uuid.UUID, status entity.AttendanceStatus) error {
	return s.attendanceRepo.UpdateStatus(ctx, attendanceID, status)
}

func (s *attendanceService) GetAttendanceSummary(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*entity.AttendanceSummary, error) {
	return s.attendanceRepo.GetAttendanceSummary(ctx, employeeID, startDate, endDate)
}

func (s *attendanceService) GetTodayAttendance(ctx context.Context, employeeID uuid.UUID) (*entity.Attendance, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.attendanceRepo.FindByEmployeeAndDate(ctx, employeeID, today)
}
