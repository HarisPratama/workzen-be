package entity

import "time"

type ClientEntity struct {
	ID          int64        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TenantID    int64        `json:"tenant_id" gorm:"column:tenant_id;not null;index"`
	CompanyName string       `json:"company_name" gorm:"column:company_name;not null"`
	Address     string       `json:"address" gorm:"column:address"`
	CreatedAt   time.Time    `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Tenant      TenantEntity `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

func (ClientEntity) TableName() string {
	return "clients"
}

type ClientQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
