package repository

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type OfferRepository interface {
	GetOffersByTenant(ctx context.Context, tenantID int64, query entity.OfferQueryString) ([]entity.OfferEntity, int64, int64, error)
	GetOfferByID(ctx context.Context, id int64) (*entity.OfferEntity, error)
	CreateOffer(ctx context.Context, req entity.OfferEntity, tenantID int64) error
	UpdateOffer(ctx context.Context, id int64, req entity.OfferUpdateRequest) error
	DeleteOffer(ctx context.Context, id int64) error
}

type offerRepository struct {
	db *gorm.DB
}

func (r *offerRepository) CreateOffer(ctx context.Context, req entity.OfferEntity, tenantID int64) error {
	modelOffer := model.Offer{
		TenantID:               tenantID,
		CandidateApplicationID: req.CandidateApplicationID,
		Position:               req.Position,
		Department:             req.Department,
		EmploymentType:         req.EmploymentType,
		BaseSalary:             req.BaseSalary,
		Currency:               req.Currency,
		Bonus:                  req.Bonus,
		Benefits:               req.Benefits,
		ProbationPeriodMonths:  req.ProbationPeriodMonths,
		NoticePeriodDays:       req.NoticePeriodDays,
		StartDate:              req.StartDate,
		ExpiryDate:             req.ExpiryDate,
		Status:                 "SENT",
	}

	err := r.db.Create(&modelOffer).Error
	if err != nil {
		code := "[REPOSITORY] CreateOffer - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (r *offerRepository) GetOffersByTenant(ctx context.Context, tenantID int64, query entity.OfferQueryString) ([]entity.OfferEntity, int64, int64, error) {
	var modelOffers []model.Offer
	var countData int64

	sqlMain := r.db.WithContext(ctx).
		Model(&model.Offer{}).
		Preload("CandidateApplication").
		Preload("CandidateApplication.Candidate").
		Where("tenant_id = ?", tenantID)

	if query.Search != "" {
		search := "%" + query.Search + "%"
		sqlMain = sqlMain.Where(`candidate_applications.position ILIKE ?`, search)
	}

	if query.Status != "" {
		sqlMain = sqlMain.Where("offers.status = ?", query.Status)
	}

	if query.CandidateApplicationID != 0 {
		sqlMain = sqlMain.Where("offers.candidate_application_id = ?", query.CandidateApplicationID)
	}

	if query.CandidateID != 0 {
		sqlMain = sqlMain.Joins("JOIN candidate_applications ca ON ca.id = offers.candidate_application_id").
			Where("ca.candidate_id = ?", query.CandidateID)
	}

	if query.StartDate != "" && query.EndDate != "" {
		sqlMain = sqlMain.Where("offers.created_at BETWEEN ? AND ?", query.StartDate, query.EndDate)
	}

	countQuery := sqlMain.Session(&gorm.Session{})

	if err := countQuery.Count(&countData).Error; err != nil {
		return nil, 0, 0, err
	}

	allowedOrder := map[string]string{
		"created_at": "offers.created_at",
		"updated_at": "offers.updated_at",
	}

	orderBy, ok := allowedOrder[query.OrderBy]
	if !ok {
		orderBy = "offers.created_at"
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
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelOffers).Error; err != nil {
		return nil, 0, 0, err
	}

	var result []entity.OfferEntity
	for _, item := range modelOffers {
		result = append(result, entity.OfferEntity{
			ID:                     item.ID,
			TenantID:               item.TenantID,
			CandidateApplicationID: item.CandidateApplicationID,
			Position:               item.Position,
			Department:             item.Department,
			EmploymentType:         item.EmploymentType,
			BaseSalary:             item.BaseSalary,
			Currency:               item.Currency,
			Bonus:                  item.Bonus,
			Benefits:               item.Benefits,
			ProbationPeriodMonths:  item.ProbationPeriodMonths,
			NoticePeriodDays:       item.NoticePeriodDays,
			StartDate:              item.StartDate,
			ExpiryDate:             item.ExpiryDate,
			Status:                 item.Status,
			SentAt:                 item.SentAt,
			RespondedAt:            item.RespondedAt,
			Notes:                  item.Notes,
			Terms:                  item.Terms,
			Feedback:               item.Feedback,
			NegotiationCounter:     item.NegotiationCounter,
			NegotiationNotes:       item.NegotiationNotes,
			Tenant: entity.TenantEntity{
				ID:          item.Tenant.ID,
				CompanyName: item.Tenant.CompanyName,
				Plan:        item.Tenant.Plan,
				Status:      item.Tenant.Status,
				Address:     item.Tenant.Address,
			},
			CandidateApplication: entity.CandidateApplicationEntity{
				ID:          item.CandidateApplication.ID,
				TenantID:    item.CandidateApplication.TenantID,
				CandidateID: item.CandidateApplication.CandidateID,
				Status:      item.CandidateApplication.Status,
				Candidate: entity.CandidateEntity{
					ID:       item.CandidateApplication.Candidate.ID,
					FullName: item.CandidateApplication.Candidate.FullName,
					Email:    item.CandidateApplication.Candidate.Email,
				},
			},
		})
	}

	return result, countData, totalPages, nil
}

func (r *offerRepository) GetOfferByID(ctx context.Context, id int64) (*entity.OfferEntity, error) {
	var modelOffer model.Offer
	err := r.db.WithContext(ctx).
		Preload("CandidateApplication").
		Preload("CandidateApplication.Candidate").
		First(&modelOffer, id).Error

	if err != nil {
		return nil, err
	}

	return &entity.OfferEntity{
		ID:                     modelOffer.ID,
		TenantID:               modelOffer.TenantID,
		CandidateApplicationID: modelOffer.CandidateApplicationID,
		Position:               modelOffer.Position,
		Department:             modelOffer.Department,
		EmploymentType:         modelOffer.EmploymentType,
		BaseSalary:             modelOffer.BaseSalary,
		Currency:               modelOffer.Currency,
		Bonus:                  modelOffer.Bonus,
		Benefits:               modelOffer.Benefits,
		ProbationPeriodMonths:  modelOffer.ProbationPeriodMonths,
		NoticePeriodDays:       modelOffer.NoticePeriodDays,
		StartDate:              modelOffer.StartDate,
		ExpiryDate:             modelOffer.ExpiryDate,
		Status:                 modelOffer.Status,
		SentAt:                 modelOffer.SentAt,
		RespondedAt:            modelOffer.RespondedAt,
		Notes:                  modelOffer.Notes,
		Terms:                  modelOffer.Terms,
		Feedback:               modelOffer.Feedback,
		NegotiationCounter:     modelOffer.NegotiationCounter,
		NegotiationNotes:       modelOffer.NegotiationNotes,
		Tenant: entity.TenantEntity{
			ID:          modelOffer.Tenant.ID,
			CompanyName: modelOffer.Tenant.CompanyName,
			Plan:        modelOffer.Tenant.Plan,
			Status:      modelOffer.Tenant.Status,
			Address:     modelOffer.Tenant.Address,
		},
		CandidateApplication: entity.CandidateApplicationEntity{
			ID:          modelOffer.CandidateApplication.ID,
			TenantID:    modelOffer.CandidateApplication.TenantID,
			CandidateID: modelOffer.CandidateApplication.CandidateID,
			Status:      modelOffer.CandidateApplication.Status,
			Candidate: entity.CandidateEntity{
				ID:       modelOffer.CandidateApplication.Candidate.ID,
				FullName: modelOffer.CandidateApplication.Candidate.FullName,
				Email:    modelOffer.CandidateApplication.Candidate.Email,
			},
		},
	}, nil
}

func (r *offerRepository) UpdateOffer(ctx context.Context, id int64, req entity.OfferUpdateRequest) error {
	updates := map[string]interface{}{}

	if req.Status != "" {
		status := strings.ToUpper(req.Status)
		updates["status"] = status

		if status == "SENT" {
			updates["sent_at"] = time.Now()
		} else if status == "ACCEPTED" || status == "REJECTED" {
			updates["responded_at"] = time.Now()
		}
	}

	if req.Feedback != "" {
		updates["feedback"] = req.Feedback
	}

	if req.NegotiatedSalary != nil {
		updates["negotiation_counter"] = req.NegotiatedSalary
	}

	if req.StartDate != "" {
		updates["start_date"] = req.StartDate
	}

	err := r.db.WithContext(ctx).
		Model(&model.Offer{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		code := "[REPOSITORY] UpdateOffer - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (r *offerRepository) DeleteOffer(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).
		Delete(&model.Offer{}, id).Error

	if err != nil {
		code := "[REPOSITORY] DeleteOffer - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewOfferRepository(db *gorm.DB) OfferRepository {
	return &offerRepository{db: db}
}
