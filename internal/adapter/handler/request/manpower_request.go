package request

import "time"

type ManPowerRequest struct {
	ClientID       int64     `json:"client_id" validate:"required"`
	Position       string    `json:"position" validate:"required"`
	RequiredCount  int       `json:"required_count" validate:"required"`
	SalaryMin      float64   `json:"salary_min" validate:"required"`
	SalaryMax      float64   `json:"salary_max" validate:"required"`
	WorkLocation   string    `json:"work_location" validate:"required"`
	JobDescription string    `json:"job_description" validate:"required"`
	DeadlineDate   time.Time `json:"deadline_date" validate:"required"`
}
