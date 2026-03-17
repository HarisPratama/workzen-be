package repository

import (
	"bwanews/internal/core/domain/entity"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.Attendance) error
	Update(ctx context.Context, attendance *entity.Attendance) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error)
	FindByEmployeeAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*entity.Attendance, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error)
	FindByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Attendance, int64, error)
	GetAttendanceSummary(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*entity.AttendanceSummary, error)
	CheckIn(ctx context.Context, attendanceID uuid.UUID, checkInTime time.Time, location *string) error
	CheckOut(ctx context.Context, attendanceID uuid.UUID, checkOutTime time.Time) error
	UpdateStatus(ctx context.Context, attendanceID uuid.UUID, status entity.AttendanceStatus) error
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *entity.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

func (r *attendanceRepository) Update(ctx context.Context, attendance *entity.Attendance) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}

func (r *attendanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Attendance{}, id).Error
}

func (r *attendanceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error) {
	var attendance entity.Attendance
	err := r.db.WithContext(ctx).Preload("Employee").Preload("Tenant").First(&attendance, id).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error) {
	var attendances []entity.Attendance
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&entity.Attendance{}).Where("employee_id = ?", employeeID)
	db.Count(&total)
	err := db.Order("date desc").Offset(offset).Limit(limit).Find(&attendances).Error

	return attendances, total, err
}

func (r *attendanceRepository) FindByEmployeeAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*entity.Attendance, error) {
	var attendance entity.Attendance
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND date BETWEEN ? AND ?", employeeID, startOfDay, endOfDay).
		First(&attendance).Error

	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Attendance, int64, error) {
	var attendances []entity.Attendance
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&entity.Attendance{}).Where("tenant_id = ?", tenantID)
	db.Count(&total)
	err := db.Order("date desc").Offset(offset).Limit(limit).Preload("Employee").Find(&attendances).Error

	return attendances, total, err
}

func (r *attendanceRepository) FindByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Attendance, int64, error) {
	var attendances []entity.Attendance
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Where("tenant_id = ?", tenantID).
		Where("date BETWEEN ? AND ?", startDate, endDate)
	db.Count(&total)
	err := db.Order("date desc").Offset(offset).Limit(limit).Preload("Employee").Find(&attendances).Error

	return attendances, total, err
}

func (r *attendanceRepository) GetAttendanceSummary(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*entity.AttendanceSummary, error) {
	var summary entity.AttendanceSummary

	db := r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Where("employee_id = ?", employeeID).
		Where("date BETWEEN ? AND ?", startDate, endDate)

	// Count total days
	db.Count(&summary.TotalDays)

	// Count by status
	r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Where("employee_id = ? AND status = ? AND date BETWEEN ? AND ?", employeeID, entity.AttendanceStatusPresent, startDate, endDate).
		Count(&summary.PresentDays)

	r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Where("employee_id = ? AND status = ? AND date BETWEEN ? AND ?", employeeID, entity.AttendanceStatusAbsent, startDate, endDate).
		Count(&summary.AbsentDays)

	r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Where("employee_id = ? AND status = ? AND date BETWEEN ? AND ?", employeeID, entity.AttendanceStatusLate, startDate, endDate).
		Count(&summary.LateDays)

	r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Where("employee_id = ? AND status = ? AND date BETWEEN ? AND ?", employeeID, entity.AttendanceStatusHalfDay, startDate, endDate).
		Count(&summary.HalfDays)

	r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Where("employee_id = ? AND status = ? AND date BETWEEN ? AND ?", employeeID, entity.AttendanceStatusOnLeave, startDate, endDate).
		Count(&summary.OnLeaveDays)

	// Calculate total work hours
	var totalWorkHours float64
	r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Select("COALESCE(SUM(work_hours), 0)").
		Where("employee_id = ? AND date BETWEEN ? AND ?", employeeID, startDate, endDate).
		Scan(&totalWorkHours)
	summary.TotalWorkHours = totalWorkHours

	// Calculate total overtime
	var totalOvertime float64
	r.db.WithContext(ctx).Model(&entity.Attendance{}).
		Select("COALESCE(SUM(overtime_hours), 0)").
		Where("employee_id = ? AND date BETWEEN ? AND ?", employeeID, startDate, endDate).
		Scan(&totalOvertime)
	summary.TotalOvertime = totalOvertime

	return &summary, nil
}

func (r *attendanceRepository) CheckIn(ctx context.Context, attendanceID uuid.UUID, checkInTime time.Time, location *string) error {
	updates := map[string]interface{}{
		"check_in": checkInTime,
	}
	if location != nil {
		updates["location"] = *location
	}
	return r.db.WithContext(ctx).Model(&entity.Attendance{}).Where("id = ?", attendanceID).Updates(updates).Error
}

func (r *attendanceRepository) CheckOut(ctx context.Context, attendanceID uuid.UUID, checkOutTime time.Time) error {
	return r.db.WithContext(ctx).Model(&entity.Attendance{}).Where("id = ?", attendanceID).Update("check_out", checkOutTime).Error
}

func (r *attendanceRepository) UpdateStatus(ctx context.Context, attendanceID uuid.UUID, status entity.AttendanceStatus) error {
	return r.db.WithContext(ctx).Model(&entity.Attendance{}).Where("id = ?", attendanceID).Update("status", status).Error
}

// Ensure attendanceRepository implements AttendanceRepository interface
var _ AttendanceRepository = (*attendanceRepository)(nil)