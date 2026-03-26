package request

type PayrollRequest struct {
	EmployeeID  string  `json:"employee_id" validate:"required,uuid"`
	PeriodStart string  `json:"period_start" validate:"required,datetime=2006-01-02"`
	PeriodEnd   string  `json:"period_end" validate:"required,datetime=2006-01-02"`
	BasicSalary float64 `json:"basic_salary" validate:"required,gte=0"`
	Allowances  float64 `json:"allowances" validate:"gte=0"`
	Deductions  float64 `json:"deductions" validate:"gte=0"`
	Tax         float64 `json:"tax" validate:"gte=0"`
	NetSalary   float64 `json:"net_salary" validate:"gte=0"`
	Notes       string  `json:"notes"`
}

type PayrollUpdateRequest struct {
	BasicSalary float64 `json:"basic_salary" validate:"gte=0"`
	Allowances  float64 `json:"allowances" validate:"gte=0"`
	Deductions  float64 `json:"deductions" validate:"gte=0"`
	Tax         float64 `json:"tax" validate:"gte=0"`
	NetSalary   float64 `json:"net_salary" validate:"gte=0"`
	Notes       string  `json:"notes"`
}

type ProcessPayrollRequest struct {
	PayrollIDs []string `json:"payroll_ids" validate:"required,dive,uuid"`
}

type MarkAsPaidRequest struct {
	PaymentDate   string `json:"payment_date" validate:"required,datetime=2006-01-02"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	TransactionID string `json:"transaction_id"`
}

type PayrollItemRequest struct {
	ItemType string  `json:"item_type" validate:"required"`
	Amount   float64 `json:"amount" validate:"required,gte=0"`
	Notes    string  `json:"notes"`
}
