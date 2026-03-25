package service

import (
	"bwanews/internal/adapter/repository"
	"bwanews/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type OfferService interface {
	GetOffersByTenant(ctx context.Context, tenantID int64, query entity.OfferQueryString) ([]entity.OfferEntity, int64, int64, error)
	GetOfferByID(ctx context.Context, id int64) (*entity.OfferEntity, error)
	CreateOffer(ctx context.Context, req entity.OfferEntityRequest, tenantID int64) error
	UpdateOffer(ctx context.Context, id int64, req entity.OfferUpdateRequest) error
	DeleteOffer(ctx context.Context, id int64) error
}

type offerService struct {
	offerRepo repository.OfferRepository
}

func (s *offerService) GetOffersByTenant(ctx context.Context, tenantID int64, query entity.OfferQueryString) ([]entity.OfferEntity, int64, int64, error) {
	results, totalData, totalPages, err := s.offerRepo.GetOffersByTenant(ctx, tenantID, query)

	if err != nil {
		code := "[SERVICE] GetOffersByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPages, nil
}

func (s *offerService) GetOfferByID(ctx context.Context, id int64) (*entity.OfferEntity, error) {
	result, err := s.offerRepo.GetOfferByID(ctx, id)

	if err != nil {
		code := "[SERVICE] GetOfferByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (s *offerService) CreateOffer(ctx context.Context, req entity.OfferEntityRequest, tenantID int64) error {
	reqEntity := entity.OfferEntity{
		TenantID:               tenantID,
		CandidateApplicationID: req.CandidateApplicationID,
		Position:               req.Position,
		Department:             req.Department,
		EmploymentType:         req.EmploymentType,
		BaseSalary:             req.BaseSalary,
		Currency:               req.Currency,
		Bonus:                  req.Bonus,
		Benefits:               req.Benefits,
		ProbationPeriodMonths:  req.ProbationPeriodMonths,
		NoticePeriodDays:       req.NoticePeriodDays,
		StartDate:              req.StartDate,
		ExpiryDate:             req.ExpiryDate,
	}

	err := s.offerRepo.CreateOffer(ctx, reqEntity)

	if err != nil {
		code := "[SERVICE] CreateOffer - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (s *offerService) UpdateOffer(ctx context.Context, id int64, req entity.OfferUpdateRequest) error {
	err := s.offerRepo.UpdateOffer(ctx, id, req)

	if err != nil {
		code := "[SERVICE] UpdateOffer - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (s *offerService) DeleteOffer(ctx context.Context, id int64) error {
	err := s.offerRepo.DeleteOffer(ctx, id)

	if err != nil {
		code := "[SERVICE] DeleteOffer - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewOfferService(offerRepo repository.OfferRepository) OfferService {
	return &offerService{offerRepo: offerRepo}
}