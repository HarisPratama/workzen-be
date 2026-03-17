package response

type EmployeeAssignmentResponse struct {
	ID                   string  `json:"id"`
	TenantID             string  `json:"tenant_id"`
	EmployeeID           string  `json:"employee_id"`
	ProjectID            string  `json:"project_id"`
	Role                 string  `json:"role"`
	StartDate            string  `json:"start_date"`
	EndDate              *string `json:"end_date,omitempty"`
	AllocationPercentage int     `json:"allocation_percentage"`
	Billable             bool    `json:"billable"`
	Status               string  `json:"status"`
	Notes                string  `json:"notes,omitempty"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

type EmployeeAssignmentListResponse struct {
	Assignments []EmployeeAssignmentResponse `json:"assignments"`
	Meta        PaginationMeta                 `json:"meta"`
}

type EmployeeAssignmentSummaryResponse struct {
	TotalAssignments      int     `json:"total_assignments"`
	ActiveAssignments     int     `json:"active_assignments"`
	CompletedAssignments  int     `json:"completed_assignments"`
	PendingAssignments    int     `json:"pending_assignments"`
	CancelledAssignments  int     `json:"cancelled_assignments"`
	TotalBillableProjects int     `json:"total_billable_projects"`
	TotalAllocationHours  float64 `json:"total_allocation_hours"`
	AverageAllocationPercentage float64 `json:"average_allocation_percentage"`
}

type ActiveAssignmentsResponse struct {
	Assignments []EmployeeAssignmentResponse `json:"assignments"`
}

type ProjectUtilizationResponse struct {
	ProjectID            string  `json:"project_id"`
	ProjectName          string  `json:"project_name"`
	TotalAssignments     int     `json:"total_assignments"`
	ActiveAssignments    int     `json:"active_assignments"`
	TotalAllocatedHours  float64 `json:"total_allocated_hours"`
	BillableHours        float64 `json:"billable_hours"`
	UtilizationRate      float64 `json:"utilization_rate"`
}

type EmployeeUtilizationResponse struct {
	EmployeeID           string                         `json:"employee_id"`
	EmployeeName         string                         `json:"employee_name"`
	CurrentAssignments   []EmployeeAssignmentResponse   `json:"current_assignments"`
	TotalAllocationPercentage int                       `json:"total_allocation_percentage"`
	AvailablePercentage  int                            `json:"available_percentage"`
	BillableProjects     int                            `json:"billable_projects"`
}

type AssignmentHistoryResponse struct {
	Assignments []EmployeeAssignmentResponse `json:"assignments"`
}

type AssignmentOverlapCheckResponse struct {
	HasOverlap      bool   `json:"has_overlap"`
	OverlappingWith *string `json:"overlapping_with,omitempty"`
	Message         string `json:"message"`
}
