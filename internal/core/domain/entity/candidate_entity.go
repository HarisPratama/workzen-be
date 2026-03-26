package entity

import "time"

type CandidateEntity struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID  int64     `json:"tenant_id" gorm:"column:tenant_id;not null;index"`
	FullName  string    `json:"full_name" gorm:"column:full_name;not null"`
	Email     string    `json:"email" gorm:"column:email;not null"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	CitizenID string    `json:"citizen_id" gorm:"column:citizen_id"`
	BirthDate time.Time `json:"birth_date" gorm:"column:birth_date"`
	Address   string    `json:"address" gorm:"column:address"`
	Source    string    `json:"source" gorm:"column:source"`
	Status    string    `json:"status" gorm:"column:status;default:'ACTIVE'"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (CandidateEntity) TableName() string {
	return "candidates"
}

type CandidateQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
