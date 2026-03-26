package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"workzen-be/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeeRepository interface {
	GetEmployees(ctx context.Context, query entity.EmployeeQueryString) ([]entity.EmployeeEntity, int64, int64, error)
	GetEmployee(ctx context.Context, id int64) (*entity.EmployeeEntity, error)
	GetEmployeesByTenant(ctx context.Context, tenantId int64, query entity.EmployeeQueryString) ([]entity.EmployeeEntity, int64, int64, error)
	GetDetailEmployeeByTenant(ctx context.Context, tenantID int64, employeeID int64) (*entity.EmployeeEntity, error)
	CreateEmployee(ctx context.Context, employee *entity.EmployeeEntity) error
	UpdateEmployee(ctx context.Context, employee *entity.EmployeeEntity) error
	DeleteEmployee(ctx context.Context, tenantID int64, employeeID int64) error
}

type employeeRepository struct {
	db *gorm.DB
}

func (e *employeeRepository) CreateEmployee(ctx context.Context, employee *entity.EmployeeEntity) error {
	result := e.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "citizen_id"},
			},
			DoNothing: true,
		}).
		Create(employee)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("citizen ID already exists")
	}

	return nil
}

func (e *employeeRepository) GetEmployee(ctx context.Context, id int64) (*entity.EmployeeEntity, error) {
	var employee entity.EmployeeEntity

	err := e.db.WithContext(ctx).
		Where("id = ?", id).
		Preload(clause.Associations).
		First(&employee).Error

	if err != nil {
		code = "[REPOSITORY] GetEmployee - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &employee, nil
}

func (e *employeeRepository) GetDetailEmployeeByTenant(ctx context.Context, tenantID int64, employeeID int64) (*entity.EmployeeEntity, error) {
	var employee entity.EmployeeEntity

	err := e.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Preload(clause.Associations).
		First(&employee, employeeID).Error

	if err != nil {
		return nil, err
	}

	return &employee, nil
}

func (e *employeeRepository) UpdateEmployee(ctx context.Context, employee *entity.EmployeeEntity) error {
	err := e.db.WithContext(ctx).Save(employee).Error
	if err != nil {
		code = "[REPOSITORY] UpdateEmployee - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (e *employeeRepository) DeleteEmployee(ctx context.Context, tenantID int64, employeeID int64) error {
	err := e.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&entity.EmployeeEntity{}, employeeID).Error
	if err != nil {
		code = "[REPOSITORY] DeleteEmployee - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (e *employeeRepository) GetEmployees(ctx context.Context, query entity.EmployeeQueryString) ([]entity.EmployeeEntity, int64, int64, error) {
	var employees []entity.EmployeeEntity
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit
	status := "ACTIVE"
	if query.Status != "" {
		status = query.Status
	}

	sqlMain := e.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("name ILIKE ?", "%"+query.Search+"%").
		Where("status = ?", status)

	err := sqlMain.Model(&entity.EmployeeEntity{}).Count(&countData).Error
	if err != nil {
		code = "[REPOSITORY] GetEmployees - 1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))

	err = sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&employees).Error
	if err != nil {
		code = "[REPOSITORY] GetEmployees - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	return employees, countData, int64(totalPages), err
}

func (e *employeeRepository) GetEmployeesByTenant(
	ctx context.Context,
	tenantID int64,
	query entity.EmployeeQueryString,
) ([]entity.EmployeeEntity, int64, int64, error) {

	var employees []entity.EmployeeEntity
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	status := "ACTIVE"
	if query.Status != "" {
		status = query.Status
	}

	sqlMain := e.db.WithContext(ctx).
		Preload(clause.Associations).
		Where("tenant_id = ?", tenantID).
		Where("name ILIKE ?", "%"+query.Search+"%").
		Where("status = ?", status)

	if err := sqlMain.Model(&entity.EmployeeEntity{}).Count(&countData).Error; err != nil {
		log.Errorw("[REPOSITORY] GetEmployeesByTenant - 1")
		return nil, 0, 0, err
	}

	totalPages := int64(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.
		Order(order).
		Limit(query.Limit).
		Offset(offset).
		Find(&employees).Error; err != nil {
		log.Errorw("[REPOSITORY] GetEmployeesByTenant - 2", "tenant_id", tenantID, "error", err)
		return nil, 0, 0, err
	}

	return employees, countData, totalPages, nil
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}
