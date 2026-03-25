package entity

import "time"

type EmployeeAssignmentEntity struct {
	ID                int64           `json:"id"`
	TenantID          int64           `json:"tenant_id"`
	EmployeeID        int64           `json:"employee_id"`
	ClientID          int64           `json:"client_id"`
	ProjectID         int64           `json:"project_id"`
	DepartmentID      int64           `json:"department_id"`
	AssignmentType    string          `json:"assignment_type"`
	StartDate         time.Time       `json:"start_date"`
	EndDate           time.Time       `json:"end_date"`
	ExpectedEndDate   time.Time       `json:"expected_end_date"`
	Status            string          `json:"status"`
	Role              string          `json:"role"`
	Position          string          `json:"position"`
	Location          string          `json:"location"`
	RemoteType        string          `json:"remote_type"`
	BillingRate       float64         `json:"billing_rate"`
	CostRate          float64         `json:"cost_rate"`
	Currency          string          `json:"currency"`
	HoursPerWeek      int             `json:"hours_per_week"`
	Notes             string          `json:"notes"`
	Reason            string          `json:"reason"`
	HandoverNotes     string          `json:"handover_notes"`
	ApprovedByID      int64           `json:"approved_by_id"`
	ApprovedAt        time.Time       `json:"approved_at"`
	TerminatedByID    int64           `json:"terminated_by_id"`
	TerminatedAt      time.Time       `json:"terminated_at"`
	TerminationReason string          `json:"termination_reason"`
	Tenant            TenantEntity    `json:"tenant"`
	Employee          EmployeeEntity  `json:"employee"`
	Client            ClientEntity    `json:"client"`
	Project           ProjectEntity   `json:"project"`
	Department        DepartmentEntity `json:"department"`
}

type EmployeeAssignmentQueryString struct {
	Limit           int
	Page            int
	OrderBy         string
	OrderType       string
	Search          string
	Status          string
	EmployeeID      int64
	StartDate       string
	EndDate         string
	AssignmentType  string
}

type EmployeeAssignmentEntityRequest struct {
	EmployeeID        int64  `json:"employee_id" validate:"required"`
	ClientID          int64  `json:"client_id"`
	ProjectID         int64  `json:"project_id"`
	DepartmentID      int64  `json:"department_id"`
	AssignmentType    string `json:"assignment_type" validate:"required"`
	StartDate         string `json:"start_date" validate:"required"`
	EndDate           string `json:"end_date"`
	ExpectedEndDate  string `json:"expected_end_date"`
	Role              string `json:"role"`
	Position          string `json:"position"`
	Location          string `json:"location"`
	RemoteType        string `json:"remote_type"`
	BillingRate       float64 `json:"billing_rate"`
	CostRate          float64 `json:"cost_rate"`
	Currency          string  `json:"currency"`
	HoursPerWeek      int     `json:"hours_per_week"`
	Notes             string  `json:"notes"`
}

type EmployeeAssignmentUpdateRequest struct {
	Status           string `json:"status"`
	EndDate          string `json:"end_date"`
	Notes            string `json:"notes"`
	TerminationReason string `json:"termination_reason"`
}