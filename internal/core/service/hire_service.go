package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var (
	ErrApplicationNotFound = errors.New("candidate application not found")
	ErrApplicationNotHired = errors.New("candidate application status must be OFFERED to hire")
	ErrCandidateNotFound   = errors.New("candidate data not found")
)

type HireRequest struct {
	StartDate string `json:"start_date" validate:"required"`
}

type HireService interface {
	HireCandidate(ctx context.Context, tenantID int64, applicationID int64, req HireRequest) error
}

type hireService struct {
	db *gorm.DB
}

func (s *hireService) HireCandidate(ctx context.Context, tenantID int64, applicationID int64, req HireRequest) error {
	jakartaTZ, _ := time.LoadLocation("Asia/Jakarta")

	startDate, err := time.ParseInLocation("2006-01-02", req.StartDate, jakartaTZ)
	if err != nil {
		return fmt.Errorf("invalid start_date format, use YYYY-MM-DD")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Get candidate application with relations
		var application model.CandidateApplication
		if err := tx.
			Preload("Candidate").
			Preload("ManpowerRequest").
			Where("tenant_id = ?", tenantID).
			First(&application, applicationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrApplicationNotFound
			}
			return err
		}

		// 2. Create employee from candidate data
		employee := entity.EmployeeEntity{
			Name:        application.Candidate.FullName,
			PhoneNumber: application.Candidate.Phone,
			CitizenID:   application.Candidate.CitizenID,
			TenantID:    tenantID,
			Status:      "ACTIVE",
		}

		if err := tx.Create(&employee).Error; err != nil {
			code := "[SERVICE] HireCandidate - CreateEmployee"
			log.Errorw(code, err)
			return fmt.Errorf("failed to create employee: %w", err)
		}

		// 3. Create assignment (outsource) from manpower request data
		assignment := model.EmployeeAssignment{
			TenantID:       tenantID,
			EmployeeID:     employee.ID,
			ClientID:       application.ManpowerRequest.ClientID,
			AssignmentType: "OUTSOURCE",
			StartDate:      startDate,
			Status:         "ACTIVE",
			Position:       application.ManpowerRequest.Position,
			Location:       application.ManpowerRequest.WorkLocation,
		}

		if err := tx.Omit(
			"ProjectID", "DepartmentID",
			"ApprovedByID", "ApprovedAt",
			"TerminatedByID", "TerminatedAt",
			"EndDate", "ExpectedEndDate",
		).Create(&assignment).Error; err != nil {
			code := "[SERVICE] HireCandidate - CreateAssignment"
			log.Errorw(code, err)
			return fmt.Errorf("failed to create assignment: %w", err)
		}

		// 4. Update candidate application status to HIRED
		if err := tx.Model(&model.CandidateApplication{}).
			Where("id = ?", applicationID).
			Update("status", "HIRED").Error; err != nil {
			code := "[SERVICE] HireCandidate - UpdateStatus"
			log.Errorw(code, err)
			return fmt.Errorf("failed to update application status: %w", err)
		}

		return nil
	})
}

func NewHireService(db *gorm.DB) HireService {
	return &hireService{db: db}
}
