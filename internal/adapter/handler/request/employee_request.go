package request

type EmployeeRequest struct {
	Name        string `json:"name" validate:"required"`
	CitizenID   string `json:"citizen_id" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
