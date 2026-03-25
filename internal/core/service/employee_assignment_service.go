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

type EmployeeAssignmentService interface {
	CreateAssignment(ctx context.Context, assignment entity.EmployeeAssignment) (*entity.EmployeeAssignment, error)
	UpdateAssignment(ctx context.Context, id uint, assignment entity.EmployeeAssignment) (*entity.EmployeeAssignment, error)
	DeleteAssignment(ctx context.Context, id uint) error
	GetAssignmentByID(ctx context.Context, id uint) (*entity.EmployeeAssignment, error)
	GetAssignmentsByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]entity.EmployeeAssignment, int64, error)
	GetAssignmentsByProject(ctx context.Context, projectID uint, page, limit int) ([]entity.EmployeeAssignment, int64, error)
	GetAssignmentsByTenant(ctx context.Context, tenantID uint, page, limit int) ([]entity.EmployeeAssignment, int64, error)
	GetActiveAssignmentsByEmployee(ctx context.Context, employeeID uint) ([]entity.EmployeeAssignment, error)
	StartAssignment(ctx context.Context, id uint) error
	EndAssignment(ctx context.Context, id uint, endDate string, reason string) error
	UpdateAssignmentStatus(ctx context.Context, id uint, status entity.AssignmentStatus, notes string) error
	GetAssignmentUtilization(ctx context.Context, employeeID uint) (float64, error)
}

type employeeAssignmentService struct {
	assignmentRepo repository.EmployeeAssignmentRepository
}

func NewEmployeeAssignmentService(assignmentRepo repository.EmployeeAssignmentRepository) EmployeeAssignmentService {
	return &employeeAssignmentService{
		assignmentRepo: assignmentRepo,
	}
}

func (s *employeeAssignmentService) CreateAssignment(ctx context.Context, assignment entity.EmployeeAssignment) (*entity.EmployeeAssignment, error) {
	if err := validator.ValidateStruct(assignment); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	assignment.Status = entity.AssignmentStatusPending

	if err := s.assignmentRepo.Create(ctx, &assignment); err != nil {
		log.Errorw("failed to create assignment", "error", err)
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	return &assignment, nil
}

func (s *employeeAssignmentService) UpdateAssignment(ctx context.Context, id uint, assignment entity.EmployeeAssignment) (*entity.EmployeeAssignment, error) {
	existing, err := s.assignmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}

	if existing.Status == entity.AssignmentStatusCompleted || existing.Status == entity.AssignmentStatusTerminated {
		return nil, fmt.Errorf("cannot update assignment in status: %s", existing.Status)
	}

	existing.Role = assignment.Role
	existing.Position = assignment.Position
	existing.Location = assignment.Location
	existing.BillingRate = assignment.BillingRate
	existing.CostRate = assignment.CostRate
	existing.HoursPerWeek = assignment.HoursPerWeek
	existing.Notes = assignment.Notes

	if assignment.EndDate != nil {
		existing.EndDate = assignment.EndDate
	}

	if err := s.assignmentRepo.Update(ctx, existing); err != nil {
		log.Errorw("failed to update assignment", "error", err)
		return nil, fmt.Errorf("failed to update assignment: %w", err)
	}

	return existing, nil
}

func (s *employeeAssignmentService) DeleteAssignment(ctx context.Context, id uint) error {
	existing, err := s.assignmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}

	if existing.Status == entity.AssignmentStatusActive {
		return fmt.Errorf("cannot delete active assignment")
	}

	return s.assignmentRepo.Delete(ctx, id)
}

func (s *employeeAssignmentService) GetAssignmentByID(ctx context.Context, id uint) (*entity.EmployeeAssignment, error) {
	assignment, err := s.assignmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}
	return assignment, nil
}

func (s *employeeAssignmentService) GetAssignmentsByEmployee(ctx context.Context, employeeID uint, page, limit int) ([]entity.EmployeeAssignment, int64, error) {
	offset := (page - 1) * limit
	count, err := s.assignmentRepo.CountByCompanyID(ctx, 0) // Fix this - need proper count
	if err != nil {
		count = 0
	}
	assignments, err := s.assignmentRepo.FindByEmployeeID(ctx, employeeID, limit, offset)
	return assignments, count, err
}

func (s *employeeAssignmentService) GetAssignmentsByProject(ctx context.Context, projectID uint, page, limit int) ([]entity.EmployeeAssignment, int64, error) {
	offset := (page - 1) * limit
	count, err := s.assignmentRepo.CountByCompanyID(ctx, 0) // Fix this - need proper count
	if err != nil {
		count = 0
	}
	assignments, err := s.assignmentRepo.FindByProjectID(ctx, projectID, limit, offset)
	return assignments, count, err
}

func (s *employeeAssignmentService) GetAssignmentsByTenant(ctx context.Context, tenantID uint, page, limit int) ([]entity.EmployeeAssignment, int64, error) {
	offset := (page - 1) * limit
	count, err := s.assignmentRepo.CountByCompanyID(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	assignments, err := s.assignmentRepo.FindByCompanyID(ctx, tenantID, limit, offset)
	return assignments, count, err
}

func (s *employeeAssignmentService) GetActiveAssignmentsByEmployee(ctx context.Context, employeeID uint) ([]entity.EmployeeAssignment, error) {
	return s.assignmentRepo.FindActiveByEmployeeID(ctx, employeeID)
}

func (s *employeeAssignmentService) StartAssignment(ctx context.Context, id uint) error {
	assignment, err := s.assignmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}

	if assignment.Status != entity.AssignmentStatusPending {
		return fmt.Errorf("cannot start assignment in status: %s", assignment.Status)
	}

	assignment.Status = entity.AssignmentStatusActive
	now := time.Now()
	assignment.StartDate = now

	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *employeeAssignmentService) EndAssignment(ctx context.Context, id uint, endDate string, reason string) error {
	assignment, err := s.assignmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}

	if assignment.Status != entity.AssignmentStatusActive {
		return fmt.Errorf("cannot end assignment in status: %s", assignment.Status)
	}

	parsedEndDate, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return fmt.Errorf("invalid end date format: %w", err)
	}

	assignment.Status = entity.AssignmentStatusCompleted
	assignment.EndDate = &parsedEndDate
	if reason != "" {
		assignment.TerminationReason = &reason
	}

	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *employeeAssignmentService) UpdateAssignmentStatus(ctx context.Context, id uint, status entity.AssignmentStatus, notes string) error {
	assignment, err := s.assignmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}

	// Validate status transition
	validTransition := false
	switch assignment.Status {
	case entity.AssignmentStatusPending:
		if status == entity.AssignmentStatusActive || status == entity.AssignmentStatusCancelled {
			validTransition = true
		}
	case entity.AssignmentStatusActive:
		if status == entity.AssignmentStatusCompleted || status == entity.AssignmentStatusTerminated {
			validTransition = true
		}
	}

	if !validTransition {
		return fmt.Errorf("invalid status transition from %s to %s", assignment.Status, status)
	}

	assignment.Status = status
	if notes != "" {
		assignment.Notes = &notes
	}

	return s.assignmentRepo.Update(ctx, assignment)
}

func (s *employeeAssignmentService) GetAssignmentUtilization(ctx context.Context, employeeID uint) (float64, error) {
	assignments, err := s.assignmentRepo.FindActiveByEmployeeID(ctx, employeeID)
	if err != nil {
		return 0, err
	}

	var totalUtilization float64
	for _, assignment := range assignments {
		if assignment.HoursPerWeek != nil {
			totalUtilization += float64(*assignment.HoursPerWeek) / 40.0 * 100.0
		}
	}

	if totalUtilization > 100.0 {
		totalUtilization = 100.0
	}

	return totalUtilization, nil
}