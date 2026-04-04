package response

type CVAnalysisResponse struct {
	Summary             string   `json:"summary"`
	Recommendation      string   `json:"recommendation"`
	Skills              []string `json:"skills"`
	ExperienceHighlights []string `json:"experience_highlights"`
	FitScore            int32    `json:"fit_score"`
}

type JobMatchResponse struct {
	Score         int32    `json:"score"`
	MatchedSkills []string `json:"matched_skills"`
	MissingSkills []string `json:"missing_skills"`
	Explanation   string   `json:"explanation"`
	Verdict       string   `json:"verdict"`
}
