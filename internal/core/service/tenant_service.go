package service

import (
	"context"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type TenantService interface {
	RegisterTenant(ctx context.Context, req entity.RegisterTenantEntity) error
}

type tenantService struct {
	tenantRepository repository.TenantRepository
}

func (t *tenantService) RegisterTenant(ctx context.Context, req entity.RegisterTenantEntity) error {
	err := t.tenantRepository.RegisterTenant(ctx, req)

	if err != nil {
		code = "[SERVICE] RegisterTenant - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewTenantService(tenantRepository repository.TenantRepository) TenantService {
	return &tenantService{tenantRepository: tenantRepository}
}
