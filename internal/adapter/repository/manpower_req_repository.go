package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"context"
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ManpowerReqRepository interface {
	GetManpowerReqsByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error)
	CreateManpowerReq(ctx context.Context, manpowerReq entity.ManpowerReqEntity) error
}

type manpowerReqRepository struct {
	db *gorm.DB
}

func (m *manpowerReqRepository) GetManpowerReqsByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error) {
	var modelManpowerReqs []model.ManpowerRequest
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := m.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ?", tenantID)

	if err := sqlMain.Model(&modelManpowerReqs).Count(&countData).Error; err != nil {
		code = "[REPOSITORY] GetManpowerReqsByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelManpowerReqs).Error; err != nil {
		code = "[REPOSITORY] GetManpowerReqsByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := make([]entity.ManpowerReqEntity, 0, len(modelManpowerReqs))
	for _, val := range modelManpowerReqs {
		resps = append(resps, entity.ManpowerReqEntity{
			ID:             val.ID,
			TenantID:       val.TenantID,
			ClientID:       val.ClientID,
			Position:       val.Position,
			RequiredCount:  int(val.RequiredCount),
			SalaryMin:      val.SalaryMin,
			SalaryMax:      val.SalaryMax,
			WorkLocation:   val.WorkLocation,
			JobDescription: val.JobDescription,
			DeadlineDate:   val.DeadlineDate,
			Status:         val.Status,
			Client: entity.ClientEntity{
				ID:          val.ClientID,
				CompanyName: val.Client.CompanyName,
				Address:     val.Client.Address,
			},
		})
	}

	return resps, countData, totalPages, nil
}

func (m *manpowerReqRepository) CreateManpowerReq(ctx context.Context, req entity.ManpowerReqEntity) error {
	modelManpowerReq := model.ManpowerRequest{
		TenantID:       req.TenantID,
		ClientID:       req.ClientID,
		Position:       req.Position,
		RequiredCount:  int64(req.RequiredCount),
		SalaryMin:      req.SalaryMin,
		SalaryMax:      req.SalaryMax,
		WorkLocation:   req.WorkLocation,
		JobDescription: req.JobDescription,
		DeadlineDate:   req.DeadlineDate,
		Status:         "OPEN",
	}

	err := m.db.Create(&modelManpowerReq).Error
	if err != nil {
		code = "[REPOSITORY] CreateManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewManpowerReqRepository(db *gorm.DB) ManpowerReqRepository {
	return &manpowerReqRepository{db: db}
}
