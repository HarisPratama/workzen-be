package repository

import (
	"context"
	"fmt"
	"math"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClientRepository interface {
	GetClientsByTenant(ctx context.Context, tenantID int64, query entity.ClientQueryString) ([]entity.ClientEntity, int64, int64, error)
	CreateClient(ctx context.Context, client entity.ClientEntity) error
}

type clientRepository struct {
	db *gorm.DB
}

func (c *clientRepository) CreateClient(ctx context.Context, req entity.ClientEntity) error {
	modelClient := model.Client{
		CompanyName: req.CompanyName,
		Address:     req.Address,
		TenantID:    req.TenantID,
	}

	err := c.db.Create(&modelClient).Error
	if err != nil {
		code = "[REPOSITORY] CreateClient - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *clientRepository) GetClientsByTenant(ctx context.Context, tenantID int64, query entity.ClientQueryString) ([]entity.ClientEntity, int64, int64, error) {
	var modelClients []model.Client
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := c.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ?", tenantID)

	if err := sqlMain.Model(&modelClients).Count(&countData).Error; err != nil {
		code = "[REPOSITORY] GetClientByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelClients).Error; err != nil {
		code = "[REPOSITORY] GetClientByTenant - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := make([]entity.ClientEntity, 0, len(modelClients))
	for _, val := range modelClients {
		resps = append(resps, entity.ClientEntity{
			ID:          val.ID,
			CompanyName: val.CompanyName,
			Address:     val.Address,
			CreatedAt:   val.CreatedAt,
		})
	}

	return resps, countData, totalPages, nil
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}
