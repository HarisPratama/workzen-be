package repository

import (
	"context"
	"fmt"
	"math"
	"strings"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ManpowerReqRepository interface {
	GetManpowerReqsByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error)
	GetDetailManpowerRequestByTenant(ctx context.Context, tenantID int64, manpowerRequestID int64) (*entity.ManpowerReqEntity, error)
	GetManpowerReqByPublicToken(ctx context.Context, token string) (*entity.ManpowerReqEntity, error)
	UpdatePublicToken(ctx context.Context, tenantID int64, manpowerReqID int64, token string) error
	CreateManpowerReq(ctx context.Context, manpowerReq *entity.ManpowerReqEntity) error
	UpdateManpowerReq(ctx context.Context, manpowerReq *entity.ManpowerReqEntity) error
	DeleteManpowerReq(ctx context.Context, tenantID int64, manpowerReqID int64) error
}

type manpowerReqRepository struct {
	db *gorm.DB
}

func (m *manpowerReqRepository) GetDetailManpowerRequestByTenant(ctx context.Context, tenantID int64, manpowerRequestID int64) (*entity.ManpowerReqEntity, error) {
	var manpowerRequest entity.ManpowerReqEntity
	var hiredCount int64

	err := m.db.WithContext(ctx).
		Where("id = ?", manpowerRequestID).
		Where("tenant_id = ?", tenantID).
		Preload(clause.Associations).
		First(&manpowerRequest).Error

	if err != nil {
		code = "[REPOSITORY] GetDetailManpowerRequestByTenant - 1"
		log.Errorw(code, err)
		return nil, err
	}

	err = m.db.WithContext(ctx).
		Model(&entity.CandidateApplicationEntity{}).
		Where("tenant_id = ?", tenantID).
		Where("manpower_request_id = ?", manpowerRequestID).
		Where("status = ?", "HIRED").
		Count(&hiredCount).Error

	if err != nil {
		return nil, err
	}

	manpowerRequest.Hired = int(hiredCount)

	return &manpowerRequest, nil
}

func (m *manpowerReqRepository) GetManpowerReqsByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error) {
	var manpowerRequests []entity.ManpowerReqEntity
	var countData int64

	sqlMain := m.db.WithContext(ctx).
		Model(&entity.ManpowerReqEntity{}).
		Joins("LEFT JOIN clients c ON c.id = manpower_requests.client_id").
		Where("manpower_requests.tenant_id = ?", tenantID)

	if query.Search != "" {
		search := "%" + query.Search + "%"
		sqlMain = sqlMain.Where(`
		(
			manpower_requests.position ILIKE ? OR
			manpower_requests.work_location ILIKE ? OR
			manpower_requests.job_description ILIKE ? OR
			c.company_name ILIKE ?
		)
		`, search, search, search, search)
	}

	if query.Status != "" {
		sqlMain = sqlMain.Where("manpower_requests.status = ?", query.Status)
	}

	countQuery := sqlMain.Session(&gorm.Session{})

	if err := countQuery.Count(&countData).Error; err != nil {
		return nil, 0, 0, err
	}

	allowedOrder := map[string]string{
		"created_at": "manpower_requests.created_at",
		"position":   "manpower_requests.position",
		"deadline":   "manpower_requests.deadline_date",
	}

	orderBy, ok := allowedOrder[query.OrderBy]
	if !ok {
		orderBy = "manpower_requests.created_at"
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
		Preload(clause.Associations).
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&manpowerRequests).Error; err != nil {
		code = "[REPOSITORY] GetManpowerReqsByTenant - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return manpowerRequests, countData, totalPages, nil
}

func (m *manpowerReqRepository) CreateManpowerReq(ctx context.Context, manpowerReq *entity.ManpowerReqEntity) error {
	err := m.db.WithContext(ctx).Create(manpowerReq).Error
	if err != nil {
		code = "[REPOSITORY] CreateManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (m *manpowerReqRepository) UpdateManpowerReq(ctx context.Context, manpowerReq *entity.ManpowerReqEntity) error {
	err := m.db.WithContext(ctx).Save(manpowerReq).Error
	if err != nil {
		code = "[REPOSITORY] UpdateManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (m *manpowerReqRepository) DeleteManpowerReq(ctx context.Context, tenantID int64, manpowerReqID int64) error {
	err := m.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&entity.ManpowerReqEntity{}, manpowerReqID).Error
	if err != nil {
		code = "[REPOSITORY] DeleteManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (m *manpowerReqRepository) GetManpowerReqByPublicToken(ctx context.Context, token string) (*entity.ManpowerReqEntity, error) {
	var manpowerRequest entity.ManpowerReqEntity
	err := m.db.WithContext(ctx).
		Where("public_token = ?", token).
		Preload("Client").
		First(&manpowerRequest).Error

	if err != nil {
		return nil, err
	}

	return &manpowerRequest, nil
}

func (m *manpowerReqRepository) UpdatePublicToken(ctx context.Context, tenantID int64, manpowerReqID int64, token string) error {
	err := m.db.WithContext(ctx).
		Model(&entity.ManpowerReqEntity{}).
		Where("id = ? AND tenant_id = ?", manpowerReqID, tenantID).
		Update("public_token", token).Error

	if err != nil {
		code = "[REPOSITORY] UpdatePublicToken - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewManpowerReqRepository(db *gorm.DB) ManpowerReqRepository {
	return &manpowerReqRepository{db: db}
}
