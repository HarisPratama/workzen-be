package entity

import "time"

type OfferEntity struct {
	ID                     int64           `json:"id"`
	TenantID               int64           `json:"tenant_id"`
	CandidateApplicationID int64           `json:"candidate_application_id"`
	Position               string          `json:"position"`
	Department             string          `json:"department"`
	EmploymentType         string          `json:"employment_type"`
	BaseSalary             float64         `json:"base_salary"`
	Currency               string          `json:"currency"`
	Bonus                  *float64        `json:"bonus"`
	Benefits               string          `json:"benefits"`
	ProbationPeriodMonths  int             `json:"probation_period_months"`
	NoticePeriodDays       int             `json:"notice_period_days"`
	StartDate              string          `json:"start_date"`
	ExpiryDate             string          `json:"expiry_date"`
	Status                 string          `json:"status"`
	SentAt                 time.Time       `json:"sent_at"`
	RespondedAt            time.Time       `json:"responded_at"`
	Notes                  string          `json:"notes"`
	Terms                  string          `json:"terms"`
	NegotiationCounter     *float64        `json:"negotiation_counter"`
	NegotiationNotes       string          `json:"negotiation_notes"`
	Tenant                 TenantEntity    `json:"tenant"`
	CandidateApplication   CandidateApplicationEntity `json:"candidate_application"`
}

type OfferQueryString struct {
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

type OfferEntityRequest struct {
	CandidateApplicationID int64  `json:"candidate_application_id" validate:"required"`
	Position              string `json:"position" validate:"required"`
	Department            string `json:"department"`
	EmploymentType        string `json:"employment_type"`
	BaseSalary            float64 `json:"base_salary" validate:"required"`
	Currency              string  `json:"currency"`
	Bonus                 *float64 `json:"bonus"`
	Benefits              string  `json:"benefits"`
	ProbationPeriodMonths int     `json:"probation_period_months"`
	NoticePeriodDays      int     `json:"notice_period_days"`
	StartDate             string  `json:"start_date"`
	ExpiryDate            string  `json:"expiry_date"`
}

type OfferUpdateRequest struct {
	Status         string    `json:"status"`
	Feedback       string    `json:"feedback"`
	NegotiatedSalary *float64 `json:"negotiated_salary"`
	StartDate      string    `json:"start_date"`
}