package model

import "time"

type EmployeeAssignment struct {
	ID                int64     `gorm:"id"`
	TenantID          int64     `gorm:"tenant_id"`
	EmployeeID        int64     `gorm:"employee_id"`
	ClientID          int64     `gorm:"client_id"`
	ProjectID         int64     `gorm:"project_id"`
	DepartmentID      int64     `gorm:"department_id"`
	AssignmentType    string    `gorm:"assignment_type"`
	StartDate         time.Time `gorm:"start_date"`
	EndDate           time.Time `gorm:"end_date"`
	ExpectedEndDate   time.Time `gorm:"expected_end_date"`
	Status            string    `gorm:"status"`
	Role              string    `gorm:"role"`
	Position          string    `gorm:"position"`
	Location          string    `gorm:"location"`
	RemoteType        string    `gorm:"remote_type"`
	BillingRate       float64   `gorm:"billing_rate"`
	CostRate          float64   `gorm:"cost_rate"`
	Currency          string    `gorm:"currency"`
	HoursPerWeek      int       `gorm:"hours_per_week"`
	Notes             string    `gorm:"notes"`
	Reason            string    `gorm:"reason"`
	HandoverNotes     string    `gorm:"handover_notes"`
	ApprovedByID      int64     `gorm:"approved_by_id"`
	ApprovedAt        time.Time `gorm:"approved_at"`
	TerminatedByID    int64     `gorm:"terminated_by_id"`
	TerminatedAt      time.Time `gorm:"terminated_at"`
	TerminationReason string    `gorm:"termination_reason"`
	Tenant            Tenant    `gorm:"foreignkey:TenantID"`
	Employee          Employee  `gorm:"foreignkey:EmployeeID"`
	Client            Client    `gorm:"foreignkey:ClientID"`
}
