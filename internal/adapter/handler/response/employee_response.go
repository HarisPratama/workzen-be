package response

type EmployeeResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CitizenID   string `json:"citizen_id"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}
