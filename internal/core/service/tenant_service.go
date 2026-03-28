package service

import (
	"context"
	"errors"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

var (
	ErrEmailAlreadyExists = errors.New("email is already registered, please use a different email")
	ErrCreateTenantFailed = errors.New("failed to create company, please try again later")
	ErrHashPasswordFailed = errors.New("failed to process password, please try again later")
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

		if err.Error() == "email already exists" {
			return ErrEmailAlreadyExists
		}

		return err
	}

	return nil
}

func NewTenantService(tenantRepository repository.TenantRepository) TenantService {
	return &tenantService{tenantRepository: tenantRepository}
}
