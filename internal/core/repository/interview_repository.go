package repository

import (
	"context"

	"github.com/google/uuid"
	"workzen-be/internal/core/domain/entity"
)

type InterviewRepository interface {
	Create(ctx context.Context, interview *entity.Interview) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Interview, error)
	FindByCandidateApplicationID(ctx context.Context, candidateApplicationID uuid.UUID) ([]entity.Interview, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]entity.Interview, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.Interview, error)
	Update(ctx context.Context, interview *entity.Interview) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int64, error)
}