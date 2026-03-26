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

type CandidateRepository interface {
	GetCandidatesByTenant(ctx context.Context, tenantID int64, query entity.CandidateQueryString) ([]entity.CandidateEntity, int64, int64, error)
	CreateCandidate(ctx context.Context, req entity.CandidateEntity) error
}

type candidateRepository struct {
	db *gorm.DB
}

func (c *candidateRepository) GetCandidatesByTenant(ctx context.Context, tenantID int64, query entity.CandidateQueryString) ([]entity.CandidateEntity, int64, int64, error) {
	var modelCandidate []model.Candidate
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := c.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ?", tenantID)

	if err := sqlMain.Model(&model.Candidate{}).Count(&countData).Error; err != nil {
		code = "[REPOSITORY] GetCandidateByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelCandidate).Error; err != nil {
		code = "[REPOSITORY] GetCandidateByTenant - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := make([]entity.CandidateEntity, 0, len(modelCandidate))
	for _, val := range modelCandidate {
		resps = append(resps, entity.CandidateEntity{
			ID:        val.ID,
			FullName:  val.FullName,
			Email:     val.Email,
			Phone:     val.Phone,
			Address:   val.Address,
			Source:    val.Source,
			Status:    val.Status,
			BirthDate: val.BirthDate,
		})
	}

	return resps, countData, totalPages, nil
}

func (c *candidateRepository) CreateCandidate(ctx context.Context, req entity.CandidateEntity) error {
	modelCandidate := model.Candidate{
		TenantID:  req.TenantID,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CitizenID: req.CitizenID,
		Source:    req.Source,
		Status:    "ACTIVE",
		BirthDate: req.BirthDate,
	}

	err := c.db.Create(&modelCandidate).Error
	if err != nil {
		code = "[REPOSITORY] CreateCandidate - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &candidateRepository{db: db}
}
