package request

type RegisterTenantRequest struct {
	CompanyName    string `json:"company_name" validate:"required"`
	CompanyAddress string `json:"company_address"`
	Plan           string `json:"plan"`
	AdminName      string `json:"admin_name" validate:"required"`
	AdminEmail     string `json:"admin_email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
}
