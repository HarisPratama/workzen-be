package entity

import (
	"time"
)

type PayrollStatus string

const (
	PayrollStatusDraft     PayrollStatus = "DRAFT"
	PayrollStatusProcessed PayrollStatus = "PROCESSED"
	PayrollStatusPaid      PayrollStatus = "PAID"
	PayrollStatusCancelled PayrollStatus = "CANCELLED"
)

type Payroll struct {
	ID          int64         `json:"id" gorm:"primary_key;autoIncrement"`
	TenantID    int64         `json:"tenant_id" gorm:"not null;index"`
	EmployeeID  int64         `json:"employee_id" gorm:"not null;index"`
	PeriodStart time.Time     `json:"period_start" gorm:"not null"`
	PeriodEnd   time.Time     `json:"period_end" gorm:"not null"`
	BasicSalary float64       `json:"basic_salary" gorm:"type:decimal(15,2);not null"`
	Allowances  float64       `json:"allowances" gorm:"type:decimal(15,2);default:0"`
	Deductions  float64       `json:"deductions" gorm:"type:decimal(15,2);default:0"`
	Tax         float64       `json:"tax" gorm:"type:decimal(15,2);default:0"`
	NetSalary   float64       `json:"net_salary" gorm:"type:decimal(15,2);not null"`
	Status      PayrollStatus `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	PaidAt      *time.Time    `json:"paid_at,omitempty"`
	Notes       string        `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time    `json:"deleted_at,omitempty" gorm:"index"`
}

func (Payroll) TableName() string {
	return "payrolls"
}

type PayrollItem struct {
	ID          int64     `json:"id" gorm:"primary_key;autoIncrement"`
	PayrollID   int64     `json:"payroll_id" gorm:"not null;index"`
	Type        string    `json:"type" gorm:"type:varchar(50);not null"` // e.g., "OVERTIME", "BONUS", "DEDUCTION"
	Description string    `json:"description" gorm:"type:varchar(255)"`
	Amount      float64   `json:"amount" gorm:"type:decimal(15,2);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Payroll *Payroll `json:"payroll,omitempty" gorm:"foreignKey:PayrollID"`
}

func (PayrollItem) TableName() string {
	return "payroll_details"
}

type PayrollCalculation struct {
	BasicSalary float64 `json:"basic_salary"`
	Allowances  float64 `json:"allowances"`
	Deductions  float64 `json:"deductions"`
	Tax         float64 `json:"tax"`
	NetSalary   float64 `json:"net_salary"`
	TaxDetails  TaxInfo `json:"tax_details,omitempty"`
}

type TaxInfo struct {
	TaxableIncome float64 `json:"taxable_income"`
	TaxRate       float64 `json:"tax_rate"`
	TaxAmount     float64 `json:"tax_amount"`
}

// CalculateNetSalary calculates the net salary based on components
func (p *Payroll) CalculateNetSalary() {
	p.NetSalary = p.BasicSalary + p.Allowances - p.Deductions - p.Tax
}

// IsEditable checks if the payroll can be edited
func (p *Payroll) IsEditable() bool {
	return p.Status == PayrollStatusDraft
}

// CanProcess checks if the payroll can be processed
func (p *Payroll) CanProcess() bool {
	return p.Status == PayrollStatusDraft
}

// CanPay checks if the payroll can be marked as paid
func (p *Payroll) CanPay() bool {
	return p.Status == PayrollStatusProcessed
}

type PayrollSummary struct {
	TotalPayrolls    int64     `json:"total_payrolls"`
	TotalBasicSalary float64   `json:"total_basic_salary"`
	TotalNetSalary   float64   `json:"total_net_salary"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
}
