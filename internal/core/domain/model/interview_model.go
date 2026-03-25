package model

import "time"

type Interview struct {
	ID                int64           `gorm:"id"`
	TenantID          int64           `gorm:"tenant_id"`
	CandidateID       int64           `gorm:"candidate_id"`
	ManpowerRequestID int64           `gorm:"manpower_request_id"`
	InterviewType     string          `gorm:"interview_type"`
	ScheduledAt       time.Time       `gorm:"scheduled_at"`
	DurationMinutes   int             `gorm:"duration_minutes"`
	Location          string          `gorm:"location"`
	MeetingLink       string          `gorm:"meeting_link"`
	Status            string          `gorm:"status"`
	Feedback          string          `gorm:"feedback"`
	Rating            int             `gorm:"rating"`
	CompletedAt       time.Time       `gorm:"completed_at"`
	CancelledAt       time.Time       `gorm:"cancelled_at"`
	CancelReason      string          `gorm:"cancel_reason"`
	Tenant            Tenant          `gorm:"foreignkey:TenantID"`
	Candidate         Candidate       `gorm:"foreignkey:CandidateID"`
	ManpowerRequest   ManpowerRequest `gorm:"foreignkey:ManpowerRequestID"`
}