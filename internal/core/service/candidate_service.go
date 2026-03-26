package service

import (
	"context"
	"fmt"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type CandidateService interface {
	GetCandidatesByTenant(ctx context.Context, tenantID int64, query entity.CandidateQueryString) ([]entity.CandidateEntity, int64, int64, error)
	GetDetailCandidateByTenant(ctx context.Context, tenantID int64, candidateID int64) (*entity.CandidateEntity, error)
	CreateCandidate(ctx context.Context, candidate entity.CandidateEntity) error
}

type candidateService struct {
	candidateRepository repository.CandidateRepository
}

func (c *candidateService) GetDetailCandidateByTenant(ctx context.Context, tenantID int64, candidateID int64) (*entity.CandidateEntity, error) {
	candidate, err := c.candidateRepository.GetDetailCandidateByTenant(ctx, tenantID, candidateID)
	if err != nil {
		return nil, fmt.Errorf("attendance not found: %w", err)
	}

	return candidate, nil
}

func (c *candidateService) GetCandidatesByTenant(ctx context.Context, tenantID int64, query entity.CandidateQueryString) ([]entity.CandidateEntity, int64, int64, error) {
	results, totalData, totalPage, err := c.candidateRepository.GetCandidatesByTenant(ctx, tenantID, query)
	if err != nil {
		code = "[SERVICE] GetCandidatesByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPage, nil
}

func (c *candidateService) CreateCandidate(ctx context.Context, req entity.CandidateEntity) error {
	err := c.candidateRepository.CreateCandidate(ctx, req)

	if err != nil {
		code = "[SERVICE] CreateCandidate - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewCandidateService(candidateRepository repository.CandidateRepository) CandidateService {
	return &candidateService{candidateRepository: candidateRepository}
}
