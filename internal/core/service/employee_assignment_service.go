package service

import (
	"context"
	"time"
	"workzen-be/internal/adapter/repository"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type EmployeeAssignmentService interface {
	GetEmployeeAssignmentsByTenant(ctx context.Context, tenantID int64, query entity.EmployeeAssignmentQueryString) ([]entity.EmployeeAssignmentEntity, int64, int64, error)
	GetEmployeeAssignmentByID(ctx context.Context, id int64) (*entity.EmployeeAssignmentEntity, error)
	CreateEmployeeAssignment(ctx context.Context, req entity.EmployeeAssignmentEntityRequest, tenantID int64) error
	UpdateEmployeeAssignment(ctx context.Context, id int64, req entity.EmployeeAssignmentUpdateRequest) error
	DeleteEmployeeAssignment(ctx context.Context, id int64) error
}

type employeeAssignmentService struct {
	employeeAssignmentRepo repository.EmployeeAssignmentRepository
}

func (s *employeeAssignmentService) GetEmployeeAssignmentsByTenant(ctx context.Context, tenantID int64, query entity.EmployeeAssignmentQueryString) ([]entity.EmployeeAssignmentEntity, int64, int64, error) {
	results, totalData, totalPages, err := s.employeeAssignmentRepo.GetEmployeeAssignmentsByTenant(ctx, tenantID, query)

	if err != nil {
		code := "[SERVICE] GetEmployeeAssignmentsByTenant - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return results, totalData, totalPages, nil
}

func (s *employeeAssignmentService) GetEmployeeAssignmentByID(ctx context.Context, id int64) (*entity.EmployeeAssignmentEntity, error) {
	result, err := s.employeeAssignmentRepo.GetEmployeeAssignmentByID(ctx, id)

	if err != nil {
		code := "[SERVICE] GetEmployeeAssignmentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (s *employeeAssignmentService) CreateEmployeeAssignment(ctx context.Context, req entity.EmployeeAssignmentEntityRequest, tenantID int64) error {
	jakartaTZ, _ := time.LoadLocation("Asia/Jakarta")
	var startDate, endDate, expectedEndDate time.Time
	if req.StartDate != "" {
		startDate, _ = time.ParseInLocation("2006-01-02", req.StartDate, jakartaTZ)
	}
	if req.EndDate != "" {
		endDate, _ = time.ParseInLocation("2006-01-02", req.EndDate, jakartaTZ)
	}
	if req.ExpectedEndDate != "" {
		expectedEndDate, _ = time.ParseInLocation("2006-01-02", req.ExpectedEndDate, jakartaTZ)
	}

	reqEntity := entity.EmployeeAssignmentEntity{
		TenantID:        tenantID,
		EmployeeID:      req.EmployeeID,
		ClientID:        req.ClientID,
		ProjectID:       req.ProjectID,
		DepartmentID:    req.DepartmentID,
		AssignmentType:  req.AssignmentType,
		StartDate:       startDate,
		EndDate:         endDate,
		ExpectedEndDate: expectedEndDate,
		Status:          "PENDING",
		Role:            req.Role,
		Position:        req.Position,
		Location:        req.Location,
		RemoteType:      req.RemoteType,
		BillingRate:     req.BillingRate,
		CostRate:        req.CostRate,
		Currency:        req.Currency,
		HoursPerWeek:    req.HoursPerWeek,
		Notes:           req.Notes,
	}

	err := s.employeeAssignmentRepo.CreateEmployeeAssignment(ctx, reqEntity)

	if err != nil {
		code := "[SERVICE] CreateEmployeeAssignment - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (s *employeeAssignmentService) UpdateEmployeeAssignment(ctx context.Context, id int64, req entity.EmployeeAssignmentUpdateRequest) error {
	err := s.employeeAssignmentRepo.UpdateEmployeeAssignment(ctx, id, req)

	if err != nil {
		code := "[SERVICE] UpdateEmployeeAssignment - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (s *employeeAssignmentService) DeleteEmployeeAssignment(ctx context.Context, id int64) error {
	err := s.employeeAssignmentRepo.DeleteEmployeeAssignment(ctx, id)

	if err != nil {
		code := "[SERVICE] DeleteEmployeeAssignment - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewEmployeeAssignmentService(employeeAssignmentRepo repository.EmployeeAssignmentRepository) EmployeeAssignmentService {
	return &employeeAssignmentService{employeeAssignmentRepo: employeeAssignmentRepo}
}
