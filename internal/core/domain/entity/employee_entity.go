package entity

import "time"

type EmployeeEntity struct {
	ID          int64
	Name        string
	PhoneNumber string
	CitizenID   string
	UserID      int64
	TenantID    int64
	Status      string
	CreatedAt   time.Time
	User        UserEntity
	Tenant      TenantEntity
}

type EmployeeQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
