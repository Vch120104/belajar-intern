package transactionjpcbpayloads

import "time"

type SettingTechnicianPayload struct {
	SettingTechnicianSystemNumber int       `json:"setting_technician_system_number"`
	CompanyId                     int       `json:"company_id"`
	EffectiveDate                 time.Time `json:"effective_date"`
}

type SettingTechnicianGetAllResponse struct {
	SettingTechnicianSystemNumber int       `json:"setting_technician_system_number"`
	EffectiveDate                 time.Time `json:"effective_date"`
	SettingId                     int       `json:"setting_id"`
}

type SettingTechnicianGetAllDetailPayload struct {
	SettingTechnicianDetailSystemNumber int  `json:"setting_technician_detail_system_number"`
	SettingTechnicianSystemNumber       int  `json:"setting_technician_system_number"`
	TechnicianNumber                    int  `json:"technician_number"`
	TechnicianEmployeeNumberId          int  `json:"technician_employee_number_id"`
	GroupExpress                        int  `json:"group_express"`
	ShiftGroupId                        int  `json:"shift_group_id"`
	IsBooking                           bool `json:"is_booking"`
}

type SettingTechnicianGetAllDetailResponse struct {
	SettingTechnicianDetailSystemNumber int    `json:"setting_technician_detail_system_number"`
	SettingTechnicianSystemNumber       int    `json:"setting_technician_system_number"`
	TechnicianNumber                    int    `json:"technician_number"`
	UserEmployeeId                      int    `json:"technician_user_id"`
	EmployeeName                        string `json:"technician_employee_name"`
	GroupExpress                        int    `json:"group_express"`
	ShiftGroupId                        int    `json:"shift_group_id"`
	IsBooking                           bool   `json:"is_booking"`
}

type SettingTechnicianEmployeeResponse struct {
	EmployeeName   string `json:"employee_name"`
	UserEmployeeId int    `json:"user_employee_id"`
}

type SettingTechnicianGetByIdResponse struct {
	SettingTechnicianSystemNumber int       `json:"setting_technician_system_number"`
	CompanyId                     int       `json:"company_id"`
	SettingId                     int       `json:"setting_id"`
	EffectiveDate                 time.Time `json:"effective_date"`
}

type SettingTechnicianDetailGetByIdResponse struct {
	SettingTechnicianDetailSystemNumber int  `json:"setting_technician_detail_system_number"`
	ShiftGroupId                        int  `json:"shift_group_id"`
	IsBooking                           bool `json:"is_booking"`
	TechnicianNumber                    int  `json:"technician_number"`
	GroupExpress                        int  `json:"group_express"`
}

type SettingTechnicianSaveRequest struct {
	CompanyId int `json:"company_id" validate:"required"`
}

type SettingTechnicianDetailSaveRequest struct {
	SettingTechnicianSystemNumber int  `json:"setting_technician_system_number"`
	CompanyId                     int  `json:"company_id" validate:"required"`
	TechnicianEmployeeNumberId    int  `json:"technician_employee_number_id"`
	ShiftGroupId                  int  `json:"shift_group_id"`
	TechnicianNumber              int  `json:"technician_number"`
	IsBooking                     bool `json:"is_booking"`
	GroupExpress                  int  `json:"group_express"`
}

type SettingTechnicianDetailUpdateRequest struct {
	ShiftGroupId     int  `json:"shift_group_id" validate:"required"`
	TechnicianNumber int  `json:"technician_number" validate:"required"`
	IsBooking        bool `json:"is_booking"`
	GroupExpress     int  `json:"group_express" validate:"required"`
}
