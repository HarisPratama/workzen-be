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

type PayrollService interface {
	CreatePayroll(ctx context.Context, req entity.Payroll) (*entity.Payroll, error)
	UpdatePayroll(ctx context.Context, id uuid.UUID, req entity.Payroll) (*entity.Payroll, error)
	DeletePayroll(ctx context.Context, id uuid.UUID) error
	GetPayrollByID(ctx context.Context, id uuid.UUID) (*entity.Payroll, error)
	GetPayrollsByTenant(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error)
	GetPayrollsByEmployee(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error)
	GetPayrollsByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Payroll, int64, error)
	ProcessPayroll(ctx context.Context, id uuid.UUID) error
	MarkAsPaid(ctx context.Context, id uuid.UUID, paidAt time.Time) error
	AddPayrollItem(ctx context.Context, payrollID uuid.UUID, item entity.PayrollItem) error
	RemovePayrollItem(ctx context.Context, itemID uuid.UUID) error
	CalculatePayrollSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*entity.PayrollSummary, error)
}

type payrollService struct {
	payrollRepo repository.PayrollRepository
}

func NewPayrollService(payrollRepo repository.PayrollRepository) PayrollService {
	return &payrollService{
		payrollRepo: payrollRepo,
	}
}

func (s *payrollService) CreatePayroll(ctx context.Context, req entity.Payroll) (*entity.Payroll, error) {
	// Validate required fields
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Calculate net salary
	req.CalculateNetSalary()

	// Set default status
	if req.Status == "" {
		req.Status = entity.PayrollStatusDraft
	}

	if err := s.payrollRepo.Create(ctx, &req); err != nil {
		log.Errorw("failed to create payroll", "error", err)
		return nil, fmt.Errorf("failed to create payroll: %w", err)
	}

	return &req, nil
}

func (s *payrollService) UpdatePayroll(ctx context.Context, id uuid.UUID, req entity.Payroll) (*entity.Payroll, error) {
	// Check if payroll exists
	existing, err := s.payrollRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("payroll not found: %w", err)
	}

	// Check if payroll can be edited
	if !existing.IsEditable() {
		return nil, fmt.Errorf("payroll cannot be edited in status: %s", existing.Status)
	}

	// Update fields
	existing.BasicSalary = req.BasicSalary
	existing.Allowances = req.Allowances
	existing.Deductions = req.Deductions
	existing.Tax = req.Tax
	existing.Notes = req.Notes

	// Recalculate net salary
	existing.CalculateNetSalary()

	if err := s.payrollRepo.Update(ctx, existing); err != nil {
		log.Errorw("failed to update payroll", "error", err)
		return nil, fmt.Errorf("failed to update payroll: %w", err)
	}

	return existing, nil
}

func (s *payrollService) DeletePayroll(ctx context.Context, id uuid.UUID) error {
	// Check if payroll exists
	existing, err := s.payrollRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("payroll not found: %w", err)
	}

	// Only allow deletion of draft payrolls
	if existing.Status != entity.PayrollStatusDraft {
		return fmt.Errorf("only draft payrolls can be deleted")
	}

	if err := s.payrollRepo.Delete(ctx, id); err != nil {
		log.Errorw("failed to delete payroll", "error", err)
		return fmt.Errorf("failed to delete payroll: %w", err)
	}

	return nil
}

func (s *payrollService) GetPayrollByID(ctx context.Context, id uuid.UUID) (*entity.Payroll, error) {
	payroll, err := s.payrollRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("payroll not found: %w", err)
	}
	return payroll, nil
}

func (s *payrollService) GetPayrollsByTenant(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error) {
	return s.payrollRepo.FindByTenantID(ctx, tenantID, page, limit)
}

func (s *payrollService) GetPayrollsByEmployee(ctx context.Context, employeeID uuid.UUID, page, limit int) ([]entity.Payroll, int64, error) {
	return s.payrollRepo.FindByEmployeeID(ctx, employeeID, page, limit)
}

func (s *payrollService) GetPayrollsByPeriod(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, page, limit int) ([]entity.Payroll, int64, error) {
	return s.payrollRepo.FindByPeriod(ctx, tenantID, startDate, endDate, page, limit)
}

func (s *payrollService) ProcessPayroll(ctx context.Context, id uuid.UUID) error {
	payroll, err := s.payrollRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("payroll not found: %w", err)
	}

	if !payroll.CanProcess() {
		return fmt.Errorf("payroll cannot be processed in status: %s", payroll.Status)
	}

	if err := s.payrollRepo.ProcessPayroll(ctx, id); err != nil {
		log.Errorw("failed to process payroll", "error", err)
		return fmt.Errorf("failed to process payroll: %w", err)
	}

	return nil
}

func (s *payrollService) MarkAsPaid(ctx context.Context, id uuid.UUID, paidAt time.Time) error {
	payroll, err := s.payrollRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("payroll not found: %w", err)
	}

	if !payroll.CanPay() {
		return fmt.Errorf("payroll cannot be marked as paid in status: %s", payroll.Status)
	}

	if err := s.payrollRepo.MarkAsPaid(ctx, id, paidAt); err != nil {
		log.Errorw("failed to mark payroll as paid", "error", err)
		return fmt.Errorf("failed to mark payroll as paid: %w", err)
	}

	return nil
}

func (s *payrollService) AddPayrollItem(ctx context.Context, payrollID uuid.UUID, item entity.PayrollItem) error {
	payroll, err := s.payrollRepo.FindByID(ctx, payrollID)
	if err != nil {
		return fmt.Errorf("payroll not found: %w", err)
	}

	if !payroll.IsEditable() {
		return fmt.Errorf("cannot add items to payroll in status: %s", payroll.Status)
	}

	item.PayrollID = payrollID

	if err := s.payrollRepo.AddPayrollItem(ctx, &item); err != nil {
		log.Errorw("failed to add payroll item", "error", err)
		return fmt.Errorf("failed to add payroll item: %w", err)
	}

	return nil
}

func (s *payrollService) RemovePayrollItem(ctx context.Context, itemID uuid.UUID) error {
	if err := s.payrollRepo.DeletePayrollItem(ctx, itemID); err != nil {
		log.Errorw("failed to remove payroll item", "error", err)
		return fmt.Errorf("failed to remove payroll item: %w", err)
	}
	return nil
}

func (s *payrollService) CalculatePayrollSummary(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) (*entity.PayrollSummary, error) {
	summary, err := s.payrollRepo.CalculatePayrollSummary(ctx, tenantID, startDate, endDate)
	if err != nil {
		log.Errorw("failed to calculate payroll summary", "error", err)
		return nil, fmt.Errorf("failed to calculate payroll summary: %w", err)
	}
	return summary, nil
}
