package model

import "time"

type CandidateApplication struct {
	ID                int64           `gorm:"id"`
	TenantID          int64           `gorm:"tenant_id"`
	CandidateID       int64           `gorm:"candidate_id"`
	ManpowerRequestID int64           `gorm:"manpower_request_id"`
	Status            string          `gorm:"status"`
	AppliedAt         time.Time       `gorm:"applied_at"`
	Tenant            Tenant          `gorm:"foreignkey:TenantID"`
	Candidate         Candidate       `gorm:"foreignkey:CandidateID"`
	ManpowerRequest   ManpowerRequest `gorm:"foreignkey:ManpowerRequestID"`
}
