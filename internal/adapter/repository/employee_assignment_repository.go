package repository

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"workzen-be/internal/core/domain/entity"
	"workzen-be/internal/core/domain/model"
)

type EmployeeAssignmentRepository interface {
	GetEmployeeAssignmentsByTenant(ctx context.Context, tenantID int64, query entity.EmployeeAssignmentQueryString) ([]entity.EmployeeAssignmentEntity, int64, int64, error)
	GetEmployeeAssignmentByID(ctx context.Context, id int64) (*entity.EmployeeAssignmentEntity, error)
	CreateEmployeeAssignment(ctx context.Context, req entity.EmployeeAssignmentEntity) error
	UpdateEmployeeAssignment(ctx context.Context, id int64, req entity.EmployeeAssignmentUpdateRequest) error
	DeleteEmployeeAssignment(ctx context.Context, id int64) error
}

type employeeAssignmentRepository struct {
	db *gorm.DB
}

func NewEmployeeAssignmentRepository(db *gorm.DB) EmployeeAssignmentRepository {
	return &employeeAssignmentRepository{db: db}
}

func (r *employeeAssignmentRepository) CreateEmployeeAssignment(ctx context.Context, req entity.EmployeeAssignmentEntity) error {
	modelAssignment := model.EmployeeAssignment{
		TenantID:       req.TenantID,
		EmployeeID:     req.EmployeeID,
		ClientID:       req.ClientID,
		AssignmentType: req.AssignmentType,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         "PENDING",
		Role:           req.Role,
		Position:       req.Position,
		Location:       req.Location,
		RemoteType:     req.RemoteType,
		BillingRate:    req.BillingRate,
		CostRate:       req.CostRate,
		Currency:       req.Currency,
		HoursPerWeek:   req.HoursPerWeek,
		Notes:          req.Notes,
	}

	err := r.db.Create(&modelAssignment).Error
	if err != nil {
		code := "[REPOSITORY] CreateEmployeeAssignment - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (r *employeeAssignmentRepository) GetEmployeeAssignmentsByTenant(ctx context.Context, tenantID int64, query entity.EmployeeAssignmentQueryString) ([]entity.EmployeeAssignmentEntity, int64, int64, error) {
	var modelAssignments []model.EmployeeAssignment
	var countData int64

	sqlMain := r.db.WithContext(ctx).
		Model(&model.EmployeeAssignment{}).
		Preload("Employee").
		Preload("Client").
		Preload("Tenant").
		Where("tenant_id = ?", tenantID)

	if query.Status != "" {
		sqlMain = sqlMain.Where("employee_assignments.status = ?", query.Status)
	}

	if query.EmployeeID != 0 {
		sqlMain = sqlMain.Where("employee_assignments.employee_id = ?", query.EmployeeID)
	}

	if query.AssignmentType != "" {
		sqlMain = sqlMain.Where("employee_assignments.assignment_type = ?", query.AssignmentType)
	}

	if query.StartDate != "" && query.EndDate != "" {
		sqlMain = sqlMain.Where("employee_assignments.start_date BETWEEN ? AND ?", query.StartDate, query.EndDate)
	}

	countQuery := sqlMain.Session(&gorm.Session{})
	if err := countQuery.Count(&countData).Error; err != nil {
		return nil, 0, 0, err
	}

	allowedOrder := map[string]string{
		"start_date": "employee_assignments.start_date",
		"created_at": "employee_assignments.created_at",
	}

	orderBy, ok := allowedOrder[query.OrderBy]
	if !ok {
		orderBy = "employee_assignments.start_date"
	}

	orderType := "DESC"
	if strings.ToUpper(query.OrderType) == "ASC" {
		orderType = "ASC"
	}

	order := fmt.Sprintf("%s %s", orderBy, orderType)

	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	offset := (query.Page - 1) * query.Limit
	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&modelAssignments).Error; err != nil {
		return nil, 0, 0, err
	}

	var result []entity.EmployeeAssignmentEntity
	for _, item := range modelAssignments {
		result = append(result, entity.EmployeeAssignmentEntity{
			ID:             item.ID,
			TenantID:       item.TenantID,
			EmployeeID:     item.EmployeeID,
			ClientID:       item.ClientID,
			AssignmentType: item.AssignmentType,
			StartDate:      item.StartDate,
			EndDate:        item.EndDate,
			Status:         item.Status,
			Role:           item.Role,
			Position:       item.Position,
			Location:       item.Location,
			RemoteType:     item.RemoteType,
			BillingRate:    item.BillingRate,
			CostRate:       item.CostRate,
			Currency:       item.Currency,
			HoursPerWeek:   item.HoursPerWeek,
			Notes:          item.Notes,
			Tenant: entity.TenantEntity{
				ID:          item.Tenant.ID,
				CompanyName: item.Tenant.CompanyName,
			},
			Employee: entity.EmployeeEntity{
				ID:   item.Employee.ID,
				Name: item.Employee.Name,
			},
			Client: entity.ClientEntity{
				ID:          item.Client.ID,
				CompanyName: item.Client.CompanyName,
			},
		})
	}

	return result, countData, totalPages, nil
}

func (r *employeeAssignmentRepository) GetEmployeeAssignmentByID(ctx context.Context, id int64) (*entity.EmployeeAssignmentEntity, error) {
	var modelAssignment model.EmployeeAssignment
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("Client").
		Preload("Tenant").
		First(&modelAssignment, id).Error

	if err != nil {
		return nil, err
	}

	return &entity.EmployeeAssignmentEntity{
		ID:             modelAssignment.ID,
		TenantID:       modelAssignment.TenantID,
		EmployeeID:     modelAssignment.EmployeeID,
		ClientID:       modelAssignment.ClientID,
		AssignmentType: modelAssignment.AssignmentType,
		StartDate:      modelAssignment.StartDate,
		EndDate:        modelAssignment.EndDate,
		Status:         modelAssignment.Status,
		Role:           modelAssignment.Role,
		Position:       modelAssignment.Position,
		Location:       modelAssignment.Location,
		RemoteType:     modelAssignment.RemoteType,
		BillingRate:    modelAssignment.BillingRate,
		CostRate:       modelAssignment.CostRate,
		Currency:       modelAssignment.Currency,
		HoursPerWeek:   modelAssignment.HoursPerWeek,
		Notes:          modelAssignment.Notes,
		Tenant: entity.TenantEntity{
			ID:          modelAssignment.Tenant.ID,
			CompanyName: modelAssignment.Tenant.CompanyName,
		},
		Employee: entity.EmployeeEntity{
			ID:   modelAssignment.Employee.ID,
			Name: modelAssignment.Employee.Name,
		},
		Client: entity.ClientEntity{
			ID:          modelAssignment.Client.ID,
			CompanyName: modelAssignment.Client.CompanyName,
		},
	}, nil
}

func (r *employeeAssignmentRepository) UpdateEmployeeAssignment(ctx context.Context, id int64, req entity.EmployeeAssignmentUpdateRequest) error {
	updates := map[string]interface{}{
		"status": req.Status,
		"notes":  req.Notes,
	}

	if req.TerminationReason != "" {
		updates["termination_reason"] = req.TerminationReason
	}

	err := r.db.WithContext(ctx).
		Model(&model.EmployeeAssignment{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		code := "[REPOSITORY] UpdateEmployeeAssignment - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (r *employeeAssignmentRepository) DeleteEmployeeAssignment(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).
		Delete(&model.EmployeeAssignment{}, id).Error

	if err != nil {
		code := "[REPOSITORY] DeleteEmployeeAssignment - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}
