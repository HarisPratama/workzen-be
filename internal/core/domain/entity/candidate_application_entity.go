package entity

import "time"

type CandidateApplicationEntity struct {
	ID                int64             `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID          int64             `json:"tenant_id" gorm:"column:tenant_id;not null;index"`
	CandidateID       int64             `json:"candidate_id" gorm:"column:candidate_id;not null;index"`
	ManpowerRequestID int64             `json:"manpower_request_id" gorm:"column:manpower_request_id;not null;index"`
	MatchScore        int32             `json:"match_score" gorm:"column:match_score;default:0"`
	MatchVerdict      string            `json:"match_verdict" gorm:"column:match_verdict"`
	Status            string            `json:"status" gorm:"column:status;default:'APPLIED'"`
	AppliedAt         time.Time         `json:"applied_at" gorm:"column:applied_at;autoCreateTime"`
	UpdatedAt         time.Time         `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Tenant            TenantEntity      `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	Candidate         CandidateEntity   `json:"candidate,omitempty" gorm:"foreignKey:CandidateID"`
	ManpowerRequest   ManpowerReqEntity `json:"manpower_request,omitempty" gorm:"foreignKey:ManpowerRequestID"`
}

func (CandidateApplicationEntity) TableName() string {
	return "candidate_applications"
}

type CandidateApplicationQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
