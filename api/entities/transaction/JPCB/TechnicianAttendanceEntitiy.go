package transactionjpcbentities

import "time"

const TableNameTechnicianAttendance = "mtr_technician_attendance"

type TechnicianAttendance struct {
	TechnicianAttendanceId int       `gorm:"column:technician_attendance_id;size:30;not null;primaryKey" json:"technician_attendance_id"`
	CompanyId              int       `gorm:"column:company_id;size:30;not null" json:"company_id"` // FK from mtr_company in general-service
	ServiceDate            time.Time `gorm:"column:service_date;uniqueIndex:idx_technician_attendance;not null" json:"service_date"`
	UserId                 int       `gorm:"column:user_id;size:30;uniqueIndex:idx_technician_attendance;not null" json:"user_id"` // FK from mtr_user_details in general-service
	Attendance             bool      `gorm:"column:attendance;default:true;null" json:"attendance"`
}

func (*TechnicianAttendance) TableName() string {
	return TableNameTechnicianAttendance
}
