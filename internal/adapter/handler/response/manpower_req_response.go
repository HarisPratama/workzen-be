package response

type ManpowerReqResponse struct {
	ID             int64          `json:"id"`
	Position       string         `json:"position"`
	RequiredCount  int            `json:"required_count"`
	Hired          int            `json:"hired"`
	SalaryMin      float64        `json:"salary_min"`
	SalaryMax      float64        `json:"salary_max"`
	WorkLocation   string         `json:"work_location"`
	JobDescription string         `json:"job_description"`
	DeadlineDate   string         `json:"deadline_date"`
	PublicToken    string         `json:"public_token,omitempty"`
	Status         string         `json:"status"`
	CreatedAt      string         `json:"created_at"`
	Client         ClientResponse `json:"client"`
}

type JobPostingResponse struct {
	Position       string  `json:"position"`
	CompanyName    string  `json:"company_name"`
	WorkLocation   string  `json:"work_location"`
	SalaryMin      float64 `json:"salary_min"`
	SalaryMax      float64 `json:"salary_max"`
	JobDescription string  `json:"job_description"`
	DeadlineDate   string  `json:"deadline_date"`
	Status         string  `json:"status"`
}

type JobApplyResponse struct {
	Message       string   `json:"message"`
	Score         int32    `json:"score"`
	Verdict       string   `json:"verdict"`
	MatchedSkills []string `json:"matched_skills"`
	MissingSkills []string `json:"missing_skills"`
	Explanation   string   `json:"explanation"`
}
