package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewStatus string

const (
	InterviewStatusScheduled   InterviewStatus = "scheduled"
	InterviewStatusCompleted   InterviewStatus = "completed"
	InterviewStatusCancelled   InterviewStatus = "cancelled"
	InterviewStatusNoShow      InterviewStatus = "no_show"
	InterviewStatusRescheduled InterviewStatus = "rescheduled"
)

type InterviewType string

const (
	InterviewTypePhone      InterviewType = "phone"
	InterviewTypeVideo      InterviewType = "video"
	InterviewTypeInPerson   InterviewType = "in_person"
	InterviewTypeTechnical  InterviewType = "technical"
	InterviewTypeHR         InterviewType = "hr"
	InterviewTypeFinal      InterviewType = "final"
)

type Interview struct {
	ID                     uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
	DeletedAt              gorm.DeletedAt  `gorm:"index" json:"-"`
	TenantID               uuid.UUID       `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CandidateApplicationID uuid.UUID       `gorm:"type:uuid;not null;index" json:"candidate_application_id"`
	InterviewerID          uuid.UUID       `gorm:"type:uuid;index" json:"interviewer_id"`
	InterviewType          InterviewType   `gorm:"type:varchar(50);not null" json:"interview_type"`
	ScheduledAt            time.Time       `gorm:"not null" json:"scheduled_at"`
	DurationMinutes        int             `gorm:"not null;default:60" json:"duration_minutes"`
	Location               *string         `json:"location"`
	MeetingLink            *string         `json:"meeting_link"`
	Status                 InterviewStatus `gorm:"type:varchar(50);default:'scheduled'" json:"status"`
	Notes                  *string         `json:"notes"`
	Feedback               *string         `json:"feedback"`
	Rating                 *int            `json:"rating"`
	CompletedAt            *time.Time      `json:"completed_at"`
	CancelledAt            *time.Time      `json:"cancelled_at"`
	CancelReason           *string         `json:"cancel_reason"`
}