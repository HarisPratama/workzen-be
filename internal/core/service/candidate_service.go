package service

import (
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type CandidateService interface {
	GetCandidatesByTenant(ctx context.Context, tenantID int64, query entity.CandidateQueryString) ([]entity.CandidateEntity, int64, int64, error)
	CreateCandidate(ctx context.Context, candidate entity.CandidateEntity) error
}

type candidateService struct {
	candidateRepository repository.CandidateRepository
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
