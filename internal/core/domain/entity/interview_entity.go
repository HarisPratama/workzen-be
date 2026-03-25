package entity

import "time"

type InterviewEntity struct {
	ID                int64           `json:"id"`
	TenantID          int64           `json:"tenant_id"`
	CandidateID       int64           `json:"candidate_id"`
	ManpowerRequestID int64           `json:"manpower_request_id"`
	Status            string          `json:"status"`
	InterviewType     string          `json:"interview_type"`
	ScheduledAt       time.Time       `json:"scheduled_at"`
	DurationMinutes   int             `json:"duration_minutes"`
	Location          string          `json:"location"`
	MeetingLink       string          `json:"meeting_link"`
	Feedback          string          `json:"feedback"`
	Rating            int             `json:"rating"`
	CompletedAt       time.Time       `json:"completed_at"`
	CancelledAt       time.Time       `json:"cancelled_at"`
	CancelReason      string          `json:"cancel_reason"`
	Tenant            TenantEntity    `json:"tenant"`
	Candidate         CandidateEntity `json:"candidate"`
	ManpowerRequest   ManpowerReqEntity `json:"manpower_request"`
}

type InterviewQueryString struct {
	Limit      int
	Page       int
	OrderBy    string
	OrderType  string
	Search     string
	Status     string
	CandidateID int64
	StartDate  string
	EndDate    string
}

type InterviewEntityRequest struct {
	CandidateID       int64  `json:"candidate_id" validate:"required"`
	ManpowerRequestID int64  `json:"manpower_request_id" validate:"required"`
	InterviewType     string `json:"interview_type" validate:"required"`
	ScheduledAt       string `json:"scheduled_at" validate:"required"`
	DurationMinutes   int    `json:"duration_minutes" validate:"required"`
	Location          string `json:"location"`
	MeetingLink       string `json:"meeting_link"`
}

type InterviewUpdateRequest struct {
	Status    string  `json:"status"`
	Feedback  string  `json:"feedback"`
	Rating    int     `json:"rating"`
	Location  string  `json:"location"`
	ScheduledAt string `json:"scheduled_at"`
}