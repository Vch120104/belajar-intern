package transactionjpcbpayloads

import "time"

type TechnicianAttendanceGetAllResponse struct {
	TechnicianAttendanceId int       `json:"technician_attendance_id"`
	CompanyId              int       `json:"company_id"`
	ServiceDate            time.Time `json:"service_date"`
	UserId                 int       `json:"user_id"`
	Attendance             bool      `json:"attendance"`
}

type TechnicianAttendanceSaveRequest struct {
	CompanyId   int       `json:"company_id"`
	ServiceDate time.Time `json:"service_date"`
	UserId      int       `json:"user_id"`
	Attendance  bool      `json:"attendance"`
}
