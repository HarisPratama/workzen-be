package repository

import (
	"context"
	"fmt"
	"math"
	"strings"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CandidateApplicationRepository interface {
	GetCandidateApplicationsByTenant(ctx context.Context, tenantID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error)
	GetCandidateApplicationByTenantMR(ctx context.Context, tenantID int64, manpowerRequestID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error)
	CreateCandidateApplication(ctx context.Context, req entity.CandidateApplicationEntity) error
}

type candidateApplicationRepository struct {
	db *gorm.DB
}

func (c *candidateApplicationRepository) CreateCandidateApplication(ctx context.Context, req entity.CandidateApplicationEntity) error {
	modelCandidateApplication := model.CandidateApplication{
		TenantID:          req.TenantID,
		ManpowerRequestID: req.ManpowerRequestID,
		CandidateID:       req.CandidateID,
		Status:            "APPLIED",
	}

	err := c.db.Create(&modelCandidateApplication).Error
	if err != nil {
		code = "[REPOSITORY] CreateCandidateApplication - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateApplicationRepository) GetCandidateApplicationByTenantMR(ctx context.Context, tenantID int64, manpowerRequestID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error) {
	var modelCandidateApplications []model.CandidateApplication
	var countData int64

	sqlMain := c.db.WithContext(ctx).
		Model(&model.CandidateApplication{}).
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
		Find(&modelCandidateApplications).Error; err != nil {
		code = "[REPOSITORY] GetCandidateApplicationByTenantMR - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := make([]entity.CandidateApplicationEntity, 0, len(modelCandidateApplications))
	for _, val := range modelCandidateApplications {
		resps = append(resps, entity.CandidateApplicationEntity{
			ID:       val.ID,
			TenantID: val.TenantID,
			Candidate: entity.CandidateEntity{
				ID:        val.CandidateID,
				FullName:  val.Candidate.FullName,
				Email:     val.Candidate.Email,
				Phone:     val.Candidate.Phone,
				BirthDate: val.Candidate.BirthDate,
				Address:   val.Candidate.Address,
				Source:    val.Candidate.Source,
			},
			ManpowerRequest: entity.ManpowerReqEntity{
				ID: val.ManpowerRequestID,
			},
		})
	}

	return resps, countData, totalPages, nil
}

func (c *candidateApplicationRepository) GetCandidateApplicationsByTenant(ctx context.Context, tenantID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewCandidateApplicationRepository(db *gorm.DB) CandidateApplicationRepository {
	return &candidateApplicationRepository{db: db}
}
