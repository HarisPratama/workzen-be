package request

type RegisterTenantRequest struct {
	CompanyName    string `json:"company_name" binding:"required"`
	CompanyAddress string `json:"company_address"`
	Plan           string `json:"plan"`
	AdminName      string `json:"admin_name" binding:"required"`
	AdminEmail     string `json:"admin_email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
}
