package response

type InterviewResponse struct {
	ID              int64  `json:"id"`
	InterviewType   string `json:"interview_type"`
	ScheduledAt     string `json:"scheduled_at"`
	DurationMinutes int    `json:"duration_minutes"`
	Location        string `json:"location"`
	MeetingLink     string `json:"meeting_link"`
	Status          string `json:"status"`
	Feedback        string `json:"feedback"`
	Rating          int    `json:"rating"`
	CompletedAt     string `json:"completed_at"`
	CancelledAt     string `json:"cancelled_at"`
	CancelReason    string `json:"cancel_reason"`
}
