package response

type CandidateApplicationResponse struct {
	ID        int64             `json:"id"`
	Status    string            `json:"status"`
	AppliedAt string            `json:"applied_at"`
	Candidate CandidateResponse `json:"candidate"`
}
