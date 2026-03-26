package request

type InterviewRequest struct {
	ManpowerRequestID      int64  `json:"manpower_request_id"`
	CandidateApplicationID int64  `json:"candidate_application_id" validate:"required"`
	InterviewerID          int64  `json:"interviewer_id" validate:"required"`
	ScheduledAt            string `json:"scheduled_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	InterviewType          string `json:"interview_type"`
	DurationMinutes        int    `json:"duration_minutes" validate:"required,gte=15"`
	Type                   string `json:"type" validate:"required,oneof=phone video in-person technical hr final"`
	Location               string `json:"location"`
	MeetingLink            string `json:"meeting_link"`
}

type InterviewUpdateRequest struct {
	ScheduledAt     string `json:"scheduled_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	DurationMinutes int    `json:"duration_minutes" validate:"omitempty,gte=15"`
	Type            string `json:"type" validate:"omitempty,oneof=phone video in-person technical hr final"`
	Location        string `json:"location"`
	MeetingLink     string `json:"meeting_link"`
}

type RescheduleRequest struct {
	NewScheduledAt string `json:"new_scheduled_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Reason         string `json:"reason" validate:"required"`
}

type CancelRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type SubmitFeedbackRequest struct {
	Rating          int    `json:"rating" validate:"required,gte=1,lte=5"`
	Strengths       string `json:"strengths"`
	Weaknesses      string `json:"weaknesses"`
	OverallFeedback string `json:"overall_feedback"`
	Recommendation  string `json:"recommendation" validate:"required,oneof=hire no_hire strong_hire"`
}

type CompleteInterviewRequest struct {
	Feedback string `json:"feedback"`
}

type InterviewFilterRequest struct {
	Status    string `json:"status" validate:"omitempty,oneof=scheduled completed cancelled no_show"`
	Type      string `json:"type" validate:"omitempty,oneof=phone video in-person technical hr final"`
	StartDate string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}
