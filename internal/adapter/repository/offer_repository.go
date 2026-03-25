package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"workzen-be/internal/core/domain/entity"
)

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) OfferRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) Create(ctx context.Context, offer *entity.Offer) error {
	return r.db.WithContext(ctx).Create(offer).Error
}

func (r *offerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Offer, error) {
	var offer entity.Offer
	err := r.db.WithContext(ctx).First(&offer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (r *offerRepository) FindByCandidateApplicationID(ctx context.Context, candidateApplicationID uuid.UUID) (*entity.Offer, error) {
	var offer entity.Offer
	err := r.db.WithContext(ctx).Where("candidate_application_id = ?", candidateApplicationID).First(&offer).Error
	return &offer, err
}

func (r *offerRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.Offer, error) {
	var offers []entity.Offer
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&offers).Error
	return offers, err
}

func (r *offerRepository) FindByStatus(ctx context.Context, tenantID uuid.UUID, status entity.OfferStatus, limit, offset int) ([]entity.Offer, error) {
	var offers []entity.Offer
	query := r.db.WithContext(ctx).Where("tenant_id = ? AND status = ?", tenantID, status)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&offers).Error
	return offers, err
}

func (r *offerRepository) Update(ctx context.Context, offer *entity.Offer) error {
	return r.db.WithContext(ctx).Save(offer).Error
}

func (r *offerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Offer{}, "id = ?", id).Error
}

func (r *offerRepository) CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Offer{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}
