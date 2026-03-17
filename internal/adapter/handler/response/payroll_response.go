package response

type PayrollResponse struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenant_id"`
	EmployeeID  string  `json:"employee_id"`
	PeriodStart string  `json:"period_start"`
	PeriodEnd   string  `json:"period_end"`
	BasicSalary float64 `json:"basic_salary"`
	Allowances  float64 `json:"allowances"`
	Deductions  float64 `json:"deductions"`
	Tax         float64 `json:"tax"`
	NetSalary   float64 `json:"net_salary"`
	Status      string  `json:"status"`
	PaidAt      *string `json:"paid_at,omitempty"`
	Notes       string  `json:"notes,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type PayrollListResponse struct {
	Payrolls []PayrollResponse `json:"payrolls"`
	Meta     PaginationMeta    `json:"meta"`
}

type PayrollSummaryResponse struct {
	TotalPayrolls      int     `json:"total_payrolls"`
	TotalBasicSalary   float64 `json:"total_basic_salary"`
	TotalAllowances    float64 `json:"total_allowances"`
	TotalDeductions    float64 `json:"total_deductions"`
	TotalTax           float64 `json:"total_tax"`
	TotalNetSalary     float64 `json:"total_net_salary"`
	PaidCount          int     `json:"paid_count"`
	PendingCount       int     `json:"pending_count"`
	ProcessedCount     int     `json:"processed_count"`
}

type PayrollItemResponse struct {
	ID       string  `json:"id"`
	ItemType string  `json:"item_type"`
	Amount   float64 `json:"amount"`
	Notes    string  `json:"notes,omitempty"`
}
