package service

import (
	"context"
	"fmt"
	"time"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/repository"
	"workzen-be/lib/validator"

	"github.com/gofiber/fiber/v2/log"
)

type InterviewService interface {
	CreateInterview(ctx context.Context, interview entity.Interview) (*entity.Interview, error)
	UpdateInterview(ctx context.Context, id uint, interview entity.Interview) (*entity.Interview, error)
	DeleteInterview(ctx context.Context, id uint) error
	GetInterviewByID(ctx context.Context, id uint) (*entity.Interview, error)
	GetInterviewsByTenant(ctx context.Context, tenantID uint, page, limit int) ([]entity.Interview, int64, error)
	GetInterviewsByCandidate(ctx context.Context, candidateApplicationID uint, page, limit int) ([]entity.Interview, int64, error)
	GetInterviewsByInterviewer(ctx context.Context, employeeID uint, page, limit int) ([]entity.Interview, int64, error)
	RescheduleInterview(ctx context.Context, id uint, newScheduledAt time.Time, reason string) error
	CancelInterview(ctx context.Context, id uint, reason string) error
	CompleteInterview(ctx context.Context, id uint, feedback string) error
	SubmitFeedback(ctx context.Context, id uint, feedback string, rating int, recommendation string) error
	GetInterviewMetrics(ctx context.Context, tenantID uint) (map[string]interface{}, error)
}

type interviewService struct {
	interviewRepo repository.InterviewRepository
}

func NewInterviewService(interviewRepo repository.InterviewRepository) InterviewService {
	return &interviewService{
		interviewRepo: interviewRepo,
	}
}

func (s *interviewService) CreateInterview(ctx context.Context, interview entity.Interview) (*entity.Interview, error) {
	if err := validator.ValidateStruct(interview); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	interview.Status = entity.InterviewStatusScheduled

	if err := s.interviewRepo.Create(ctx, &interview); err != nil {
		log.Errorw("failed to create interview", "error", err)
		return nil, fmt.Errorf("failed to create interview: %w", err)
	}

	return &interview, nil
}

func (s *interviewService) UpdateInterview(ctx context.Context, id uint, interview entity.Interview) (*entity.Interview, error) {
	existing, err := s.interviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("interview not found: %w", err)
	}

	if existing.Status != entity.InterviewStatusScheduled && existing.Status != entity.InterviewStatusRescheduled {
		return nil, fmt.Errorf("cannot update interview in status: %s", existing.Status)
	}

	existing.ScheduledAt = interview.ScheduledAt
	existing.DurationMinutes = interview.DurationMinutes
	existing.InterviewType = interview.InterviewType
	existing.Location = interview.Location
	existing.MeetingLink = interview.MeetingLink

	if err := s.interviewRepo.Update(ctx, existing); err != nil {
		log.Errorw("failed to update interview", "error", err)
		return nil, fmt.Errorf("failed to update interview: %w", err)
	}

	return existing, nil
}

func (s *interviewService) DeleteInterview(ctx context.Context, id uint) error {
	existing, err := s.interviewRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	if existing.Status != entity.InterviewStatusScheduled && existing.Status != entity.InterviewStatusRescheduled {
		return fmt.Errorf("cannot delete interview in status: %s", existing.Status)
	}

	return s.interviewRepo.Delete(ctx, id)
}

func (s *interviewService) GetInterviewByID(ctx context.Context, id uint) (*entity.Interview, error) {
	interview, err := s.interviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("interview not found: %w", err)
	}
	return interview, nil
}

func (s *interviewService) GetInterviewsByTenant(ctx context.Context, tenantID uint, page, limit int) ([]entity.Interview, int64, error) {
	offset := (page - 1) * limit
	count, err := s.interviewRepo.CountByCompanyID(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	interviews, err := s.interviewRepo.FindByCompanyID(ctx, tenantID, limit, offset)
	return interviews, count, err
}

func (s *interviewService) GetInterviewsByCandidate(ctx context.Context, candidateApplicationID uint, page, limit int) ([]entity.Interview, int64, error) {
	interviews, err := s.interviewRepo.FindByCandidateApplicationID(ctx, candidateApplicationID)
	if err != nil {
		return nil, 0, err
	}
	return interviews, int64(len(interviews)), nil
}

func (s *interviewService) GetInterviewsByInterviewer(ctx context.Context, employeeID uint, page, limit int) ([]entity.Interview, int64, error) {
	interviews, err := s.interviewRepo.FindByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, 0, err
	}
	return interviews, int64(len(interviews)), nil
}

func (s *interviewService) RescheduleInterview(ctx context.Context, id uint, newScheduledAt time.Time, reason string) error {
	interview, err := s.interviewRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	if interview.Status != entity.InterviewStatusScheduled && interview.Status != entity.InterviewStatusRescheduled {
		return fmt.Errorf("cannot reschedule interview in status: %s", interview.Status)
	}

	interview.ScheduledAt = newScheduledAt
	interview.Status = entity.InterviewStatusRescheduled
	if reason != "" {
		interview.Notes = &reason
	}

	return s.interviewRepo.Update(ctx, interview)
}

func (s *interviewService) CancelInterview(ctx context.Context, id uint, reason string) error {
	interview, err := s.interviewRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	if interview.Status == entity.InterviewStatusCompleted || interview.Status == entity.InterviewStatusCancelled {
		return fmt.Errorf("cannot cancel interview in status: %s", interview.Status)
	}

	interview.Status = entity.InterviewStatusCancelled
	cancelTime := time.Now()
	interview.CancelledAt = &cancelTime
	if reason != "" {
		interview.CancelReason = &reason
	}

	return s.interviewRepo.Update(ctx, interview)
}

func (s *interviewService) CompleteInterview(ctx context.Context, id uint, feedback string) error {
	interview, err := s.interviewRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	if interview.Status != entity.InterviewStatusScheduled && interview.Status != entity.InterviewStatusRescheduled {
		return fmt.Errorf("cannot complete interview in status: %s", interview.Status)
	}

	interview.Status = entity.InterviewStatusCompleted
	completedTime := time.Now()
	interview.CompletedAt = &completedTime
	if feedback != "" {
		interview.Feedback = &feedback
	}

	return s.interviewRepo.Update(ctx, interview)
}

func (s *interviewService) SubmitFeedback(ctx context.Context, id uint, feedback string, rating int, recommendation string) error {
	interview, err := s.interviewRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("interview not found: %w", err)
	}

	if interview.Status != entity.InterviewStatusCompleted {
		return fmt.Errorf("can only submit feedback for completed interviews")
	}

	interview.Feedback = &feedback
	interview.Rating = &rating

	return s.interviewRepo.Update(ctx, interview)
}

func (s *interviewService) GetInterviewMetrics(ctx context.Context, tenantID uint) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	
	// This is a placeholder - implement actual metrics logic based on repository queries
	metrics["total_interviews"] = 0
	metrics["completed_interviews"] = 0
	metrics["cancelled_interviews"] = 0
	metrics["pending_interviews"] = 0
	
	return metrics, nil
}