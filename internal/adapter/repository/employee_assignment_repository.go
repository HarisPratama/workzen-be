package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"workzen-be/internal/core/domain/entity"
)

type employeeAssignmentRepository struct {
	db *gorm.DB
}

func NewEmployeeAssignmentRepository(db *gorm.DB) EmployeeAssignmentRepository {
	return &employeeAssignmentRepository{db: db}
}

func (r *employeeAssignmentRepository) Create(ctx context.Context, assignment *entity.EmployeeAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *employeeAssignmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.EmployeeAssignment, error) {
	var assignment entity.EmployeeAssignment
	err := r.db.WithContext(ctx).First(&assignment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *employeeAssignmentRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, limit, offset int) ([]entity.EmployeeAssignment, error) {
	var assignments []entity.EmployeeAssignment
	query := r.db.WithContext(ctx).Where("employee_id = ?", employeeID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&assignments).Error
	return assignments, err
}

func (r *employeeAssignmentRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]entity.EmployeeAssignment, error) {
	var assignments []entity.EmployeeAssignment
	query := r.db.WithContext(ctx).Where("project_id = ?", projectID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&assignments).Error
	return assignments, err
}

func (r *employeeAssignmentRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.EmployeeAssignment, error) {
	var assignments []entity.EmployeeAssignment
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&assignments).Error
	return assignments, err
}

func (r *employeeAssignmentRepository) FindByStatus(ctx context.Context, tenantID uuid.UUID, status entity.AssignmentStatus, limit, offset int) ([]entity.EmployeeAssignment, error) {
	var assignments []entity.EmployeeAssignment
	query := r.db.WithContext(ctx).Where("tenant_id = ? AND status = ?", tenantID, status)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&assignments).Error
	return assignments, err
}

func (r *employeeAssignmentRepository) FindActiveByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]entity.EmployeeAssignment, error) {
	var assignments []entity.EmployeeAssignment
	err := r.db.WithContext(ctx).
		Where("employee_id = ? AND status = ?", employeeID, entity.AssignmentStatusActive).
		Find(&assignments).Error
	return assignments, err
}

func (r *employeeAssignmentRepository) Update(ctx context.Context, assignment *entity.EmployeeAssignment) error {
	return r.db.WithContext(ctx).Save(assignment).Error
}

func (r *employeeAssignmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.EmployeeAssignment{}, "id = ?", id).Error
}

func (r *employeeAssignmentRepository) CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.EmployeeAssignment{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}
