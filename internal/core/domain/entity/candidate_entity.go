package entity

import "time"

type CandidateEntity struct {
	ID        int64
	TenantID  int64
	FullName  string
	Email     string
	Phone     string
	CitizenID string
	BirthDate time.Time
	Address   string
	Source    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type CandidateQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
