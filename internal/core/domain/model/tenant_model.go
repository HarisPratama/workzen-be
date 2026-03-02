package model

import "time"

type Tenant struct {
	ID          int64     `gorm:"id"`
	CompanyName string    `gorm:"company_name"`
	Plan        string    `gorm:"plan"`
	Status      string    `gorm:"status"`
	Address     string    `gorm:"address"`
	CreatedAt   time.Time `gorm:"created_at"`
}
