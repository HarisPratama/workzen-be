package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"bwanews/lib/conv"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TenantRepository interface {
	RegisterTenant(ctx context.Context, req entity.RegisterTenantEntity) error
}

type tenantRepository struct {
	db *gorm.DB
}

func (t *tenantRepository) RegisterTenant(ctx context.Context, req entity.RegisterTenantEntity) error {

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		modelTenant := model.Tenant{
			CompanyName: req.CompanyName,
			Address:     req.Address,
			Plan:        "FREE",
			Status:      "ACTIVE",
		}

		err := tx.Create(&modelTenant).Error
		if err != nil {
			code = "[REPOSITORY] RegisterTenant - 1"
			log.Errorw(code, err)
			return err
		}

		bytes, err := conv.HashPassword(req.Password)
		if err != nil {
			code = "[REPOSITORY] RegisterTenant - 2"
			log.Errorw(code, err)
			return err
		}

		adminUser := model.User{
			TenantID: &modelTenant.ID,
			Name:     req.Name,
			Email:    req.Email,
			Password: string(bytes),
			Role:     "TENANT_ADMIN",
			Status:   "ACTIVE",
		}

		//if err := tx.Create(&adminUser).Error; err != nil {
		//	return err
		//}

		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "email"},
			},
			DoNothing: true,
		}).
			Create(&adminUser)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("email already exists")
		}

		return nil
	})
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}
