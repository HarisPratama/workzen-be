package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferStatus string

const (
	OfferStatusDraft       OfferStatus = "draft"
	OfferStatusPending     OfferStatus = "pending"
	OfferStatusSent        OfferStatus = "sent"
	OfferStatusAccepted    OfferStatus = "accepted"
	OfferStatusRejected    OfferStatus = "rejected"
	OfferStatusNegotiating OfferStatus = "negotiating"
	OfferStatusWithdrawn   OfferStatus = "withdrawn"
	OfferStatusExpired     OfferStatus = "expired"
)

type Offer struct {
	ID                     uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
	TenantID               uuid.UUID     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CandidateApplicationID uuid.UUID     `gorm:"type:uuid;not null;index" json:"candidate_application_id"`
	Position               string        `gorm:"not null" json:"position"`
	Department             string        `json:"department"`
	EmploymentType         string        `json:"employment_type"`
	BaseSalary             float64       `gorm:"not null" json:"base_salary"`
	Currency               string        `gorm:"size:3;default:'IDR'" json:"currency"`
	Bonus                  *float64      `json:"bonus"`
	Benefits               *string       `json:"benefits"`
	ProbationPeriodMonths  *int          `json:"probation_period_months"`
	NoticePeriodDays       *int          `json:"notice_period_days"`
	StartDate              *time.Time    `json:"start_date"`
	ExpiryDate             *time.Time    `json:"expiry_date"`
	Status                 OfferStatus   `gorm:"type:varchar(50);default:'draft'" json:"status"`
	SentAt                 *time.Time    `json:"sent_at"`
	RespondedAt            *time.Time    `json:"responded_at"`
	Notes                  *string       `json:"notes"`
	Terms                  *string       `json:"terms"`
	NegotiationCounter     *float64      `json:"negotiation_counter"`
	NegotiationNotes       *string       `json:"negotiation_notes"`
}