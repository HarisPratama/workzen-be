package entity

import "time"

type UserEntity struct {
	ID        int64      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"column:name"`
	Email     string     `json:"email" gorm:"column:email"`
	Password  string     `json:"-" gorm:"column:password"`
	Role      string     `json:"role" gorm:"column:role"`
	Status    string     `json:"status" gorm:"column:status"`
	TenantID  *int64     `json:"tenant_id" gorm:"column:tenant_id"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (UserEntity) TableName() string {
	return "users"
}
