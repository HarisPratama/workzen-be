package entity

import "time"

type ManpowerReqEntity struct {
	ID             int64        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID       int64        `json:"tenant_id" gorm:"column:tenant_id;not null;index"`
	ClientID       int64        `json:"client_id" gorm:"column:client_id;not null;index"`
	Position       string       `json:"position" gorm:"column:position;not null"`
	RequiredCount  int          `json:"required_count" gorm:"column:required_count"`
	SalaryMin      float64      `json:"salary_min" gorm:"column:salary_min"`
	SalaryMax      float64      `json:"salary_max" gorm:"column:salary_max"`
	Hired          int          `json:"hired" gorm:"-"`
	WorkLocation   string       `json:"work_location" gorm:"column:work_location"`
	JobDescription string       `json:"job_description" gorm:"column:job_description"`
	DeadlineDate   time.Time    `json:"deadline_date" gorm:"column:deadline_date"`
	Status         string       `json:"status" gorm:"column:status;default:'OPEN'"`
	CreatedAt      time.Time    `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time    `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Tenant         TenantEntity `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	Client         ClientEntity `json:"client,omitempty" gorm:"foreignKey:ClientID"`
}

func (ManpowerReqEntity) TableName() string {
	return "manpower_requests"
}

type ManpowerReqQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
