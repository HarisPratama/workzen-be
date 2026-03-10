package response

type ManpowerReqResponse struct {
	ID             int64   `json:"id"`
	Position       string  `json:"position"`
	RequiredCount  int     `json:"required_count"`
	SalaryMin      float64 `json:"salary_min"`
	SalaryMax      float64 `json:"salary_max"`
	WorkLocation   string  `json:"work_location"`
	JobDescription string  `json:"job_description"`
	DeadlineDate   string  `json:"deadline_date"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	Client         string  `json:"client"`
}
