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

type OfferService interface {
	CreateOffer(ctx context.Context, offer entity.Offer) (*entity.Offer, error)
	UpdateOffer(ctx context.Context, id uint, offer entity.Offer) (*entity.Offer, error)
	DeleteOffer(ctx context.Context, id uint) error
	GetOfferByID(ctx context.Context, id uint) (*entity.Offer, error)
	GetOffersByTenant(ctx context.Context, tenantID uint, page, limit int) ([]entity.Offer, int64, error)
	GetOffersByCandidate(ctx context.Context, candidateApplicationID uint, page, limit int) ([]entity.Offer, int64, error)
	SendOffer(ctx context.Context, id uint) error
	WithdrawOffer(ctx context.Context, id uint, reason string) error
	AcceptOffer(ctx context.Context, id uint) error
	RejectOffer(ctx context.Context, id uint, reason string) error
	NegotiateOffer(ctx context.Context, id uint, proposedSalary float64, proposedBenefits string, justification string) error
	GetOfferMetrics(ctx context.Context, tenantID uint) (map[string]interface{}, error)
}

type offerService struct {
	offerRepo repository.OfferRepository
}

func NewOfferService(offerRepo repository.OfferRepository) OfferService {
	return &offerService{
		offerRepo: offerRepo,
	}
}

func (s *offerService) CreateOffer(ctx context.Context, offer entity.Offer) (*entity.Offer, error) {
	if err := validator.ValidateStruct(offer); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	offer.Status = entity.OfferStatusDraft

	if err := s.offerRepo.Create(ctx, &offer); err != nil {
		log.Errorw("failed to create offer", "error", err)
		return nil, fmt.Errorf("failed to create offer: %w", err)
	}

	return &offer, nil
}

func (s *offerService) UpdateOffer(ctx context.Context, id uint, offer entity.Offer) (*entity.Offer, error) {
	existing, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("offer not found: %w", err)
	}

	// Only allow updates for draft or negotiating offers
	if existing.Status != entity.OfferStatusDraft && existing.Status != entity.OfferStatusNegotiating {
		return nil, fmt.Errorf("cannot update offer in status: %s", existing.Status)
	}

	existing.Position = offer.Position
	existing.Department = offer.Department
	existing.EmploymentType = offer.EmploymentType
	existing.BaseSalary = offer.BaseSalary
	existing.Currency = offer.Currency
	existing.Bonus = offer.Bonus
	existing.Benefits = offer.Benefits
	existing.ProbationPeriodMonths = offer.ProbationPeriodMonths
	existing.NoticePeriodDays = offer.NoticePeriodDays
	existing.StartDate = offer.StartDate
	existing.ExpiryDate = offer.ExpiryDate
	existing.Notes = offer.Notes
	existing.Terms = offer.Terms

	if err := s.offerRepo.Update(ctx, existing); err != nil {
		log.Errorw("failed to update offer", "error", err)
		return nil, fmt.Errorf("failed to update offer: %w", err)
	}

	return existing, nil
}

func (s *offerService) DeleteOffer(ctx context.Context, id uint) error {
	existing, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offer not found: %w", err)
	}

	// Only allow deletion of draft offers
	if existing.Status != entity.OfferStatusDraft {
		return fmt.Errorf("cannot delete offer in status: %s", existing.Status)
	}

	return s.offerRepo.Delete(ctx, id)
}

func (s *offerService) GetOfferByID(ctx context.Context, id uint) (*entity.Offer, error) {
	offer, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("offer not found: %w", err)
	}
	return offer, nil
}

func (s *offerService) GetOffersByTenant(ctx context.Context, tenantID uint, page, limit int) ([]entity.Offer, int64, error) {
	offset := (page - 1) * limit
	count, err := s.offerRepo.CountByCompanyID(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	offers, err := s.offerRepo.FindByCompanyID(ctx, tenantID, limit, offset)
	return offers, count, err
}

func (s *offerService) GetOffersByCandidate(ctx context.Context, candidateApplicationID uint, page, limit int) ([]entity.Offer, int64, error) {
	offer, err := s.offerRepo.FindByCandidateApplicationID(ctx, candidateApplicationID)
	if err != nil {
		return nil, 0, err
	}
	return []entity.Offer{*offer}, 1, nil
}

func (s *offerService) SendOffer(ctx context.Context, id uint) error {
	offer, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offer not found: %w", err)
	}

	if offer.Status != entity.OfferStatusDraft {
		return fmt.Errorf("cannot send offer in status: %s", offer.Status)
	}

	offer.Status = entity.OfferStatusSent
	sentAt := time.Now()
	offer.SentAt = &sentAt

	return s.offerRepo.Update(ctx, offer)
}

func (s *offerService) WithdrawOffer(ctx context.Context, id uint, reason string) error {
	offer, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offer not found: %w", err)
	}

	if offer.Status != entity.OfferStatusSent && offer.Status != entity.OfferStatusNegotiating {
		return fmt.Errorf("cannot withdraw offer in status: %s", offer.Status)
	}

	offer.Status = entity.OfferStatusWithdrawn
	// Note: If you want to store reason, you'd need to add a field to the entity

	return s.offerRepo.Update(ctx, offer)
}

func (s *offerService) AcceptOffer(ctx context.Context, id uint) error {
	offer, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offer not found: %w", err)
	}

	if offer.Status != entity.OfferStatusSent && offer.Status != entity.OfferStatusNegotiating {
		return fmt.Errorf("cannot accept offer in status: %s", offer.Status)
	}

	offer.Status = entity.OfferStatusAccepted
	respondedAt := time.Now()
	offer.RespondedAt = &respondedAt

	return s.offerRepo.Update(ctx, offer)
}

func (s *offerService) RejectOffer(ctx context.Context, id uint, reason string) error {
	offer, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offer not found: %w", err)
	}

	if offer.Status != entity.OfferStatusSent && offer.Status != entity.OfferStatusNegotiating {
		return fmt.Errorf("cannot reject offer in status: %s", offer.Status)
	}

	offer.Status = entity.OfferStatusRejected
	respondedAt := time.Now()
	offer.RespondedAt = &respondedAt
	// Note: If you want to store reason, you'd need to add a field to the entity

	return s.offerRepo.Update(ctx, offer)
}

func (s *offerService) NegotiateOffer(ctx context.Context, id uint, proposedSalary float64, proposedBenefits string, justification string) error {
	offer, err := s.offerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("offer not found: %w", err)
	}

	if offer.Status != entity.OfferStatusSent {
		return fmt.Errorf("cannot negotiate offer in status: %s", offer.Status)
	}

	offer.Status = entity.OfferStatusNegotiating
	offer.NegotiationCounter = &proposedSalary
	offer.NegotiationNotes = &justification

	return s.offerRepo.Update(ctx, offer)
}

func (s *offerService) GetOfferMetrics(ctx context.Context, tenantID uint) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	
	// This is a placeholder - implement actual metrics logic based on repository queries
	metrics["total_offers"] = 0
	metrics["draft_offers"] = 0
	metrics["sent_offers"] = 0
	metrics["accepted_offers"] = 0
	metrics["rejected_offers"] = 0
	metrics["negotiating_offers"] = 0
	metrics["withdrawn_offers"] = 0
	
	return metrics, nil
}