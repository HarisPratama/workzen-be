package request

type AnalyzeCVRequest struct {
	CvText string `json:"cv_text" validate:"required"`
}

type MatchJobRequest struct {
	CvText string `json:"cv_text" validate:"required"`
	JdText string `json:"jd_text" validate:"required"`
}
