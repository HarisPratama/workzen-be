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
	GetEmployees(ctx context.Context, query entity.EmployeeQueryString) ([]entity.EmployeeEntity, int64, int64, error)
	CreateEmployee(ctx context.Context, req entity.EmployeeEntity, tenantId int64) error
}

type employeeService struct {
	employeeRepo repository.EmployeeRepository
	cfg          *config.Config
	r2           cloudflare.CloudflareR2Adapter
}

func (e *employeeService) CreateEmployee(ctx context.Context, req entity.EmployeeEntity, tenantId int64) error {
	err := e.employeeRepo.CreateEmployee(ctx, req, tenantId)

	if err != nil {
		code = "[SERVICE] CreateEmployee - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (e *employeeService) GetEmployees(ctx context.Context, query entity.EmployeeQueryString) ([]entity.EmployeeEntity, int64, int64, error) {
	results, totalData, totalPage, err := e.employeeRepo.GetEmployees(ctx, query)

	if err != nil {
		code = "[SERVICE] GetEmployees - 1"
		log.Error(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPage, nil
}

func NewEmployeeService(repo repository.EmployeeRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) EmployeeService {
	return &employeeService{
		employeeRepo: repo,
		cfg:          cfg,
		r2:           r2,
	}
}
