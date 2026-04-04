package service

import (
	"context"
	"errors"
	"fmt"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var ErrScoreTooLow = errors.New("match score below minimum threshold")

type JobPostingService interface {
	ApplyToJob(ctx context.Context, manpowerReq *entity.ManpowerReqEntity, fullName, email, phone, cvText string) (*JobApplyResult, error)
}

type JobApplyResult struct {
	Score         int32
	Verdict       string
	MatchedSkills []string
	MissingSkills []string
	Explanation   string
	Saved         bool
}

type jobPostingService struct {
	db            *gorm.DB
	candidateRepo repository.CandidateRepository
	appRepo       repository.CandidateApplicationRepository
	aiService     AIService
}

func NewJobPostingService(db *gorm.DB, candidateRepo repository.CandidateRepository, appRepo repository.CandidateApplicationRepository, aiService AIService) JobPostingService {
	return &jobPostingService{
		db:            db,
		candidateRepo: candidateRepo,
		appRepo:       appRepo,
		aiService:     aiService,
	}
}

func (s *jobPostingService) ApplyToJob(ctx context.Context, manpowerReq *entity.ManpowerReqEntity, fullName, email, phone, cvText string) (*JobApplyResult, error) {
	matchResp, err := s.aiService.MatchJob(ctx, cvText, manpowerReq.JobDescription)
	if err != nil {
		code := "[SERVICE] ApplyToJob - MatchJob"
		log.Errorw(code, err)
		return nil, fmt.Errorf("failed to analyze CV: %w", err)
	}

	result := &JobApplyResult{
		Score:         matchResp.Score,
		Verdict:       matchResp.Verdict,
		MatchedSkills: matchResp.MatchedSkills,
		MissingSkills: matchResp.MissingSkills,
		Explanation:   matchResp.Explanation,
		Saved:         false,
	}

	if matchResp.Score < 70 {
		return result, ErrScoreTooLow
	}

	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	candidate := &entity.CandidateEntity{
		TenantID: manpowerReq.TenantID,
		FullName: fullName,
		Email:    email,
		Phone:    phone,
		Source:   "JOB_POSTING",
		Status:   "ACTIVE",
	}

	if err := tx.Create(candidate).Error; err != nil {
		tx.Rollback()
		code := "[SERVICE] ApplyToJob - CreateCandidate"
		log.Errorw(code, err)
		return nil, fmt.Errorf("failed to save candidate: %w", err)
	}

	application := &entity.CandidateApplicationEntity{
		TenantID:          manpowerReq.TenantID,
		CandidateID:       candidate.ID,
		ManpowerRequestID: manpowerReq.ID,
		MatchScore:        matchResp.Score,
		MatchVerdict:      matchResp.Verdict,
		Status:            "APPLIED",
	}

	if err := tx.Create(application).Error; err != nil {
		tx.Rollback()
		code := "[SERVICE] ApplyToJob - CreateApplication"
		log.Errorw(code, err)
		return nil, fmt.Errorf("failed to save application: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		code := "[SERVICE] ApplyToJob - Commit"
		log.Errorw(code, err)
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	result.Saved = true
	return result, nil
}
