package response

type OfferResponse struct {
	ID                       string  `json:"id"`
	TenantID                 string  `json:"tenant_id"`
	CandidateApplicationID   string  `json:"candidate_application_id"`
	JobTitle                 string  `json:"job_title"`
	Department               string  `json:"department"`
	OfferType                string  `json:"offer_type"`
	EmploymentLevel          string  `json:"employment_level"`
	BaseSalary               float64 `json:"base_salary"`
	Currency                 string  `json:"currency"`
	SignOnBonus              float64 `json:"sign_on_bonus"`
	AnnualBonus              float64 `json:"annual_bonus"`
	BenefitsPackage          string  `json:"benefits_package,omitempty"`
	StockOptions             float64 `json:"stock_options"`
	VestingSchedule          string  `json:"vesting_schedule,omitempty"`
	ProbationPeriodDays      int     `json:"probation_period_days"`
	NoticePeriodDays         int     `json:"notice_period_days"`
	PaidTimeOffDays          int     `json:"paid_time_off_days"`
	Status                   string  `json:"status"`
	SentAt                   *string `json:"sent_at,omitempty"`
	RespondedAt             *string `json:"responded_at,omitempty"`
	ResponseDeadline        string  `json:"response_deadline"`
	CandidateResponse       string  `json:"candidate_response,omitempty"`
	NegotiationNotes        string  `json:"negotiation_notes,omitempty"`
	InternalNotes           string  `json:"internal_notes,omitempty"`
	CreatedAt               string  `json:"created_at"`
	UpdatedAt               string  `json:"updated_at"`
}

type OfferListResponse struct {
	Offers []OfferResponse `json:"offers"`
	Meta   PaginationMeta  `json:"meta"`
}

type OfferMetricsResponse struct {
	TotalOffers       int     `json:"total_offers"`
	DraftCount        int     `json:"draft_count"`
	SentCount         int     `json:"sent_count"`
	AcceptedCount     int     `json:"accepted_count"`
	RejectedCount     int     `json:"rejected_count"`
	WithdrawnCount    int     `json:"withdrawn_count"`
	ExpiredCount      int     `json:"expired_count"`
	NegotiationCount  int     `json:"negotiation_count"`
	AcceptanceRate    float64 `json:"acceptance_rate"`
	AverageTimeToAccept int   `json:"average_time_to_accept_hours"`
	TotalBaseSalaryValue float64 `json:"total_base_salary_value"`
}

type PendingOffersResponse struct {
	Offers []OfferResponse `json:"offers"`
}

type OfferComparisonResponse struct {
	Offers []OfferResponse `json:"offers"`
}
