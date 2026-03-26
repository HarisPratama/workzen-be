package service

import (
	"context"
	"fmt"
	"workzen-be/config"
	"workzen-be/internal/adapter/cloudflare"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type EmployeeService interface {
	GetEmployees(ctx context.Context, query entity.EmployeeQueryString, role string, tenantID int64) ([]entity.EmployeeEntity, int64, int64, error)
	GetDetailEmployeeByTenant(ctx context.Context, tenantID int64, employeeID int64) (*entity.EmployeeEntity, error)
	CreateEmployee(ctx context.Context, employee *entity.EmployeeEntity) error
	UpdateEmployee(ctx context.Context, tenantID int64, employeeID int64, employee *entity.EmployeeEntity) error
	DeleteEmployee(ctx context.Context, tenantID int64, employeeID int64) error
}

type employeeService struct {
	employeeRepo repository.EmployeeRepository
	cfg          *config.Config
	r2           cloudflare.CloudflareR2Adapter
}

func (e *employeeService) CreateEmployee(ctx context.Context, employee *entity.EmployeeEntity) error {
	err := e.employeeRepo.CreateEmployee(ctx, employee)
	if err != nil {
		code = "[SERVICE] CreateEmployee - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (e *employeeService) GetEmployees(ctx context.Context, query entity.EmployeeQueryString, role string, tenantID int64) ([]entity.EmployeeEntity, int64, int64, error) {
	if role == "SUPER_ADMIN" {
		return e.employeeRepo.GetEmployees(ctx, query)
	}

	return e.employeeRepo.GetEmployeesByTenant(ctx, tenantID, query)
}

func (e *employeeService) GetDetailEmployeeByTenant(ctx context.Context, tenantID int64, employeeID int64) (*entity.EmployeeEntity, error) {
	employee, err := e.employeeRepo.GetDetailEmployeeByTenant(ctx, tenantID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return employee, nil
}

func (e *employeeService) UpdateEmployee(ctx context.Context, tenantID int64, employeeID int64, employee *entity.EmployeeEntity) error {
	existing, err := e.employeeRepo.GetDetailEmployeeByTenant(ctx, tenantID, employeeID)
	if err != nil {
		return fmt.Errorf("employee not found: %w", err)
	}

	employee.ID = existing.ID
	employee.TenantID = existing.TenantID
	employee.UserID = existing.UserID
	employee.CreatedAt = existing.CreatedAt

	err = e.employeeRepo.UpdateEmployee(ctx, employee)
	if err != nil {
		code = "[SERVICE] UpdateEmployee - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (e *employeeService) DeleteEmployee(ctx context.Context, tenantID int64, employeeID int64) error {
	_, err := e.employeeRepo.GetDetailEmployeeByTenant(ctx, tenantID, employeeID)
	if err != nil {
		return fmt.Errorf("employee not found: %w", err)
	}

	err = e.employeeRepo.DeleteEmployee(ctx, tenantID, employeeID)
	if err != nil {
		code = "[SERVICE] DeleteEmployee - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewEmployeeService(repo repository.EmployeeRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) EmployeeService {
	return &employeeService{
		employeeRepo: repo,
		cfg:          cfg,
		r2:           r2,
	}
}
