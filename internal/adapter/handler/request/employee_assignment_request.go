package request

type EmployeeAssignmentRequest struct {
	EmployeeID           string  `json:"employee_id" validate:"required,uuid"`
	ProjectID            string  `json:"project_id" validate:"required,uuid"`
	Role                 string  `json:"role" validate:"required"`
	StartDate            string  `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate              string  `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	AllocationPercentage int     `json:"allocation_percentage" validate:"required,gte=0,lte=100"`
	Billable             bool    `json:"billable"`
	Notes                string  `json:"notes"`
}

type EmployeeAssignmentUpdateRequest struct {
	Role                 string `json:"role"`
	EndDate              string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	AllocationPercentage int    `json:"allocation_percentage" validate:"gte=0,lte=100"`
	Billable             bool   `json:"billable"`
	Notes                string `json:"notes"`
}

type EndAssignmentRequest struct {
	EndDate string `json:"end_date" validate:"required,datetime=2006-01-02"`
	Reason  string `json:"reason"`
}

type UpdateAssignmentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending active completed cancelled"`
	Notes  string `json:"notes"`
}
