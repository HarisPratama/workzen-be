package service

import (
	"bwanews/config"
	"bwanews/internal/adapter/cloudflare"
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type EmployeeService interface {
	GetEmployees(ctx context.Context, query entity.EmployeeQueryString, role string, tenantID int64) ([]entity.EmployeeEntity, int64, int64, error)
	CreateEmployee(ctx context.Context, req entity.EmployeeEntity) error
}

type employeeService struct {
	employeeRepo repository.EmployeeRepository
	cfg          *config.Config
	r2           cloudflare.CloudflareR2Adapter
}

func (e *employeeService) CreateEmployee(ctx context.Context, req entity.EmployeeEntity) error {
	err := e.employeeRepo.CreateEmployee(ctx, req)

	if err != nil {
		code = "[SERVICE] CreateEmployee - 2"
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

func NewEmployeeService(repo repository.EmployeeRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) EmployeeService {
	return &employeeService{
		employeeRepo: repo,
		cfg:          cfg,
		r2:           r2,
	}
}
