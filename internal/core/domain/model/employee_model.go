package model

import "time"

type Employee struct {
	ID          int64     `gorm:"id"`
	Name        string    `gorm:"name"`
	PhoneNumber string    `gorm:"phone_number"`
	CitizenID   string    `gorm:"citizen_id"`
	Status      string    `gorm:"status"`
	UserID      *int64    `gorm:"user_id"`
	TenantID    int64     `gorm:"tenant_id"`
	CreatedAt   time.Time `gorm:"created_at"`
	User        User      `gorm:"foreignkey:UserID"`
	Tenant      Tenant    `gorm:"foreignkey:TenantID"`
}
