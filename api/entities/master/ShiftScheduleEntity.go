package masterentities

import "time"

const TableNameShiftSchedule = "mtr_shift_schedule"

type ShiftSchedule struct {
	IsActive        bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	ShiftScheduleId int       `gorm:"column:shift_schedule_id;not null;size:30;primaryKey"        json:"shift_schedule_id"`
	CompanyId       int       `gorm:"column:company_id;size:30;not null"        json:"company_id"`
	ShiftCode       string    `gorm:"column:shift_code;size:3;type:char(3);not null"        json:"shift_code"`
	EffectiveDate   time.Time `gorm:"column:effective_date;not null"        json:"effective_date"`
	ShiftGroup      string    `gorm:"column:shift_group;size:10;null"        json:"shift_group"`
	StartTime       float64   `gorm:"column:start_time;not null"        json:"start_time"`
	EndTime         float64   `gorm:"column:end_time;not null"        json:"end_time"`
	RestStartTime   float64   `gorm:"column:rest_start_time;not null"        json:"rest_start_time"`
	RestEndTime     float64   `gorm:"column:rest_end_time;not null"        json:"rest_end_time"`
	Monday          bool      `gorm:"column:monday;null"        json:"monday"`
	Tuesday         bool      `gorm:"column:tuesday;null"        json:"tuesday"`
	Wednesday       bool      `gorm:"column:wednesday;null"        json:"wednesday"`
	Thursday        bool      `gorm:"column:thursday;null"        json:"thursday"`
	Friday          bool      `gorm:"column:friday;null"        json:"friday"`
	Saturday        bool      `gorm:"column:saturday;null"        json:"saturday"`
	Sunday          bool      `gorm:"column:sunday;null"        json:"sunday"`
	Manpower        float64   `gorm:"column:manpower;null"        json:"manpower"`
	ManpowerBooking float64   `gorm:"column:manpower_booking;null"        json:"manpower_booking"`
}

func (*ShiftSchedule) TableName() string {
	return TableNameShiftSchedule
}
