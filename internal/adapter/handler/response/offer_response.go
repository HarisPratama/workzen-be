package response

import "time"

type OfferResponse struct {
	ID                     int64     `json:"id"`
	CandidateApplicationID int64     `json:"candidate_application_id"`
	Position               string    `json:"position"`
	Department             string    `json:"department"`
	EmploymentType         string    `json:"employment_type"`
	BaseSalary             float64   `json:"base_salary"`
	Currency               string    `json:"currency"`
	Bonus                  *float64  `json:"bonus"`
	Benefits               string    `json:"benefits"`
	ProbationPeriodMonths  int       `json:"probation_period_months"`
	NoticePeriodDays       int       `json:"notice_period_days"`
	StartDate              string    `json:"start_date"`
	ExpiryDate             string    `json:"expiry_date"`
	Status                 string    `json:"status"`
	SentAt                 time.Time `json:"sent_at"`
	RespondedAt            time.Time `json:"responded_at"`
	Notes                  string    `json:"notes"`
	Terms                  string    `json:"terms"`
}
