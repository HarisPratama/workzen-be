package entity

import "time"

type ManpowerReqEntity struct {
	ID             int64
	TenantID       int64
	ClientID       int64
	Position       string
	RequiredCount  int
	SalaryMin      float64
	SalaryMax      float64
	WorkLocation   string
	JobDescription string
	DeadlineDate   time.Time
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Tenant         TenantEntity
	Client         ClientEntity
}

type ManpowerReqQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
