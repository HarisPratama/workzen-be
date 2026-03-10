package entity

import "time"

type ClientEntity struct {
	ID          int64
	TenantID    int64
	CompanyName string
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tenant      TenantEntity
}

type ClientQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
