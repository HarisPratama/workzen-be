package model

import "time"

type ManpowerRequest struct {
	ID             int64     `gorm:"id"`
	TenantID       int64     `gorm:"tenant_id"`
	ClientID       int64     `gorm:"client_id"`
	Position       string    `gorm:"position"`
	RequiredCount  int64     `gorm:"required_count"`
	SalaryMin      float64   `gorm:"salary_min"`
	SalaryMax      float64   `gorm:"salary_max"`
	WorkLocation   string    `gorm:"work_location"`
	JobDescription string    `gorm:"job_description"`
	DeadlineDate   time.Time `gorm:"deadline_date"`
	CreatedAt      time.Time `gorm:"created_at"`
	Status         string    `gorm:"status"`
	Tenant         Tenant    `gorm:"foreignkey:TenantID"`
	Client         Client    `gorm:"foreignkey:ClientID"`
}
