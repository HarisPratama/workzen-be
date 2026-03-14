package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ManpowerReqRepository interface {
	GetManpowerReqsByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error)
	GetDetailManpowerRequestByTenant(ctx context.Context, TenantID int64, manpowerRequestID int64) (*entity.ManpowerReqEntity, error)
	CreateManpowerReq(ctx context.Context, manpowerReq entity.ManpowerReqEntity) error
}

type manpowerReqRepository struct {
	db *gorm.DB
}

func (m *manpowerReqRepository) GetDetailManpowerRequestByTenant(ctx context.Context, tenantID int64, manpowerRequestID int64) (*entity.ManpowerReqEntity, error) {
	var modelManpowerRequest model.ManpowerRequest
	var hiredCount int64

	err = m.db.
		Where("id = ?", manpowerRequestID).
		Where("tenant_id = ?", tenantID).
		Preload(clause.Associations).
		First(&modelManpowerRequest).Error

	if err != nil {
		code = "[REPOSITORY] GetDetailManpowerRequestByTenant - 1"
		log.Errorw(code, err)
		return nil, err
	}

	err = m.db.
		Model(&model.CandidateApplication{}).
		Where("tenant_id = ?", tenantID).
		Where("manpower_request_id = ?", manpowerRequestID).
		Where("status = ?", "HIRED").
		Count(&hiredCount).Error

	if err != nil {
		return nil, err
	}

	resp := entity.ManpowerReqEntity{
		ID:             modelManpowerRequest.ID,
		Position:       modelManpowerRequest.Position,
		RequiredCount:  int(modelManpowerRequest.RequiredCount),
		Hired:          int(hiredCount),
		SalaryMin:      modelManpowerRequest.SalaryMin,
		SalaryMax:      modelManpowerRequest.SalaryMax,
		WorkLocation:   modelManpowerRequest.WorkLocation,
		JobDescription: modelManpowerRequest.JobDescription,
		DeadlineDate:   modelManpowerRequest.DeadlineDate,
		CreatedAt:      modelManpowerRequest.CreatedAt,
		Status:         modelManpowerRequest.Status,
		Client: entity.ClientEntity{
			ID:          modelManpowerRequest.Client.ID,
			CompanyName: modelManpowerRequest.Client.CompanyName,
			Address:     modelManpowerRequest.Client.Address,
		},
	}

	return &resp, nil
}

func (m *manpowerReqRepository) GetManpowerReqsByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error) {
	var modelManpowerRequest []model.ManpowerRequest
	var countData int64

	sqlMain := m.db.WithContext(ctx).
		Model(&model.ManpowerRequest{}).
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
		Find(&modelManpowerRequest).Error; err != nil {
		code = "[REPOSITORY] GetManpowerReqsByTenant - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := make([]entity.ManpowerReqEntity, 0, len(modelManpowerRequest))
	for _, val := range modelManpowerRequest {
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
