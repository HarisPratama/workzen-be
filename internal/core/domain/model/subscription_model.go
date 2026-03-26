package model

import "time"

type SubscriptionPlan struct {
	ID                  int64     `gorm:"id"`
	Name                string    `gorm:"name"`
	Tier                string    `gorm:"tier"`
	Description         string    `gorm:"description"`
	Price               float64   `gorm:"price"`
	BillingCycle        string    `gorm:"billing_cycle"`
	MaxEmployees        int       `gorm:"max_employees"`
	MaxClients          int       `gorm:"max_clients"`
	MaxManpowerRequests int       `gorm:"max_manpower_requests"`
	Features            string    `gorm:"features;type:jsonb"`
	IsActive            bool      `gorm:"is_active"`
	CreatedAt           time.Time `gorm:"created_at"`
	UpdatedAt           time.Time `gorm:"updated_at"`
}

func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

type TenantSubscription struct {
	ID            int64            `gorm:"id"`
	TenantID      int64            `gorm:"tenant_id"`
	PlanID        int64            `gorm:"plan_id"`
	Status        string           `gorm:"status"`
	StartDate     time.Time        `gorm:"start_date"`
	EndDate       time.Time        `gorm:"end_date"`
	AutoRenew     bool             `gorm:"auto_renew"`
	PaymentMethod string           `gorm:"payment_method"`
	LastPaymentAt *time.Time       `gorm:"last_payment_at"`
	CancelledAt   *time.Time       `gorm:"cancelled_at"`
	CancelReason  string           `gorm:"cancel_reason"`
	CreatedAt     time.Time        `gorm:"created_at"`
	UpdatedAt     time.Time        `gorm:"updated_at"`
	Plan          SubscriptionPlan `gorm:"foreignkey:PlanID"`
	Tenant        Tenant           `gorm:"foreignkey:TenantID"`
}

func (TenantSubscription) TableName() string {
	return "tenant_subscriptions"
}
