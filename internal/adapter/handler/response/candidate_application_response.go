package response

import "time"

type CandidateApplicationResponse struct {
	ID        int64             `json:"id"`
	Status    string            `json:"status"`
	AppliedAt time.Time         `json:"applied_at"`
	Candidate CandidateResponse `json:"candidate"`
}
