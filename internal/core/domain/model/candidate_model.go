package model

import "time"

type Candidate struct {
	ID        int64     `gorm:"id"`
	TenantID  int64     `gorm:"tenant_id"`
	FullName  string    `gorm:"full_name"`
	Email     string    `gorm:"email"`
	Phone     string    `gorm:"phone"`
	CitizenID string    `gorm:"citizen_id"`
	BirthDate time.Time `gorm:"birth_date"`
	Address   string    `gorm:"address"`
	Source    string    `gorm:"source"`
	Status    string    `gorm:"status"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}
