package entity

import "time"

type TenantEntity struct {
	ID          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	CompanyName string    `json:"company_name" gorm:"column:company_name"`
	Plan        string    `json:"plan" gorm:"column:plan"`
	Status      string    `json:"status" gorm:"column:status"`
	Address     string    `json:"address" gorm:"column:address"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (TenantEntity) TableName() string {
	return "tenants"
}

type RegisterTenantEntity struct {
	CompanyName string
	Address     string
	Plan        string
	Name        string
	Email       string
	Password    string
}
