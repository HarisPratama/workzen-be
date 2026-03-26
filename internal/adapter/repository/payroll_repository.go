package repository

import (
	"context"
	"time"
	"workzen-be/internal/core/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayrollRepository interface {
	Create(ctx context.Context, payroll *entity.Payroll) error
	Update(ctx context.Context, payroll *entity.Payroll) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Payroll, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error)
	FindByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Payroll, int64, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status entity.PayrollStatus, page, limit int) ([]entity.Payroll, int64, error)
	ProcessPayroll(ctx context.Context, id uuid.UUID) error
	MarkAsPaid(ctx context.Context, id uuid.UUID, paidAt time.Time) error
	AddPayrollItem(ctx context.Context, item *entity.PayrollItem) error
	DeletePayrollItem(ctx context.Context, itemID uuid.UUID) error
	GetPayrollItems(ctx context.Context, payrollID uuid.UUID) ([]entity.PayrollItem, error)
	CalculatePayrollSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*entity.PayrollSummary, error)
}

type payrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &payrollRepository{db: db}
}

func (r *payrollRepository) Create(ctx context.Context, payroll *entity.Payroll) error {
	return r.db.WithContext(ctx).Create(payroll).Error
}

func (r *payrollRepository) Update(ctx context.Context, payroll *entity.Payroll) error {
	return r.db.WithContext(ctx).Save(payroll).Error
}

func (r *payrollRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Payroll{}, id).Error
}

func (r *payrollRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payroll, error) {
	var payroll entity.Payroll
	err := r.db.WithContext(ctx).Preload("Employee").Preload("Tenant").First(&payroll, id).Error
	if err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (r *payrollRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error) {
	var payrolls []entity.Payroll
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&entity.Payroll{}).Where("tenant_id = ?", tenantID)
	db.Count(&total)
	err := db.Order("created_at desc").Offset(offset).Limit(limit).Preload("Employee").Find(&payrolls).Error

	return payrolls, total, err
}

func (r *payrollRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error) {
	var payrolls []entity.Payroll
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&entity.Payroll{}).Where("employee_id = ?", employeeID)
	db.Count(&total)
	err := db.Order("period_start desc").Offset(offset).Limit(limit).Find(&payrolls).Error

	return payrolls, total, err
}

func (r *payrollRepository) FindByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Payroll, int64, error) {
	var payrolls []entity.Payroll
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&entity.Payroll{}).
		Where("tenant_id = ?", tenantID).
		Where("(period_start <= ? AND period_end >= ?)", endDate, startDate)
	db.Count(&total)
	err := db.Order("period_start desc").Offset(offset).Limit(limit).Preload("Employee").Find(&payrolls).Error

	return payrolls, total, err
}

func (r *payrollRepository) FindByStatus(ctx context.Context, tenantID uuid.UUID, status entity.PayrollStatus, page, limit int) ([]entity.Payroll, int64, error) {
	var payrolls []entity.Payroll
	var total int64

	offset := (page - 1) * limit

	db := r.db.WithContext(ctx).Model(&entity.Payroll{}).Where("tenant_id = ? AND status = ?", tenantID, status)
	db.Count(&total)
	err := db.Order("created_at desc").Offset(offset).Limit(limit).Preload("Employee").Find(&payrolls).Error

	return payrolls, total, err
}

func (r *payrollRepository) ProcessPayroll(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.Payroll{}).Where("id = ?", id).Update("status", entity.PayrollStatusProcessed).Error
}

func (r *payrollRepository) MarkAsPaid(ctx context.Context, id uuid.UUID, paidAt time.Time) error {
	return r.db.WithContext(ctx).Model(&entity.Payroll{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":  entity.PayrollStatusPaid,
		"paid_at": paidAt,
	}).Error
}

func (r *payrollRepository) AddPayrollItem(ctx context.Context, item *entity.PayrollItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *payrollRepository) DeletePayrollItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.PayrollItem{}, itemID).Error
}

func (r *payrollRepository) GetPayrollItems(ctx context.Context, payrollID uuid.UUID) ([]entity.PayrollItem, error) {
	var items []entity.PayrollItem
	err := r.db.WithContext(ctx).Where("payroll_id = ?", payrollID).Find(&items).Error
	return items, err
}

func (r *payrollRepository) CalculatePayrollSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*entity.PayrollSummary, error) {
	type result struct {
		TotalPayrolls    int64   `gorm:"column:total_payrolls"`
		TotalBasicSalary float64 `gorm:"column:total_basic_salary"`
		TotalNetSalary   float64 `gorm:"column:total_net_salary"`
	}

	var res result

	err := r.db.WithContext(ctx).Model(&entity.Payroll{}).
		Select("COUNT(*) as total_payrolls, SUM(basic_salary) as total_basic_salary, SUM(net_salary) as total_net_salary").
		Where("tenant_id = ? AND period_start >= ? AND period_end <= ?", tenantID, startDate, endDate).
		Scan(&res).Error

	if err != nil {
		return nil, err
	}

	summary := &entity.PayrollSummary{
		TotalPayrolls:    res.TotalPayrolls,
		TotalBasicSalary: res.TotalBasicSalary,
		TotalNetSalary:   res.TotalNetSalary,
		PeriodStart:      startDate,
		PeriodEnd:        endDate,
	}

	return summary, nil
}

// Ensure payrollRepository implements PayrollRepository interface
var _ PayrollRepository = (*payrollRepository)(nil)
