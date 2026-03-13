package response

type CandidateResponse struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CitizenID string `json:"citizen_id"`
	BirthDate string `json:"birth_date"`
	Address   string `json:"address"`
	Source    string `json:"source"`
	Status    string `json:"status"`
}
