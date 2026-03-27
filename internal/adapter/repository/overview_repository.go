package repository

import (
	"context"

	"gorm.io/gorm"
)

type OverviewData struct {
	TotalEmployees             int64 `json:"total_employees"`
	TotalClients               int64 `json:"total_clients"`
	TotalCandidates            int64 `json:"total_candidates"`
	TotalManpowerRequests      int64 `json:"total_manpower_requests"`
	TotalCandidateApplications int64 `json:"total_candidate_applications"`
	TotalInterviews            int64 `json:"total_interviews"`
	TotalOffers                int64 `json:"total_offers"`
	TotalAssignments           int64 `json:"total_assignments"`
}

type OverviewRepository interface {
	GetOverviewByTenant(ctx context.Context, tenantID int64) (*OverviewData, error)
}

type overviewRepository struct {
	db *gorm.DB
}

func (r *overviewRepository) GetOverviewByTenant(ctx context.Context, tenantID int64) (*OverviewData, error) {
	var data OverviewData

	db := r.db.WithContext(ctx)

	if err := db.Table("employees").Where("tenant_id = ?", tenantID).Count(&data.TotalEmployees).Error; err != nil {
		return nil, err
	}

	if err := db.Table("clients").Where("tenant_id = ?", tenantID).Count(&data.TotalClients).Error; err != nil {
		return nil, err
	}

	if err := db.Table("candidates").Where("tenant_id = ?", tenantID).Count(&data.TotalCandidates).Error; err != nil {
		return nil, err
	}

	if err := db.Table("manpower_requests").Where("tenant_id = ?", tenantID).Count(&data.TotalManpowerRequests).Error; err != nil {
		return nil, err
	}

	if err := db.Table("candidate_applications").Where("tenant_id = ?", tenantID).Count(&data.TotalCandidateApplications).Error; err != nil {
		return nil, err
	}

	if err := db.Table("interviews").Where("tenant_id = ?", tenantID).Count(&data.TotalInterviews).Error; err != nil {
		return nil, err
	}

	if err := db.Table("offers").Where("tenant_id = ?", tenantID).Count(&data.TotalOffers).Error; err != nil {
		return nil, err
	}

	if err := db.Table("employee_assignments").Where("tenant_id = ?", tenantID).Count(&data.TotalAssignments).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func NewOverviewRepository(db *gorm.DB) OverviewRepository {
	return &overviewRepository{db: db}
}
