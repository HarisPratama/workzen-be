package repository

import (
	"context"
	"fmt"
	"math"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClientRepository interface {
	GetClientsByTenant(ctx context.Context, tenantID int64, query entity.ClientQueryString) ([]entity.ClientEntity, int64, int64, error)
	GetDetailClientByTenant(ctx context.Context, tenantID int64, clientID int64) (*entity.ClientEntity, error)
	CreateClient(ctx context.Context, client *entity.ClientEntity) error
	UpdateClient(ctx context.Context, client *entity.ClientEntity) error
	DeleteClient(ctx context.Context, tenantID int64, clientID int64) error
}

type clientRepository struct {
	db *gorm.DB
}

func (c *clientRepository) GetDetailClientByTenant(ctx context.Context, tenantID int64, clientID int64) (*entity.ClientEntity, error) {
	var client entity.ClientEntity
	err := c.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&client, clientID).Error

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *clientRepository) CreateClient(ctx context.Context, client *entity.ClientEntity) error {
	err := c.db.WithContext(ctx).Create(client).Error
	if err != nil {
		code = "[REPOSITORY] CreateClient - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *clientRepository) UpdateClient(ctx context.Context, client *entity.ClientEntity) error {
	err := c.db.WithContext(ctx).Save(client).Error
	if err != nil {
		code = "[REPOSITORY] UpdateClient - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *clientRepository) DeleteClient(ctx context.Context, tenantID int64, clientID int64) error {
	err := c.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&entity.ClientEntity{}, clientID).Error
	if err != nil {
		code = "[REPOSITORY] DeleteClient - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *clientRepository) GetClientsByTenant(ctx context.Context, tenantID int64, query entity.ClientQueryString) ([]entity.ClientEntity, int64, int64, error) {
	var clients []entity.ClientEntity
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := c.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ?", tenantID)

	if query.Search != "" {
		search := "%" + query.Search + "%"
		sqlMain = sqlMain.Where("company_name ILIKE ?", search)
	}

	if err := sqlMain.Model(&entity.ClientEntity{}).Count(&countData).Error; err != nil {
		code = "[REPOSITORY] GetClientByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&clients).Error; err != nil {
		code = "[REPOSITORY] GetClientByTenant - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return clients, countData, totalPages, nil
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}
