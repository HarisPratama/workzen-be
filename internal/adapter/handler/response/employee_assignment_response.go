package response

type EmployeeAssignmentResponse struct {
	ID              int64                          `json:"id"`
	AssignmentType  string                         `json:"assignment_type"`
	StartDate       string                         `json:"start_date"`
	EndDate         string                         `json:"end_date"`
	ExpectedEndDate string                         `json:"expected_end_date"`
	Status          string                         `json:"status"`
	Role            string                         `json:"role"`
	Position        string                         `json:"position"`
	Location        string                         `json:"location"`
	RemoteType      string                         `json:"remote_type"`
	BillingRate     float64                        `json:"billing_rate"`
	CostRate        float64                        `json:"cost_rate"`
	Currency        string                         `json:"currency"`
	HoursPerWeek    int                            `json:"hours_per_week"`
	Notes           string                         `json:"notes"`
	Employee        EmployeeAssignmentEmployeeResp `json:"employee"`
	Client          EmployeeAssignmentClientResp   `json:"client"`
}

type EmployeeAssignmentEmployeeResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type EmployeeAssignmentClientResp struct {
	ID          int64  `json:"id"`
	CompanyName string `json:"company_name"`
}
