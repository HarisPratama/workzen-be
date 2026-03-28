package model

import "time"

type Offer struct {
	ID                     int64                `gorm:"id"`
	TenantID               int64                `gorm:"tenant_id"`
	CandidateApplicationID int64                `gorm:"candidate_application_id"`
	Position               string               `gorm:"position"`
	Department             string               `gorm:"department"`
	EmploymentType         string               `gorm:"employment_type"`
	BaseSalary             float64              `gorm:"base_salary"`
	Currency               string               `gorm:"currency"`
	Bonus                  *float64             `gorm:"bonus"`
	Benefits               string               `gorm:"benefits"`
	ProbationPeriodMonths  int                  `gorm:"probation_period_months"`
	NoticePeriodDays       int                  `gorm:"notice_period_days"`
	StartDate              time.Time            `gorm:"start_date"`
	ExpiryDate             time.Time            `gorm:"expiry_date"`
	Status                 string               `gorm:"status"`
	SentAt                 time.Time            `gorm:"sent_at"`
	RespondedAt            time.Time            `gorm:"responded_at"`
	Notes                  string               `gorm:"notes"`
	Terms                  string               `gorm:"terms"`
	NegotiationCounter     *float64             `gorm:"negotiation_counter"`
	Feedback               string               `gorm:"feedback"`
	NegotiationNotes       string               `gorm:"negotiation_notes"`
	Tenant                 Tenant               `gorm:"foreignkey:TenantID"`
	CandidateApplication   CandidateApplication `gorm:"foreignkey:CandidateApplicationID"`
}
