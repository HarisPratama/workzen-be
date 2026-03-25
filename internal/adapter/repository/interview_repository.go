package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"workzen-be/internal/core/domain/entity"
)

type interviewRepository struct {
	db *gorm.DB
}

func NewInterviewRepository(db *gorm.DB) InterviewRepository {
	return &interviewRepository{db: db}
}

func (r *interviewRepository) Create(ctx context.Context, interview *entity.Interview) error {
	return r.db.WithContext(ctx).Create(interview).Error
}

func (r *interviewRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Interview, error) {
	var interview entity.Interview
	err := r.db.WithContext(ctx).First(&interview, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &interview, nil
}

func (r *interviewRepository) FindByCandidateApplicationID(ctx context.Context, candidateApplicationID uuid.UUID) ([]entity.Interview, error) {
	var interviews []entity.Interview
	err := r.db.WithContext(ctx).Where("candidate_application_id = ?", candidateApplicationID).Find(&interviews).Error
	return interviews, err
}

func (r *interviewRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]entity.Interview, error) {
	var interviews []entity.Interview
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Find(&interviews).Error
	return interviews, err
}

func (r *interviewRepository) FindByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]entity.Interview, error) {
	var interviews []entity.Interview
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&interviews).Error
	return interviews, err
}

func (r *interviewRepository) Update(ctx context.Context, interview *entity.Interview) error {
	return r.db.WithContext(ctx).Save(interview).Error
}

func (r *interviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Interview{}, "id = ?", id).Error
}

func (r *interviewRepository) CountByTenantID(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Interview{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}
