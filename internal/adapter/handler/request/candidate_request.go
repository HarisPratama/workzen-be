package request

import "time"

type CandidateRequest struct {
	FullName  string    `json:"full_name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"required"`
	CitizenID string    `json:"citizen_id"`
	BirthDate time.Time `json:"birth_date" validate:"required"`
	Address   string    `json:"address" validate:"required"`
	Source    string    `json:"source" validate:"required"`
	Status    string    `json:"status"`
}
