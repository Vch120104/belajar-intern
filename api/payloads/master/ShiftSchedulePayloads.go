package masterpayloads

import "time"

type ShiftScheduleResponse struct {
	IsActive        bool      `json:"is_active"`
	ShiftScheduleId int       `json:"shift_schedule_id"`
	ShiftCode       string    `json:"shift_code"`
	CompanyId       int       `json:"company_id"`
	ShiftGroup      string    `json:"shift_group"`
	EffectiveDate   time.Time `json:"effective_date"`
	StartTime       float64   `json:"start_time"`
	EndTime         float64   `json:"end_time"`
	RestStartTime   float64   `json:"rest_start_time"`
	RestEndTime     float64   `json:"rest_end_time"`
	Monday          bool      `json:"monday"`
	Tuesday         bool      `json:"tuesday"`
	Wednesday       bool      `json:"wednesday"`
	Thursday        bool      `json:"thursday"`
	Friday          bool      `json:"friday"`
	Saturday        bool      `json:"saturday"`
	Sunday          bool      `json:"sunday"`
	Manpower        float64   `json:"manpower"`
	ManpowerBooking float64   `json:"manpower_booking"`
}

type ShiftScheduleUpdate struct {
	ShiftScheduleId int     `json:"shift_schedule_id"`
	CompanyId       int     `json:"company_id"`
	ShiftGroup      string  `json:"shift_group"`
	StartTime       float64 `json:"start_time"`
	EndTime         float64 `json:"end_time"`
	RestStartTime   float64 `json:"rest_start_time"`
	RestEndTime     float64 `json:"rest_end_time"`
	Monday          bool    `json:"monday"`
	Tuesday         bool    `json:"tuesday"`
	Wednesday       bool    `json:"wednesday"`
	Thursday        bool    `json:"thursday"`
	Friday          bool    `json:"friday"`
	Saturday        bool    `json:"saturday"`
	Sunday          bool    `json:"sunday"`
	Manpower        float64 `json:"manpower"`
	ManpowerBooking float64 `json:"manpower_booking"`
}

type ChangeStatusShiftScheduleRequest struct {
	IsActive bool `json:"is_active"`
}

type ShiftScheduleDropDownResponse struct {
	IsActive        bool   `json:"is_active"`
	ShiftScheduleId int    `json:"shift_schedule_id"`
	ShiftCode       string `json:"shift_code"`
}

type ShiftScheduleOutstandingJAResponse struct {
	ShiftScheduleId int       `json:"shift_schedule_id"`
	EffectiveDate   time.Time `json:"effective_date"`
	StartTime       float64   `json:"start_time"`
	EndTime         float64   `json:"end_time"`
	RestStartTime   float64   `json:"rest_start_time"`
	RestEndTime     float64   `json:"rest_end_time"`
	ManpowerBooking float64   `json:"manpower_booking"`
}
