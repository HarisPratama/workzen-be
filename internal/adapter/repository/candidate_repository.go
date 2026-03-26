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

type CandidateRepository interface {
	GetCandidatesByTenant(ctx context.Context, tenantID int64, query entity.CandidateQueryString) ([]entity.CandidateEntity, int64, int64, error)
	GetDetailCandidateByTenant(ctx context.Context, tenantID int64, candidateID int64) (*entity.CandidateEntity, error)
	CreateCandidate(ctx context.Context, candidate *entity.CandidateEntity) error
	UpdateCandidate(ctx context.Context, candidate *entity.CandidateEntity) error
	DeleteCandidate(ctx context.Context, tenantID int64, candidateID int64) error
}

type candidateRepository struct {
	db *gorm.DB
}

func (c *candidateRepository) GetDetailCandidateByTenant(ctx context.Context, tenantID int64, candidateID int64) (*entity.CandidateEntity, error) {
	var candidate entity.CandidateEntity
	err := c.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&candidate, candidateID).Error

	if err != nil {
		return nil, err
	}

	return &candidate, nil
}

func (c *candidateRepository) GetCandidatesByTenant(ctx context.Context, tenantID int64, query entity.CandidateQueryString) ([]entity.CandidateEntity, int64, int64, error) {
	var candidates []entity.CandidateEntity
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := c.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ?", tenantID)

	if query.Search != "" {
		search := "%" + query.Search + "%"
		sqlMain = sqlMain.Where("full_name ILIKE ? OR email ILIKE ?", search, search)
	}

	if query.Status != "" {
		sqlMain = sqlMain.Where("status = ?", query.Status)
	}

	if err := sqlMain.Model(&entity.CandidateEntity{}).Count(&countData).Error; err != nil {
		code = "[REPOSITORY] GetCandidateByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&candidates).Error; err != nil {
		code = "[REPOSITORY] GetCandidateByTenant - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return candidates, countData, totalPages, nil
}

func (c *candidateRepository) CreateCandidate(ctx context.Context, candidate *entity.CandidateEntity) error {
	err := c.db.WithContext(ctx).Create(candidate).Error
	if err != nil {
		code = "[REPOSITORY] CreateCandidate - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateRepository) UpdateCandidate(ctx context.Context, candidate *entity.CandidateEntity) error {
	err := c.db.WithContext(ctx).Save(candidate).Error
	if err != nil {
		code = "[REPOSITORY] UpdateCandidate - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateRepository) DeleteCandidate(ctx context.Context, tenantID int64, candidateID int64) error {
	err := c.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&entity.CandidateEntity{}, candidateID).Error
	if err != nil {
		code = "[REPOSITORY] DeleteCandidate - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &candidateRepository{db: db}
}
