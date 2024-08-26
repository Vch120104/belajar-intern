package transactionjpcbentities

import masterentities "after-sales/api/entities/master"

const TableNameSettingTechnicianDetail = "trx_setting_technician_detail"

type SettingTechnicianDetail struct {
	SettingTechnicianDetailSystemNumber int                           `gorm:"column:setting_technician_detail_system_number;size:30;not null;primaryKey" json:"setting_technician_detail_system_number"`
	SettingTechnicianSystemNumbers      int                           `gorm:"column:setting_technician_system_number;not null;size:30" json:"setting_technician_system_number"`
	SettingTechnician                   *SettingTechnician            `gorm:"foreignKey:SettingTechnicianSystemNumbers;references:SettingTechnicianSystemNumber"`
	TechnicianEmployeeNumberId          int                           `gorm:"column:technician_employee_number_id;size:30;null" json:"technician_employee_number_id"` // FK from mtr_user_details on general-service
	ShiftGroupId                        int                           `gorm:"column:shift_group_id;size:30;null" json:"shift_group_id"`
	ShiftSchedule                       *masterentities.ShiftSchedule `gorm:"foreignKey:ShiftGroupId;references:ShiftScheduleId"`
	TechnicianNumber                    int                           `gorm:"column:technician_number;size:30;default:0;null" json:"technician_number"`
	IsBooking                           bool                          `gorm:"column:is_booking;default:false;null" json:"is_booking"`
	GroupExpress                        int                           `gorm:"column:group_express;default:0;null" json:"group_express"`
}

func (*SettingTechnicianDetail) TableName() string {
	return TableNameSettingTechnicianDetail
}
