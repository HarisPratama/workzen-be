package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type InterviewRepository interface {
	GetInterviewsByTenant(ctx context.Context, tenantID int64, query entity.InterviewQueryString) ([]entity.InterviewEntity, int64, int64, error)
	GetInterviewByID(ctx context.Context, id int64) (*entity.InterviewEntity, error)
	CreateInterview(ctx context.Context, req entity.InterviewEntity) error
	UpdateInterview(ctx context.Context, id int64, req entity.InterviewUpdateRequest) error
	DeleteInterview(ctx context.Context, id int64) error
}

type interviewRepository struct {
	db *gorm.DB
}

func (r *interviewRepository) CreateInterview(ctx context.Context, req entity.InterviewEntity) error {
	modelInterview := model.Interview{
		TenantID:          req.TenantID,
		CandidateID:       req.CandidateID,
		ManpowerRequestID: req.ManpowerRequestID,
		InterviewType:     req.InterviewType,
		ScheduledAt:       req.ScheduledAt,
		DurationMinutes:   req.DurationMinutes,
		Location:          req.Location,
		MeetingLink:       req.MeetingLink,
		Status:            "SCHEDULED",
	}

	err := r.db.Create(&modelInterview).Error
	if err != nil {
		code := "[REPOSITORY] CreateInterview - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (r *interviewRepository) GetInterviewsByTenant(ctx context.Context, tenantID int64, query entity.InterviewQueryString) ([]entity.InterviewEntity, int64, int64, error) {
	var modelInterviews []model.Interview
	var countData int64

	sqlMain := r.db.WithContext(ctx).
		Model(&model.Interview{}).
		Preload("Candidate").
		Preload("ManpowerRequest").
		Where("tenant_id = ?", tenantID)

	if query.Search != "" {
		search := "%" + query.Search + "%"
		sqlMain = sqlMain.Where(`candidates.full_name ILIKE ?`, search)
	}

	if query.Status != "" {
		sqlMain = sqlMain.Where("interviews.status = ?", query.Status)
	}

	if query.CandidateID != 0 {
		sqlMain = sqlMain.Where("interviews.candidate_id = ?", query.CandidateID)
	}

	if query.StartDate != "" && query.EndDate != "" {
		sqlMain = sqlMain.Where("interviews.scheduled_at BETWEEN ? AND ?", query.StartDate, query.EndDate)
	}

	countQuery := sqlMain.Session(&gorm.Session{})

	if err := countQuery.Count(&countData).Error; err != nil {
		return nil, 0, 0, err
	}

	allowedOrder := map[string]string{
		"scheduled_at": "interviews.scheduled_at",
		"created_at":   "interviews.created_at",
	}

	orderBy, ok := allowedOrder[query.OrderBy]
	if !ok {
		orderBy = "interviews.scheduled_at"
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
		Find(&modelInterviews).Error; err != nil {
		return nil, 0, 0, err
	}

	var result []entity.InterviewEntity
	for _, item := range modelInterviews {
		result = append(result, entity.InterviewEntity{
			ID:                item.ID,
			TenantID:          item.TenantID,
			CandidateID:       item.CandidateID,
			ManpowerRequestID: item.ManpowerRequestID,
			Status:            item.Status,
			InterviewType:     item.InterviewType,
			ScheduledAt:       item.ScheduledAt,
			DurationMinutes:   item.DurationMinutes,
			Location:          item.Location,
			MeetingLink:       item.MeetingLink,
			Feedback:          item.Feedback,
			Rating:            item.Rating,
			CompletedAt:       item.CompletedAt,
			CancelledAt:       item.CancelledAt,
			CancelReason:      item.CancelReason,
			Tenant:            item.Tenant,
			Candidate:         item.Candidate,
			ManpowerRequest:   item.ManpowerRequest,
		})
	}

	return result, countData, totalPages, nil
}

func (r *interviewRepository) GetInterviewByID(ctx context.Context, id int64) (*entity.InterviewEntity, error) {
	var modelInterview model.Interview
	err := r.db.WithContext(ctx).
		Preload("Candidate").
		Preload("ManpowerRequest").
		First(&modelInterview, id).Error

	if err != nil {
		return nil, err
	}

	return &entity.InterviewEntity{
		ID:                modelInterview.ID,
		TenantID:          modelInterview.TenantID,
		CandidateID:       modelInterview.CandidateID,
		ManpowerRequestID: modelInterview.ManpowerRequestID,
		Status:            modelInterview.Status,
		InterviewType:     modelInterview.InterviewType,
		ScheduledAt:       modelInterview.ScheduledAt,
		DurationMinutes:   modelInterview.DurationMinutes,
		Location:          modelInterview.Location,
		MeetingLink:       modelInterview.MeetingLink,
		Feedback:          modelInterview.Feedback,
		Rating:            modelInterview.Rating,
		CompletedAt:       modelInterview.CompletedAt,
		CancelledAt:       modelInterview.CancelledAt,
		CancelReason:      modelInterview.CancelReason,
		Tenant:            modelInterview.Tenant,
		Candidate:         modelInterview.Candidate,
		ManpowerRequest:   modelInterview.ManpowerRequest,
	}, nil
}

func (r *interviewRepository) UpdateInterview(ctx context.Context, id int64, req entity.InterviewUpdateRequest) error {
	updates := map[string]interface{}{
		"status":    req.Status,
		"feedback":  req.Feedback,
		"rating":    req.Rating,
		"location":  req.Location,
		"scheduled_at": req.ScheduledAt,
	}

	if req.Status == "COMPLETED" {
		updates["completed_at"] = time.Now()
	} else if req.Status == "CANCELLED" {
		updates["cancelled_at"] = time.Now()
	}

	err := r.db.WithContext(ctx).
		Model(&model.Interview{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		code := "[REPOSITORY] UpdateInterview - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (r *interviewRepository) DeleteInterview(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).
		Delete(&model.Interview{}, id).Error

	if err != nil {
		code := "[REPOSITORY] DeleteInterview - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewInterviewRepository(db *gorm.DB) InterviewRepository {
	return &interviewRepository{db: db}
}