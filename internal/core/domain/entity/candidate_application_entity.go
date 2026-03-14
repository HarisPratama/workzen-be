package entity

import "time"

type CandidateApplicationEntity struct {
	ID                int64
	TenantID          int64
	CandidateID       int64
	ManpowerRequestID int64
	Status            string
	AppliedAt         time.Time
	Tenant            TenantEntity
	Candidate         CandidateEntity
	ManpowerRequest   ManpowerReqEntity
}

type CandidateApplicationQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
