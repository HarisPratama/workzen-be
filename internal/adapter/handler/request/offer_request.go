package request

type OfferRequest struct {
	CandidateApplicationID string  `json:"candidate_application_id" validate:"required,uuid"`
	JobTitle               string  `json:"job_title" validate:"required"`
	Department             string  `json:"department" validate:"required"`
	OfferType              string  `json:"offer_type" validate:"required,oneof=full_time part_time contract freelance internship"`
	EmploymentLevel        string  `json:"employment_level" validate:"required"`
	BaseSalary             float64 `json:"base_salary" validate:"required,gte=0"`
	Currency               string  `json:"currency" validate:"required,iso4217"`
	SignOnBonus            float64 `json:"sign_on_bonus" validate:"gte=0"`
	AnnualBonus            float64 `json:"annual_bonus" validate:"gte=0"`
	BenefitsPackage        string  `json:"benefits_package"`
	StockOptions           float64 `json:"stock_options" validate:"gte=0"`
	VestingSchedule        string  `json:"vesting_schedule"`
	ProbationPeriodDays    int     `json:"probation_period_days" validate:"gte=0"`
	NoticePeriodDays       int     `json:"notice_period_days" validate:"gte=0"`
	PaidTimeOffDays        int     `json:"paid_time_off_days" validate:"gte=0"`
	ResponseDeadline       string  `json:"response_deadline" validate:"required,datetime=2006-01-02"`
	InternalNotes          string  `json:"internal_notes"`
}

type OfferUpdateRequest struct {
	JobTitle            string  `json:"job_title"`
	Department          string  `json:"department"`
	OfferType           string  `json:"offer_type" validate:"omitempty,oneof=full_time part_time contract freelance internship"`
	EmploymentLevel     string  `json:"employment_level"`
	BaseSalary          float64 `json:"base_salary" validate:"gte=0"`
	Currency            string  `json:"currency" validate:"omitempty,iso4217"`
	SignOnBonus         float64 `json:"sign_on_bonus" validate:"gte=0"`
	AnnualBonus         float64 `json:"annual_bonus" validate:"gte=0"`
	BenefitsPackage     string  `json:"benefits_package"`
	StockOptions        float64 `json:"stock_options" validate:"gte=0"`
	VestingSchedule     string  `json:"vesting_schedule"`
	ProbationPeriodDays int     `json:"probation_period_days" validate:"gte=0"`
	NoticePeriodDays    int     `json:"notice_period_days" validate:"gte=0"`
	PaidTimeOffDays     int     `json:"paid_time_off_days" validate:"gte=0"`
	InternalNotes       string  `json:"internal_notes"`
}

type SendOfferRequest struct {
	EmailSubject string `json:"email_subject"`
	EmailBody    string `json:"email_body"`
}

type WithdrawRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type RejectRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type NegotiateRequest struct {
	ProposedBaseSalary float64 `json:"proposed_base_salary" validate:"gte=0"`
	ProposedBonus      float64 `json:"proposed_bonus" validate:"gte=0"`
	ProposedBenefits   string  `json:"proposed_benefits"`
	Justification      string  `json:"justification" validate:"required"`
}

type OfferFilterRequest struct {
	Status              string `json:"status" validate:"omitempty,oneof=draft sent accepted rejected withdrawn expired"`
	OfferType           string `json:"offer_type" validate:"omitempty,oneof=full_time part_time contract freelance internship"`
	Department          string `json:"department"`
	MinBaseSalary       float64 `json:"min_base_salary" validate:"gte=0"`
	MaxBaseSalary       float64 `json:"max_base_salary" validate:"gte=0"`
	ResponseDeadlineFrom string `json:"response_deadline_from" validate:"omitempty,datetime=2006-01-02"`
	ResponseDeadlineTo   string `json:"response_deadline_to" validate:"omitempty,datetime=2006-01-02"`
}
