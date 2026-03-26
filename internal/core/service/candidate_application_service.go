package service

import (
	"context"
	"fmt"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type CandidateApplicationService interface {
	GetCandidateApplicationsByTenant(ctx context.Context, tenantID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error)
	GetCandidateApplicationByTenantMR(ctx context.Context, tenantID int64, manpowerRequestID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error)
	GetDetailCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) (*entity.CandidateApplicationEntity, error)
	CreateCandidateApplication(ctx context.Context, application *entity.CandidateApplicationEntity) error
	UpdateCandidateApplication(ctx context.Context, tenantID int64, applicationID int64, application *entity.CandidateApplicationEntity) error
	DeleteCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) error
}

type candidateApplicationService struct {
	candidateApplicationRepo repository.CandidateApplicationRepository
}

func (c *candidateApplicationService) GetCandidateApplicationsByTenant(ctx context.Context, tenantID int64, query entity.CandidateApplicationQueryString) ([]entity.CandidateApplicationEntity, int64, int64, error) {
	results, totalData, totalPage, err := c.candidateApplicationRepo.GetCandidateApplicationsByTenant(ctx, tenantID, query)
	if err != nil {
		code = "[SERVICE] GetCandidateApplicationsByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPage, nil
}

func (c *candidateApplicationService) CreateCandidateApplication(ctx context.Context, application *entity.CandidateApplicationEntity) error {
	err := c.candidateApplicationRepo.CreateCandidateApplication(ctx, application)
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

func (c *candidateApplicationService) GetDetailCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) (*entity.CandidateApplicationEntity, error) {
	application, err := c.candidateApplicationRepo.GetDetailCandidateApplication(ctx, tenantID, applicationID)
	if err != nil {
		return nil, fmt.Errorf("candidate application not found: %w", err)
	}

	return application, nil
}

func (c *candidateApplicationService) UpdateCandidateApplication(ctx context.Context, tenantID int64, applicationID int64, application *entity.CandidateApplicationEntity) error {
	existing, err := c.candidateApplicationRepo.GetDetailCandidateApplication(ctx, tenantID, applicationID)
	if err != nil {
		return fmt.Errorf("candidate application not found: %w", err)
	}

	application.ID = existing.ID
	application.TenantID = existing.TenantID
	application.CandidateID = existing.CandidateID
	application.ManpowerRequestID = existing.ManpowerRequestID
	application.AppliedAt = existing.AppliedAt

	err = c.candidateApplicationRepo.UpdateCandidateApplication(ctx, application)
	if err != nil {
		code = "[SERVICE] UpdateCandidateApplication - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *candidateApplicationService) DeleteCandidateApplication(ctx context.Context, tenantID int64, applicationID int64) error {
	_, err := c.candidateApplicationRepo.GetDetailCandidateApplication(ctx, tenantID, applicationID)
	if err != nil {
		return fmt.Errorf("candidate application not found: %w", err)
	}

	err = c.candidateApplicationRepo.DeleteCandidateApplication(ctx, tenantID, applicationID)
	if err != nil {
		code = "[SERVICE] DeleteCandidateApplication - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewCandidateApplicationService(candidateApplicationRepo repository.CandidateApplicationRepository) CandidateApplicationService {
	return &candidateApplicationService{
		candidateApplicationRepo: candidateApplicationRepo,
	}
}
