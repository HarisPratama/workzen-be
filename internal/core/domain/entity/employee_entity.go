package entity

import "time"

type EmployeeEntity struct {
	ID          int64        `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string       `json:"name" gorm:"column:name;not null"`
	PhoneNumber string       `json:"phone_number" gorm:"column:phone_number"`
	CitizenID   string       `json:"citizen_id" gorm:"column:citizen_id"`
	UserID      *int64       `json:"user_id" gorm:"column:user_id"`
	TenantID    int64        `json:"tenant_id" gorm:"column:tenant_id;not null;index"`
	Status      string       `json:"status" gorm:"column:status;default:'ACTIVE'"`
	CreatedAt   time.Time    `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	User        UserEntity   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Tenant      TenantEntity `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

func (EmployeeEntity) TableName() string {
	return "employees"
}

type EmployeeQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Status    string
}
