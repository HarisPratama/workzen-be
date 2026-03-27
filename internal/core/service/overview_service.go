package service

import (
	"context"
	"workzen-be/internal/adapter/repository"

	"github.com/gofiber/fiber/v2/log"
)

type OverviewService interface {
	GetOverviewByTenant(ctx context.Context, tenantID int64) (*repository.OverviewData, error)
}

type overviewService struct {
	overviewRepo repository.OverviewRepository
}

func (s *overviewService) GetOverviewByTenant(ctx context.Context, tenantID int64) (*repository.OverviewData, error) {
	result, err := s.overviewRepo.GetOverviewByTenant(ctx, tenantID)
	if err != nil {
		code = "[SERVICE] GetOverviewByTenant - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func NewOverviewService(repo repository.OverviewRepository) OverviewService {
	return &overviewService{overviewRepo: repo}
}
