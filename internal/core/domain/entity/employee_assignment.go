package entity

import (
	"time"

	"github.com/google/uuid"
)

type AssignmentStatus string

const (
	AssignmentStatusActive    AssignmentStatus = "ACTIVE"
	AssignmentStatusCompleted AssignmentStatus = "COMPLETED"
	AssignmentStatusCancelled AssignmentStatus = "CANCELLED"
)

// EmployeeAssignment represents the assignment of an employee to a client/project
type EmployeeAssignment struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID    uuid.UUID        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	EmployeeID  uuid.UUID        `json:"employee_id" gorm:"type:uuid;not null;index"`
	ClientID    uuid.UUID        `json:"client_id" gorm:"type:uuid;not null;index"`
	ProjectID   *uuid.UUID       `json:"project_id,omitempty" gorm:"type:uuid;index"`
	StartDate   time.Time        `json:"start_date" gorm:"not null"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	Status      AssignmentStatus `json:"status" gorm:"type:varchar(20);default:'ACTIVE'"`
	HourlyRate  *float64         `json:"hourly_rate,omitempty" gorm:"type:decimal(15,2)"`
	MonthlyRate *float64         `json:"monthly_rate,omitempty" gorm:"type:decimal(15,2)"`
	Notes       string           `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Client   *Client   `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	Tenant   *Tenant   `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// AssignmentMetrics holds aggregated data about assignments
type AssignmentMetrics struct {
	TotalAssignments      int     `json:"total_assignments"`
	ActiveAssignments     int     `json:"active_assignments"`
	CompletedAssignments  int     `json:"completed_assignments"`
	CancelledAssignments  int     `json:"cancelled_assignments"`
	AverageDurationDays   float64 `json:"average_duration_days"`
}

// AssignmentCostSummary holds cost data for assignments
type AssignmentCostSummary struct {
	TotalEstimatedCost  float64 `json:"total_estimated_cost"`
	TotalHourlyCost     float64 `json:"total_hourly_cost"`
	TotalMonthlyCost    float64 `json:"total_monthly_cost"`
	ActiveCostPerMonth  float64 `json:"active_cost_per_month"`
}

// IsActive checks if the assignment is currently active
func (ea *EmployeeAssignment) IsActive() bool {
	return ea.Status == AssignmentStatusActive &&
		ea.StartDate.Before(time.Now()) &&
		(ea.EndDate == nil || ea.EndDate.After(time.Now()))
}

// CanComplete checks if the assignment can be marked as completed
func (ea *EmployeeAssignment) CanComplete() bool {
	return ea.Status == AssignmentStatusActive
}

// CanCancel checks if the assignment can be cancelled
func (ea *EmployeeAssignment) CanCancel() bool {
	return ea.Status == AssignmentStatusActive
}

// GetDuration returns the duration of the assignment in days
func (ea *EmployeeAssignment) GetDuration() int {
	endDate := time.Now()
	if ea.EndDate != nil {
		endDate = *ea.EndDate
	}
	return int(endDate.Sub(ea.StartDate).Hours() / 24)
}

// CalculateEstimatedCost calculates the estimated total cost
func (ea *EmployeeAssignment) CalculateEstimatedCost() float64 {
	duration := float64(ea.GetDuration())

	if ea.MonthlyRate != nil && *ea.MonthlyRate > 0 {
		return *ea.MonthlyRate * (duration / 30)
	}

	if ea.HourlyRate != nil && *ea.HourlyRate > 0 {
		return *ea.HourlyRate * duration * 8 // 8 hours per day
	}

	return 0
}