package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignmentStatus string

const (
	AssignmentStatusActive     AssignmentStatus = "active"
	AssignmentStatusPending    AssignmentStatus = "pending"
	AssignmentStatusCompleted  AssignmentStatus = "completed"
	AssignmentStatusTerminated AssignmentStatus = "terminated"
	AssignmentStatusOnHold     AssignmentStatus = "on_hold"
)

type AssignmentType string

const (
	AssignmentTypeProject    AssignmentType = "project"
	AssignmentTypeClient     AssignmentType = "client"
	AssignmentTypeDepartment AssignmentType = "department"
	AssignmentTypeTraining   AssignmentType = "training"
	AssignmentTypeSecondment AssignmentType = "secondment"
)

type EmployeeAssignment struct {
	ID                uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `gorm:"index" json:"-"`
	TenantID          uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenant_id"`
	EmployeeID        uuid.UUID        `gorm:"type:uuid;not null;index" json:"employee_id"`
	ClientID          *uuid.UUID       `gorm:"type:uuid;index" json:"client_id"`
	ProjectID         *uuid.UUID       `gorm:"type:uuid;index" json:"project_id"`
	DepartmentID      *uuid.UUID       `gorm:"type:uuid;index" json:"department_id"`
	AssignmentType    AssignmentType   `gorm:"type:varchar(50);not null" json:"assignment_type"`
	StartDate         time.Time        `gorm:"not null" json:"start_date"`
	EndDate           *time.Time       `json:"end_date"`
	ExpectedEndDate   *time.Time       `json:"expected_end_date"`
	Status            AssignmentStatus `gorm:"type:varchar(50);default:'active'" json:"status"`
	Role              string           `json:"role"`
	Position          string           `json:"position"`
	Location          *string          `json:"location"`
	RemoteType        string           `gorm:"default:'onsite'" json:"remote_type"`
	BillingRate       *float64         `json:"billing_rate"`
	CostRate          *float64         `json:"cost_rate"`
	Currency          string           `gorm:"size:3;default:'IDR'" json:"currency"`
	HoursPerWeek      *int             `json:"hours_per_week"`
	Notes             *string          `json:"notes"`
	Reason            *string          `json:"reason"`
	HandoverNotes     *string          `json:"handover_notes"`
	ApprovedByID      *uuid.UUID       `gorm:"type:uuid" json:"approved_by_id"`
	ApprovedAt        *time.Time       `json:"approved_at"`
	TerminatedByID    *uuid.UUID       `gorm:"type:uuid" json:"terminated_by_id"`
	TerminatedAt      *time.Time       `json:"terminated_at"`
	TerminationReason *string          `json:"termination_reason"`
}