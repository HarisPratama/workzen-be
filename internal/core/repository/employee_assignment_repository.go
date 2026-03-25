package repository

import (
	"context"

	"github.com/google/uuid"
	"workzen-be/internal/core/domain/entity"
)

type EmployeeAssignmentRepository interface {
	Create(ctx context.Context, assignment *entity.EmployeeAssignment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.EmployeeAssignment, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, limit, offset int) ([]entity.EmployeeAssignment, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]entity.EmployeeAssignment, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.EmployeeAssignment, error)
	FindByStatus(ctx context.Context, tenantID uuid.UUID, status entity.AssignmentStatus, limit, offset int) ([]entity.EmployeeAssignment, error)
	FindActiveByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]entity.EmployeeAssignment, error)
	Update(ctx context.Context, assignment *entity.EmployeeAssignment) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int64, error)
}