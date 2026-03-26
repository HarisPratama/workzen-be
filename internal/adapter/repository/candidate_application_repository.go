package repository

import (
	"context"
	"fmt"
	"math"
	"strings"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CandidateApplicationRepository interface {
	GetCandidateApplicationsByTenant(ctx context.Context, tenantID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error)
	GetCandidateApplicationByTenantMR(ctx context.Context, tenantID int64, manpowerRequestID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error)
	GetDetailCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) (*entity.CandidateApplicationEntity, error)
	CreateCandidateApplication(ctx context.Context, application *entity.CandidateApplicationEntity) error
	UpdateCandidateApplication(ctx context.Context, application *entity.CandidateApplicationEntity) error
	DeleteCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) error
}

type candidateApplicationRepository struct {
	db *gorm.DB
}

func (c *candidateApplicationRepository) GetDetailCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) (*entity.CandidateApplicationEntity, error) {
	var application entity.CandidateApplicationEntity
	err := c.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Preload("Candidate").
		Preload("ManpowerRequest").
		First(&application, applicationID).Error

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (c *candidateApplicationRepository) CreateCandidateApplication(ctx context.Context, application *entity.CandidateApplicationEntity) error {
	err := c.db.WithContext(ctx).Create(application).Error
	if err != nil {
		code = "[REPOSITORY] CreateCandidateApplication - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateApplicationRepository) UpdateCandidateApplication(ctx context.Context, application *entity.CandidateApplicationEntity) error {
	err := c.db.WithContext(ctx).Save(application).Error
	if err != nil {
		code = "[REPOSITORY] UpdateCandidateApplication - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateApplicationRepository) DeleteCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) error {
	err := c.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&entity.CandidateApplicationEntity{}, applicationID).Error
	if err != nil {
		code = "[REPOSITORY] DeleteCandidateApplication - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateApplicationRepository) GetCandidateApplicationByTenantMR(ctx context.Context, tenantID int64, manpowerRequestID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error) {
	var applications []entity.CandidateApplicationEntity
	var countData int64

	sqlMain := c.db.WithContext(ctx).
		Model(&entity.CandidateApplicationEntity{}).
		Joins("LEFT JOIN candidates c ON c.id = candidate_applications.candidate_id").
		Where("candidate_applications.tenant_id = ?", tenantID).
		Where("candidate_applications.manpower_request_id = ?", manpowerRequestID)

	if query.Search != "" {
		search := "%" + query.Search + "%"
		sqlMain = sqlMain.Where(`c.full_name ILIKE ?`, search)
	}

	if query.Status != "" {
		sqlMain = sqlMain.Where("candidate_applications.status = ?", query.Status)
	}

	countQuery := sqlMain.Session(&gorm.Session{})

	if err := countQuery.Count(&countData).Error; err != nil {
		return nil, 0, 0, err
	}

	allowedOrder := map[string]string{
		"applied_at": "candidate_applications.applied_at",
	}

	orderBy, ok := allowedOrder[query.OrderBy]
	if !ok {
		orderBy = "candidate_applications.applied_at"
	}

	orderType := "DESC"
	if strings.ToUpper(query.OrderType) == "ASC" {
		orderType = "ASC"
	}

	order := fmt.Sprintf("%s %s", orderBy, orderType)

	if query.Limit <= 0 {
		query.Limit = 10
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	offset := (query.Page - 1) * query.Limit

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Preload("Candidate").
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&applications).Error; err != nil {
		code = "[REPOSITORY] GetCandidateApplicationByTenantMR - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return applications, countData, totalPages, nil
}

func (c *candidateApplicationRepository) GetCandidateApplicationsByTenant(ctx context.Context, tenantID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error) {
	var applications []entity.CandidateApplicationEntity
	var countData int64

	sqlMain := c.db.WithContext(ctx).
		Model(&entity.CandidateApplicationEntity{}).
		Where("candidate_applications.tenant_id = ?", tenantID)

	if query.Search != "" {
		search := "%" + query.Search + "%"
		sqlMain = sqlMain.
			Joins("LEFT JOIN candidates c ON c.id = candidate_applications.candidate_id").
			Where(`c.full_name ILIKE ?`, search)
	}

	if query.Status != "" {
		sqlMain = sqlMain.Where("candidate_applications.status = ?", query.Status)
	}

	countQuery := sqlMain.Session(&gorm.Session{})
	if err := countQuery.Count(&countData).Error; err != nil {
		return nil, 0, 0, err
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	offset := (query.Page - 1) * query.Limit
	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	orderBy := "candidate_applications.applied_at"
	orderType := "DESC"
	if strings.ToUpper(query.OrderType) == "ASC" {
		orderType = "ASC"
	}
	order := fmt.Sprintf("%s %s", orderBy, orderType)

	if err := sqlMain.
		Preload("Candidate").
		Preload("ManpowerRequest").
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&applications).Error; err != nil {
		code = "[REPOSITORY] GetCandidateApplicationsByTenant - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return applications, countData, totalPages, nil
}

func NewCandidateApplicationRepository(db *gorm.DB) CandidateApplicationRepository {
	return &candidateApplicationRepository{db: db}
}
