package entity

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceStatus string

const (
	AttendanceStatusPresent   AttendanceStatus = "PRESENT"
	AttendanceStatusAbsent    AttendanceStatus = "ABSENT"
	AttendanceStatusLate      AttendanceStatus = "LATE"
	AttendanceStatusHalfDay   AttendanceStatus = "HALF_DAY"
	AttendanceStatusOnLeave   AttendanceStatus = "ON_LEAVE"
	AttendanceStatusWFH       AttendanceStatus = "WORK_FROM_HOME"
)

type AttendanceType string

const (
	AttendanceTypeRegular  AttendanceType = "REGULAR"
	AttendanceTypeOvertime AttendanceType = "OVERTIME"
	AttendanceTypeHoliday  AttendanceType = "HOLIDAY"
)

// Attendance represents an employee's attendance record
type Attendance struct {
	ID           uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TenantID     uuid.UUID        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	EmployeeID   uuid.UUID        `json:"employee_id" gorm:"type:uuid;not null;index"`
	Date         time.Time        `json:"date" gorm:"type:date;not null;index"`
	CheckIn      *time.Time       `json:"check_in,omitempty"`
	CheckOut     *time.Time       `json:"check_out,omitempty"`
	Status       AttendanceStatus `json:"status" gorm:"type:varchar(20);default:'PRESENT'"`
	Type         AttendanceType   `json:"type" gorm:"type:varchar(20);default:'REGULAR'"`
	WorkHours    float64          `json:"work_hours" gorm:"type:decimal(4,2);default:0"`
	OvertimeHours float64         `json:"overtime_hours" gorm:"type:decimal(4,2);default:0"`
	Location     *string          `json:"location,omitempty"`
	Notes        string           `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time       `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Tenant   *Tenant   `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// CalculateWorkHours calculates total work hours based on check-in and check-out
type AttendanceSummary struct {
	TotalDays       int     `json:"total_days"`
	PresentDays     int     `json:"present_days"`
	AbsentDays      int     `json:"absent_days"`
	LateDays        int     `json:"late_days"`
	HalfDays        int     `json:"half_days"`
	OnLeaveDays     int     `json:"on_leave_days"`
	TotalWorkHours  float64 `json:"total_work_hours"`
	TotalOvertime   float64 `json:"total_overtime"`
}

// AttendanceQueryString represents query parameters for filtering attendance
type AttendanceQueryString struct {
	Page       int       `form:"page" default:"1"`
	Limit      int       `form:"limit" default:"10"`
	EmployeeID uuid.UUID `form:"employee_id"`
	Status     string    `form:"status"`
	StartDate  string    `form:"start_date"`
	EndDate    string    `form:"end_date"`
}

func (a *Attendance) CalculateWorkHours() {
	if a.CheckIn != nil && a.CheckOut != nil {
		duration := a.CheckOut.Sub(*a.CheckIn)
		a.WorkHours = duration.Hours()
	}
}