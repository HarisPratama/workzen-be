package repository

import (
	"bwanews/internal/core/domain/entity"
	"bwanews/internal/core/domain/model"
	"context"
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeeRepository interface {
	GetEmployees(ctx context.Context, query entity.EmployeeQueryString) ([]entity.EmployeeEntity, int64, int64, error)
	GetEmployee(ctx context.Context, id int64) (*entity.EmployeeEntity, error)
	CreateEmployee(ctx context.Context, req entity.EmployeeEntity, tenantId int64) error
}

type employeeRepository struct {
	db *gorm.DB
}

func (e *employeeRepository) CreateEmployee(ctx context.Context, req entity.EmployeeEntity, tenantId int64) error {
	modelEmployee := model.Employee{
		Name:        req.Name,
		CitizenID:   req.CitizenID,
		PhoneNumber: req.PhoneNumber,
		TenantID:    tenantId,
	}

	err := e.db.Create(&modelEmployee).Error
	if err != nil {
		code = "[REPOSITORY] CreateEmployee - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (e *employeeRepository) GetEmployee(ctx context.Context, id int64) (*entity.EmployeeEntity, error) {
	var modelEmployee model.Employee

	err = e.db.Where("id = ?", id).Preload(clause.Associations).First(&modelEmployee).Error

	if err != nil {
		code = "[REPOSITORY] GetEmployee - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.EmployeeEntity{
		ID:          modelEmployee.ID,
		Name:        modelEmployee.Name,
		CitizenID:   modelEmployee.CitizenID,
		PhoneNumber: modelEmployee.PhoneNumber,
		Status:      modelEmployee.Status,
		CreatedAt:   modelEmployee.CreatedAt,
		User: entity.UserEntity{
			ID:   modelEmployee.User.ID,
			Name: modelEmployee.User.Name,
		},
	}

	return &resp, nil
}

func (e *employeeRepository) GetEmployees(ctx context.Context, query entity.EmployeeQueryString) ([]entity.EmployeeEntity, int64, int64, error) {
	var modelEmployees []model.Employee
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit
	status := "ACTIVE"
	if query.Status != "" {
		status = query.Status
	}

	sqlMain := e.db.Preload(clause.Associations).
		Where("name ilike ?", "%"+query.Search+"%").
		Where("status = ?", status)

	err = sqlMain.Model(&modelEmployees).Count(&countData).Error
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
		Find(&modelEmployees).Error
	if err != nil {
		code = "[REPOSITORY] GetEmployees - 2"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := []entity.EmployeeEntity{}
	for _, val := range modelEmployees {
		resp := entity.EmployeeEntity{
			ID:        val.ID,
			Name:      val.Name,
			Status:    val.Status,
			CitizenID: val.CitizenID,
			CreatedAt: val.CreatedAt,
			User: entity.UserEntity{
				ID:   val.User.ID,
				Name: val.User.Name,
			},
		}

		resps = append(resps, resp)
	}

	return resps, countData, int64(totalPages), err
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}
