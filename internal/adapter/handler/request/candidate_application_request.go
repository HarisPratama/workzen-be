package request

type CandidateApplicationRequest struct {
	CandidateID       int64 `json:"candidate_id" validate:"required"`
	ManpowerRequestID int64 `json:"manpower_request_id" validate:"required"`
}

type CandidateApplicationUpdateRequest struct {
	Status string `json:"status" validate:"required"`
}
