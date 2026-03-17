package entity

import (
	"time"

	"github.com/google/uuid"
)

type InterviewStatus string

const (
	InterviewStatusScheduled   InterviewStatus = "SCHEDULED"
	InterviewStatusCompleted   InterviewStatus = "COMPLETED"
	InterviewStatusCancelled   InterviewStatus = "CANCELLED"
	InterviewStatusNoShow      InterviewStatus = "NO_SHOW"
	InterviewStatusRescheduled InterviewStatus = "RESCHEDULED"
)

type InterviewType string

const (
	InterviewTypePhone      InterviewType = "PHONE"
	InterviewTypeVideo      InterviewType = "VIDEO"
	InterviewTypeInPerson   InterviewType = "IN_PERSON"
	InterviewTypeTechnical  InterviewType = "TECHNICAL"
	InterviewTypeHR         InterviewType = "HR"
	InterviewTypePanel      InterviewType = "PANEL"
)

type InterviewResult string

const (
	InterviewResultPass       InterviewResult = "PASS"
	InterviewResultFail         InterviewResult = "FAIL"
	InterviewResultPending      InterviewResult = "PENDING"
	InterviewResultOnHold       InterviewResult = "ON_HOLD"
	InterviewResultRecommended  InterviewResult = "RECOMMENDED"
	InterviewResultStrongReject InterviewResult = "STRONG_REJECT"
)

// Interview represents a job interview session
type Interview struct {
	ID                    uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID              uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	CandidateApplicationID uuid.UUID      `json:"candidate_application_id" gorm:"type:uuid;not null;index"`
	InterviewerID         *uuid.UUID      `json:"interviewer_id,omitempty" gorm:"type:uuid;index"`
	ScheduledAt           time.Time       `json:"scheduled_at" gorm:"not null;index"`
	DurationMinutes       int             `json:"duration_minutes" gorm:"default:60"`
	Type                  InterviewType   `json:"type" gorm:"type:varchar(20);default:'IN_PERSON'"`
	Status                InterviewStatus `json:"status" gorm:"type:varchar(20);default:'SCHEDULED'"`
	Result                InterviewResult `json:"result" gorm:"type:varchar(20);default:'PENDING'"`
	Location              *string         `json:"location,omitempty"`
	MeetingLink           *string         `json:"meeting_link,omitempty"`
	Feedback              *string         `json:"feedback,omitempty"`
	CandidateStrengths    *string         `json:"candidate_strengths,omitempty"`
	CandidateWeaknesses   *string         `json:"candidate_weaknesses,omitempty"`
	RecommendationNotes   *string         `json:"recommendation_notes,omitempty"`
	Score                 *float64        `json:"score,omitempty" gorm:"type:decimal(3,1)"`
	CompletedAt           *time.Time      `json:"completed_at,omitempty"`
	CreatedAt             time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt             *time.Time      `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	CandidateApplication *CandidateApplication `json:"candidate_application,omitempty" gorm:"foreignKey:CandidateApplicationID"`
	Interviewer          *Employee             `json:"interviewer,omitempty" gorm:"foreignKey:InterviewerID"`
	Tenant               *Tenant               `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// InterviewQueryString represents query parameters for filtering interviews
type InterviewQueryString struct {
	Page         int    `form:"page" default:"1"`
	Limit        int    `form:"limit" default:"10"`
	Status       string `form:"status"`
	Type         string `form:"type"`
	Result       string `form:"result"`
	StartDate    string `form:"start_date"`
	EndDate      string `form:"end_date"`
	InterviewerID string `form:"interviewer_id"`
}

// InterviewMetrics holds aggregated interview data
type InterviewMetrics struct {
	TotalScheduled       int     `json:"total_scheduled"`
	TotalCompleted       int     `json:"total_completed"`
	TotalCancelled       int     `json:"total_cancelled"`
	TotalNoShow          int     `json:"total_no_show"`
	PassRate             float64 `json:"pass_rate"`
	AverageScore         float64 `json:"average_score"`
	RecommendationRate   float64 `json:"recommendation_rate"`
	AverageInterviewDuration int `json:"average_interview_duration_minutes"`
}

// CanCancel checks if the interview can be cancelled
func (i *Interview) CanCancel() bool {
	return i.Status == InterviewStatusScheduled || i.Status == InterviewStatusRescheduled
}

// CanReschedule checks if the interview can be rescheduled
func (i *Interview) CanReschedule() bool {
	return i.Status == InterviewStatusScheduled || i.Status == InterviewStatusRescheduled
}

// CanComplete checks if the interview can be marked as completed
func (i *Interview) CanComplete() bool {
	return i.Status == InterviewStatusScheduled || i.Status == InterviewStatusRescheduled
}

// IsOverdue checks if the interview is overdue (scheduled time passed but not completed)
func (i *Interview) IsOverdue() bool {
	return i.Status == InterviewStatusScheduled && i.ScheduledAt.Before(time.Now())
}

// CalculateDuration calculates the actual duration of the interview if completed
func (i *Interview) CalculateDuration() int {
	if i.CompletedAt != nil {
		return int(i.CompletedAt.Sub(i.ScheduledAt).Minutes())
	}
	return i.DurationMinutes
}