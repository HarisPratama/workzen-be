package response

type InterviewResponse struct {
	ID                       string  `json:"id"`
	TenantID                 string  `json:"tenant_id"`
	CandidateApplicationID   string  `json:"candidate_application_id"`
	InterviewerID           string  `json:"interviewer_id"`
	ScheduledAt             string  `json:"scheduled_at"`
	DurationMinutes         int     `json:"duration_minutes"`
	Type                    string  `json:"type"`
	Status                  string  `json:"status"`
	Location                string  `json:"location,omitempty"`
	MeetingLink             string  `json:"meeting_link,omitempty"`
	Feedback                string  `json:"feedback,omitempty"`
	Rating                  *int    `json:"rating,omitempty"`
	Recommendation          string  `json:"recommendation,omitempty"`
	CreatedAt               string  `json:"created_at"`
	UpdatedAt               string  `json:"updated_at"`
}

type InterviewListResponse struct {
	Interviews []InterviewResponse `json:"interviews"`
	Meta       PaginationMeta      `json:"meta"`
}

type InterviewMetricsResponse struct {
	TotalInterviews      int     `json:"total_interviews"`
	ScheduledCount       int     `json:"scheduled_count"`
	CompletedCount       int     `json:"completed_count"`
	CancelledCount       int     `json:"cancelled_count"`
	NoShowCount          int     `json:"no_show_count"`
	AverageRating        float64 `json:"average_rating"`
	RecommendationRate   float64 `json:"recommendation_rate"`
	AverageDurationMins  int     `json:"average_duration_mins"`
}

type UpcomingInterviewsResponse struct {
	Interviews []InterviewResponse `json:"interviews"`
}
