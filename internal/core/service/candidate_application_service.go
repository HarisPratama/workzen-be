package service

import (
	"context"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type CandidateApplicationService interface {
	GetCandidateApplicationByTenantMR(ctx context.Context, tenantID int64, manpowerRequestID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error)
	CreateCandidateApplication(ctx context.Context, req entity.CandidateApplicationEntity) error
}

type candidateApplicationService struct {
	candidateApplicationRepo repository.CandidateApplicationRepository
}

func (c *candidateApplicationService) CreateCandidateApplication(ctx context.Context, req entity.CandidateApplicationEntity) error {
	err := c.candidateApplicationRepo.CreateCandidateApplication(ctx, req)

	if err != nil {
		code = "[SERVICE] CreateCandidateApplication - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateApplicationService) GetCandidateApplicationByTenantMR(ctx context.Context, tenantID int64, manpowerRequestID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error) {
	results, totalData, totalPage, err := c.candidateApplicationRepo.GetCandidateApplicationByTenantMR(ctx, tenantID, manpowerRequestID, query)

	if err != nil {
		code = "[SERVICE] GetCandidateApplicationByTenantMR - 1"
		log.Errorw(code, err)
		return nil, 0, 0, nil
	}

	return results, totalData, totalPage, nil
}

func NewCandidateApplicationService(candidateApplicationRepo repository.CandidateApplicationRepository) CandidateApplicationService {
	return &candidateApplicationService{
		candidateApplicationRepo: candidateApplicationRepo,
	}
}
