package service

import (
	"context"
	"fmt"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type ClientService interface {
	GetClientByTenant(ctx context.Context, query entity.ClientQueryString, tenantID int64) ([]entity.ClientEntity, int64, int64, error)
	GetDetailClientByTenant(ctx context.Context, tenantID int64, clientID int64) (*entity.ClientEntity, error)
	CreateClient(ctx context.Context, client *entity.ClientEntity) error
	UpdateClient(ctx context.Context, tenantID int64, clientID int64, client *entity.ClientEntity) error
	DeleteClient(ctx context.Context, tenantID int64, clientID int64) error
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

func (c *clientService) GetDetailClientByTenant(ctx context.Context, tenantID int64, clientID int64) (*entity.ClientEntity, error) {
	client, err := c.clientRepo.GetDetailClientByTenant(ctx, tenantID, clientID)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	return client, nil
}

func (c *clientService) CreateClient(ctx context.Context, client *entity.ClientEntity) error {
	err := c.clientRepo.CreateClient(ctx, client)
	if err != nil {
		code = "[SERVICE] CreateClient - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *clientService) UpdateClient(ctx context.Context, tenantID int64, clientID int64, client *entity.ClientEntity) error {
	existing, err := c.clientRepo.GetDetailClientByTenant(ctx, tenantID, clientID)
	if err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	client.ID = existing.ID
	client.TenantID = existing.TenantID
	client.CreatedAt = existing.CreatedAt

	err = c.clientRepo.UpdateClient(ctx, client)
	if err != nil {
		code = "[SERVICE] UpdateClient - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *clientService) DeleteClient(ctx context.Context, tenantID int64, clientID int64) error {
	_, err := c.clientRepo.GetDetailClientByTenant(ctx, tenantID, clientID)
	if err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	err = c.clientRepo.DeleteClient(ctx, tenantID, clientID)
	if err != nil {
		code = "[SERVICE] DeleteClient - 1"
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
