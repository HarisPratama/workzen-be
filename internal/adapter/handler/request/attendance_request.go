package request

type AttendanceRequest struct {
	EmployeeID string `json:"employee_id" validate:"required,uuid"`
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
	Type       string `json:"type" validate:"required,oneof=regular overtime weekend holiday"`
	Status     string `json:"status" validate:"required,oneof=present absent leave holiday"`
	CheckIn    string `json:"check_in" validate:"omitempty,datetime=15:04:05"`
	CheckOut   string `json:"check_out" validate:"omitempty,datetime=15:04:05"`
	Notes      string `json:"notes"`
}

type AttendanceUpdateRequest struct {
	Status   string `json:"status" validate:"omitempty,oneof=present absent leave holiday"`
	CheckIn  string `json:"check_in" validate:"omitempty,datetime=15:04:05"`
	CheckOut string `json:"check_out" validate:"omitempty,datetime=15:04:05"`
	Notes    string `json:"notes"`
}

type CheckInRequest struct {
	CheckInTime string `json:"check_in_time" validate:"omitempty,datetime=15:04:05"`
	Location    string `json:"location"`
}

type CheckOutRequest struct {
	CheckOutTime string `json:"check_out_time" validate:"omitempty,datetime=15:04:05"`
}

type UpdateAttendanceStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=present absent leave holiday"`
	Notes  string `json:"notes"`
}

type AttendancePeriodRequest struct {
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02"`
}
