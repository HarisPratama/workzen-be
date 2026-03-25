package repository

import (
	"context"

	"github.com/google/uuid"
	"workzen-be/internal/core/domain/entity"
)

type OfferRepository interface {
	Create(ctx context.Context, offer *entity.Offer) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Offer, error)
	FindByCandidateApplicationID(ctx context.Context, candidateApplicationID uuid.UUID) (*entity.Offer, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.Offer, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status entity.OfferStatus, limit, offset int) ([]entity.Offer, error)
	Update(ctx context.Context, offer *entity.Offer) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int64, error)
}