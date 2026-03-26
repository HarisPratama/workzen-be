package entity

import "time"

type SubscriptionPlanEntity struct {
	ID                  int64
	Name                string
	Tier                string
	Description         string
	Price               float64
	BillingCycle        string
	MaxEmployees        int
	MaxClients          int
	MaxManpowerRequests int
	Features            string
	IsActive            bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type TenantSubscriptionEntity struct {
	ID            int64
	TenantID      int64
	PlanID        int64
	Status        string
	StartDate     time.Time
	EndDate       time.Time
	AutoRenew     bool
	PaymentMethod string
	LastPaymentAt *time.Time
	CancelledAt   *time.Time
	CancelReason  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Plan          SubscriptionPlanEntity
	Tenant        TenantEntity
}

type SubscriptionPlanQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Search    string
	Tier      string
	IsActive  *bool
}

type TenantSubscriptionQueryString struct {
	Limit     int
	Page      int
	OrderBy   string
	OrderType string
	Status    string
}

type SubscriptionPlanRequest struct {
	Name                string  `json:"name" validate:"required"`
	Tier                string  `json:"tier" validate:"required,oneof=FREE PRO ENTERPRISE CUSTOM"`
	Description         string  `json:"description"`
	Price               float64 `json:"price" validate:"gte=0"`
	BillingCycle        string  `json:"billing_cycle" validate:"required,oneof=monthly yearly custom"`
	MaxEmployees        int     `json:"max_employees" validate:"gte=0"`
	MaxClients          int     `json:"max_clients" validate:"gte=0"`
	MaxManpowerRequests int     `json:"max_manpower_requests" validate:"gte=0"`
	Features            string  `json:"features"`
	IsActive            bool    `json:"is_active"`
}

type SubscriptionPlanUpdateRequest struct {
	Name                string  `json:"name"`
	Description         string  `json:"description"`
	Price               float64 `json:"price" validate:"gte=0"`
	BillingCycle        string  `json:"billing_cycle" validate:"omitempty,oneof=monthly yearly custom"`
	MaxEmployees        int     `json:"max_employees" validate:"gte=0"`
	MaxClients          int     `json:"max_clients" validate:"gte=0"`
	MaxManpowerRequests int     `json:"max_manpower_requests" validate:"gte=0"`
	Features            string  `json:"features"`
	IsActive            *bool   `json:"is_active"`
}

type SubscribeTenantRequest struct {
	PlanID        int64  `json:"plan_id" validate:"required"`
	PaymentMethod string `json:"payment_method"`
	AutoRenew     bool   `json:"auto_renew"`
}

type CancelSubscriptionRequest struct {
	CancelReason string `json:"cancel_reason" validate:"required"`
}

type ChangeSubscriptionPlanRequest struct {
	NewPlanID     int64  `json:"new_plan_id" validate:"required"`
	PaymentMethod string `json:"payment_method"`
}
