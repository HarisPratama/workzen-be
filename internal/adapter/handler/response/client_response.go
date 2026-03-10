package response

type ClientResponse struct {
	ID          int64  `json:"id"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	CreatedAt   string `json:"created_at"`
}
