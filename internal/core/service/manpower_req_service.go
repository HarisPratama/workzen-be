package service

import (
	"context"
	"fmt"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type ManpowerReqService interface {
	GetManpowerReqByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error)
	GetDetailManpowerRequestByTenant(ctx context.Context, tenantID int64, id int64) (*entity.ManpowerReqEntity, error)
	GetManpowerReqByPublicToken(ctx context.Context, token string) (*entity.ManpowerReqEntity, error)
	GeneratePublicLink(ctx context.Context, tenantID int64, manpowerReqID int64) (string, error)
	CreateManpowerReq(ctx context.Context, manpowerReq *entity.ManpowerReqEntity) error
	UpdateManpowerReq(ctx context.Context, tenantID int64, manpowerReqID int64, manpowerReq *entity.ManpowerReqEntity) error
	DeleteManpowerReq(ctx context.Context, tenantID int64, manpowerReqID int64) error
}

type manpowerReqService struct {
	manpowerReqRepo repository.ManpowerReqRepository
}

func (m *manpowerReqService) GetDetailManpowerRequestByTenant(ctx context.Context, tenantID int64, id int64) (*entity.ManpowerReqEntity, error) {
	result, err := m.manpowerReqRepo.GetDetailManpowerRequestByTenant(ctx, tenantID, id)
	if err != nil {
		code = "[SERVICE] GetDetailManpowerRequestByTenant - 1"
		log.Error(code, err)
		return nil, err
	}

	return result, nil
}

func (m *manpowerReqService) GetManpowerReqByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error) {
	results, totalData, totalPage, err := m.manpowerReqRepo.GetManpowerReqsByTenant(ctx, tenantID, query)
	if err != nil {
		code = "[SERVICE] GetManpowerReqByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPage, nil
}

func (m *manpowerReqService) CreateManpowerReq(ctx context.Context, manpowerReq *entity.ManpowerReqEntity) error {
	err := m.manpowerReqRepo.CreateManpowerReq(ctx, manpowerReq)
	if err != nil {
		code = "[SERVICE] CreateManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (m *manpowerReqService) UpdateManpowerReq(ctx context.Context, tenantID int64, manpowerReqID int64, manpowerReq *entity.ManpowerReqEntity) error {
	existing, err := m.manpowerReqRepo.GetDetailManpowerRequestByTenant(ctx, tenantID, manpowerReqID)
	if err != nil {
		return fmt.Errorf("manpower request not found: %w", err)
	}

	manpowerReq.ID = existing.ID
	manpowerReq.TenantID = existing.TenantID
	manpowerReq.CreatedAt = existing.CreatedAt

	err = m.manpowerReqRepo.UpdateManpowerReq(ctx, manpowerReq)
	if err != nil {
		code = "[SERVICE] UpdateManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (m *manpowerReqService) DeleteManpowerReq(ctx context.Context, tenantID int64, manpowerReqID int64) error {
	_, err := m.manpowerReqRepo.GetDetailManpowerRequestByTenant(ctx, tenantID, manpowerReqID)
	if err != nil {
		return fmt.Errorf("manpower request not found: %w", err)
	}

	err = m.manpowerReqRepo.DeleteManpowerReq(ctx, tenantID, manpowerReqID)
	if err != nil {
		code = "[SERVICE] DeleteManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (m *manpowerReqService) GetManpowerReqByPublicToken(ctx context.Context, token string) (*entity.ManpowerReqEntity, error) {
	result, err := m.manpowerReqRepo.GetManpowerReqByPublicToken(ctx, token)
	if err != nil {
		code = "[SERVICE] GetManpowerReqByPublicToken - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

func (m *manpowerReqService) GeneratePublicLink(ctx context.Context, tenantID int64, manpowerReqID int64) (string, error) {
	existing, err := m.manpowerReqRepo.GetDetailManpowerRequestByTenant(ctx, tenantID, manpowerReqID)
	if err != nil {
		return "", fmt.Errorf("manpower request not found: %w", err)
	}

	if existing.PublicToken != "" {
		return existing.PublicToken, nil
	}

	token := uuid.New().String()
	err = m.manpowerReqRepo.UpdatePublicToken(ctx, tenantID, manpowerReqID, token)
	if err != nil {
		code = "[SERVICE] GeneratePublicLink - 1"
		log.Errorw(code, err)
		return "", err
	}

	return token, nil
}

func NewManpowerReqService(repo repository.ManpowerReqRepository) ManpowerReqService {
	return &manpowerReqService{
		manpowerReqRepo: repo,
	}
}
