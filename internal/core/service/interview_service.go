package service

import (
	"context"
	"time"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type InterviewService interface {
	GetInterviewsByTenant(ctx context.Context, tenantID int64, query entity.InterviewQueryString) ([]entity.InterviewEntity, int64, int64, error)
	GetInterviewByID(ctx context.Context, id int64) (*entity.InterviewEntity, error)
	CreateInterview(ctx context.Context, req entity.InterviewEntityRequest, tenantID int64) error
	UpdateInterview(ctx context.Context, id int64, req entity.InterviewUpdateRequest) error
	DeleteInterview(ctx context.Context, id int64) error
}

type interviewService struct {
	interviewRepo repository.InterviewRepository
}

func (s *interviewService) GetInterviewsByTenant(ctx context.Context, tenantID int64, query entity.InterviewQueryString) ([]entity.InterviewEntity, int64, int64, error) {
	results, totalData, totalPages, err := s.interviewRepo.GetInterviewsByTenant(ctx, tenantID, query)

	if err != nil {
		code := "[SERVICE] GetInterviewsByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPages, nil
}

func (s *interviewService) GetInterviewByID(ctx context.Context, id int64) (*entity.InterviewEntity, error) {
	result, err := s.interviewRepo.GetInterviewByID(ctx, id)

	if err != nil {
		code := "[SERVICE] GetInterviewByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (s *interviewService) CreateInterview(ctx context.Context, req entity.InterviewEntityRequest, tenantID int64) error {
	jakartaTZ, _ := time.LoadLocation("Asia/Jakarta")
	// Strip timezone from input and interpret as Jakarta time
	timeStr := req.ScheduledAt
	if len(timeStr) > 19 {
		timeStr = timeStr[:19]
	}
	scheduledAt, _ := time.ParseInLocation("2006-01-02T15:04:05", timeStr, jakartaTZ)

	reqEntity := entity.InterviewEntity{
		TenantID:               tenantID,
		CandidateApplicationID: req.CandidateApplicationID,
		InterviewerID:          req.InterviewerID,
		InterviewType:          req.InterviewType,
		ScheduledAt:            scheduledAt,
		DurationMinutes:        req.DurationMinutes,
		Location:               req.Location,
		MeetingLink:            req.MeetingLink,
	}

	err := s.interviewRepo.CreateInterview(ctx, reqEntity)

	if err != nil {
		code := "[SERVICE] CreateInterview - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (s *interviewService) UpdateInterview(ctx context.Context, id int64, req entity.InterviewUpdateRequest) error {
	err := s.interviewRepo.UpdateInterview(ctx, id, req)

	if err != nil {
		code := "[SERVICE] UpdateInterview - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (s *interviewService) DeleteInterview(ctx context.Context, id int64) error {
	err := s.interviewRepo.DeleteInterview(ctx, id)

	if err != nil {
		code := "[SERVICE] DeleteInterview - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewInterviewService(interviewRepo repository.InterviewRepository) InterviewService {
	return &interviewService{interviewRepo: interviewRepo}
}
