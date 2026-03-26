package response

type SubscriptionPlanResponse struct {
	ID                  int64   `json:"id"`
	Name                string  `json:"name"`
	Tier                string  `json:"tier"`
	Description         string  `json:"description"`
	Price               float64 `json:"price"`
	BillingCycle        string  `json:"billing_cycle"`
	MaxEmployees        int     `json:"max_employees"`
	MaxClients          int     `json:"max_clients"`
	MaxManpowerRequests int     `json:"max_manpower_requests"`
	Features            string  `json:"features"`
	IsActive            bool    `json:"is_active"`
}

type TenantSubscriptionResponse struct {
	ID            int64                    `json:"id"`
	TenantID      int64                    `json:"tenant_id"`
	Status        string                   `json:"status"`
	StartDate     string                   `json:"start_date"`
	EndDate       string                   `json:"end_date"`
	AutoRenew     bool                     `json:"auto_renew"`
	PaymentMethod string                   `json:"payment_method"`
	LastPaymentAt string                   `json:"last_payment_at,omitempty"`
	CancelledAt   string                   `json:"cancelled_at,omitempty"`
	CancelReason  string                   `json:"cancel_reason,omitempty"`
	Plan          SubscriptionPlanResponse `json:"plan"`
}
