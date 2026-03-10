package model

import "time"

type Client struct {
	ID          int64     `gorm:"id"`
	CompanyName string    `gorm:"company_name"`
	Address     string    `gorm:"address"`
	TenantID    int64     `gorm:"tenant_id"`
	Tenant      Tenant    `gorm:"foreignkey:TenantID"`
	CreatedAt   time.Time `gorm:"created_at"`
}
