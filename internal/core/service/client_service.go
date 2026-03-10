package service

import (
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type ClientService interface {
	GetClientByTenant(ctx context.Context, query entity.ClientQueryString, tenantID int64) ([]entity.ClientEntity, int64, int64, error)
	CreateClient(ctx context.Context, req entity.ClientEntity) error
}

type clientService struct {
	clientRepo repository.ClientRepository
}

func (c *clientService) GetClientByTenant(ctx context.Context, query entity.ClientQueryString, tenantID int64) ([]entity.ClientEntity, int64, int64, error) {
	results, totalData, totalPage, err := c.clientRepo.GetClientsByTenant(ctx, tenantID, query)
	if err != nil {
		code = "[SERVICE] GetClientByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPage, nil
}

func (c *clientService) CreateClient(ctx context.Context, req entity.ClientEntity) error {
	err := c.clientRepo.CreateClient(ctx, req)

	if err != nil {
		code = "[SERVICE] CreateClient - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientService{
		clientRepo: repo,
	}
}
