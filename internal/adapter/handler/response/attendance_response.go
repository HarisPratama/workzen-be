package response

type AttendanceResponse struct {
	ID            string  `json:"id"`
	TenantID      string  `json:"tenant_id"`
	EmployeeID    string  `json:"employee_id"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Status        string  `json:"status"`
	CheckIn       *string `json:"check_in,omitempty"`
	CheckOut      *string `json:"check_out,omitempty"`
	WorkHours     float64 `json:"work_hours"`
	OvertimeHours float64 `json:"overtime_hours"`
	Location      string  `json:"location,omitempty"`
	Notes         string  `json:"notes,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type AttendanceListResponse struct {
	Attendances []AttendanceResponse `json:"attendances"`
	Meta        PaginationResponse   `json:"meta"`
}

type AttendanceSummaryResponse struct {
	TotalDays          int     `json:"total_days"`
	PresentDays        int     `json:"present_days"`
	AbsentDays         int     `json:"absent_days"`
	LeaveDays          int     `json:"leave_days"`
	HolidayDays        int     `json:"holiday_days"`
	TotalWorkHours     float64 `json:"total_work_hours"`
	TotalOvertimeHours float64 `json:"total_overtime_hours"`
	AverageWorkHours   float64 `json:"average_work_hours"`
}

type TodayAttendanceResponse struct {
	HasCheckedIn  bool    `json:"has_checked_in"`
	HasCheckedOut bool    `json:"has_checked_out"`
	CheckInTime   *string `json:"check_in_time,omitempty"`
	CheckOutTime  *string `json:"check_out_time,omitempty"`
	WorkHours     float64 `json:"work_hours"`
	Status        string  `json:"status"`
}
