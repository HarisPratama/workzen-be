package entity

type TenantEntity struct {
	ID          int64
	CompanyName string
	Plan        string
	Status      string
	Address     string
}

type RegisterTenantEntity struct {
	CompanyName string
	Address     string
	Plan        string
	Name        string
	Email       string
	Password    string
}
