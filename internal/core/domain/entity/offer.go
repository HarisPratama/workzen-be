package entity

import (
	"time"

	"github.com/google/uuid"
)

type OfferStatus string

const (
	OfferStatusDraft        OfferStatus = "DRAFT"
	OfferStatusPending      OfferStatus = "PENDING"
	OfferStatusSent         OfferStatus = "SENT"
	OfferStatusViewed       OfferStatus = "VIEWED"
	OfferStatusAccepted     OfferStatus = "ACCEPTED"
	OfferStatusRejected     OfferStatus = "REJECTED"
	OfferStatusNegotiating  OfferStatus = "NEGOTIATING"
	OfferStatusWithdrawn    OfferStatus = "WITHDRAWN"
	OfferStatusExpired      OfferStatus = "EXPIRED"
)

type OfferType string

const (
	OfferTypeFullTime    OfferType = "FULL_TIME"
	OfferTypePartTime    OfferType = "PART_TIME"
	OfferTypeContract    OfferType = "CONTRACT"
	OfferTypeInternship  OfferType = "INTERNSHIP"
	OfferTypeFreelance   OfferType = "FREELANCE"
)

// Offer represents a job offer sent to a candidate
type Offer struct {
	ID                     uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID               uuid.UUID       `json:"tenant_id" gorm:"type:uuid;not null;index"`
	CandidateApplicationID uuid.UUID       `json:"candidate_application_id" gorm:"type:uuid;not null;index"`
	JobTitle               string          `json:"job_title" gorm:"type:varchar(255);not null"`
	Department             *string         `json:"department,omitempty"`
	OfferType              OfferType       `json:"offer_type" gorm:"type:varchar(20);default:'FULL_TIME'"`
	EmploymentLevel        *string         `json:"employment_level,omitempty"` // e.g., Junior, Senior, Lead
	BaseSalary             float64         `json:"base_salary" gorm:"type:decimal(15,2);not null"`
	Currency               string          `json:"currency" gorm:"type:varchar(3);default:'USD'"`
	SignOnBonus            *float64        `json:"sign_on_bonus,omitempty" gorm:"type:decimal(15,2)"`
	AnnualBonus          *float64        `json:"annual_bonus,omitempty" gorm:"type:decimal(15,2)"`
	BenefitsPackage      *string         `json:"benefits_package,omitempty" gorm:"type:text"` // JSON or text description
	StockOptions         *float64        `json:"stock_options,omitempty" gorm:"type:decimal(15,2)"`
	VestingSchedule      *string         `json:"vesting_schedule,omitempty"`
	ProbationPeriodDays  *int            `json:"probation_period_days,omitempty"`
	NoticePeriodDays     *int            `json:"notice_period_days,omitempty"`
	PaidTimeOffDays      *int            `json:"paid_time_off_days,omitempty"`
	Status               OfferStatus     `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	SentAt               *time.Time      `json:"sent_at,omitempty"`
	ViewedAt             *time.Time      `json:"viewed_at,omitempty"`
	RespondedAt          *time.Time      `json:"responded_at,omitempty"`
	AcceptedAt           *time.Time      `json:"accepted_at,omitempty"`
	RejectedAt           *time.Time      `json:"rejected_at,omitempty"`
	ExpiryDate           *time.Time      `json:"expiry_date,omitempty"`
	ResponseDeadline     *time.Time      `json:"response_deadline,omitempty"`
	NegotiationCounter   *int            `json:"negotiation_counter,omitempty" gorm:"default:0"`
	RejectionReason      *string         `json:"rejection_reason,omitempty" gorm:"type:text"`
	NegotiationNotes     *string         `json:"negotiation_notes,omitempty" gorm:"type:text"`
	InternalNotes        *string         `json:"internal_notes,omitempty" gorm:"type:text"`
	CreatedBy            uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"`
	SentBy               *uuid.UUID      `json:"sent_by,omitempty"`
	CreatedAt            time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt            *time.Time      `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	CandidateApplication *CandidateApplication `json:"candidate_application,omitempty" gorm:"foreignKey:CandidateApplicationID"`
	Tenant               *Tenant               `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// OfferNegotiation represents a negotiation round for an offer
type OfferNegotiation struct {
	ID               uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OfferID          uuid.UUID     `json:"offer_id" gorm:"type:uuid;not null;index"`
	NegotiationRound int           `json:"negotiation_round" gorm:"not null"`
	ProposedBy       string        `json:"proposed_by" gorm:"type:varchar(20);not null"` // CANDIDATE or EMPLOYER
	PreviousSalary   float64       `json:"previous_salary" gorm:"type:decimal(15,2)"`
	ProposedSalary   float64       `json:"proposed_salary" gorm:"type:decimal(15,2)"`
	ProposedBenefits *string       `json:"proposed_benefits,omitempty" gorm:"type:text"`
	Justification    *string       `json:"justification,omitempty" gorm:"type:text"`
	Response         *string       `json:"response,omitempty" gorm:"type:text"`
	RespondedAt      *time.Time    `json:"responded_at,omitempty"`
	IsAccepted       *bool         `json:"is_accepted,omitempty"`
	CreatedAt        time.Time     `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Offer *Offer `json:"offer,omitempty" gorm:"foreignKey:OfferID"`
}

// OfferMetrics holds aggregated data about offers
type OfferMetrics struct {
	TotalOffers           int     `json:"total_offers"`
	DraftOffers           int     `json:"draft_offers"`
	PendingOffers         int     `json:"pending_offers"`
	SentOffers            int     `json:"sent_offers"`
	AcceptedOffers        int     `json:"accepted_offers"`
	RejectedOffers        int     `json:"rejected_offers"`
	WithdrawnOffers       int     `json:"withdrawn_offers"`
	ExpiredOffers         int     `json: