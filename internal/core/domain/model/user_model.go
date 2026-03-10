package model

import "time"

type User struct {
	ID        int64      `gorm:"id"`
	TenantID  *int64     `gorm:"tenant_id"`
	Name      string     `gorm:"name"`
	Email     string     `gorm:"email"`
	Role      string     `gorm:"role"`
	Status    string     `gorm:"status"`
	Password  string     `gorm:"password"`
	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt *time.Time `gorm:"updated_at"`
}
