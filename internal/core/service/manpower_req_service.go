package service

import (
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type ManpowerReqService interface {
	GetManpowerReqByTenant(ctx context.Context, tenantID int64, query entity.ManpowerReqQueryString) ([]entity.ManpowerReqEntity, int64, int64, error)
	CreateManpowerReq(ctx context.Context, req entity.ManpowerReqEntity) error
}

type manpowerReqService struct {
	manpowerReqRepo repository.ManpowerReqRepository
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

func (m *manpowerReqService) CreateManpowerReq(ctx context.Context, req entity.ManpowerReqEntity) error {
	err := m.manpowerReqRepo.CreateManpowerReq(ctx, req)

	if err != nil {
		code = "[SERVICE] CreateManpowerReq - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewManpowerReqService(repo repository.ManpowerReqRepository) ManpowerReqService {
	return &manpowerReqService{
		manpowerReqRepo: repo,
	}
}
