package transactionjpcbentities

import "time"

const TableNameSettingTechnician = "trx_setting_technician"

type SettingTechnician struct {
	CompanyId                     int        `gorm:"column:company_id;size:30;not null" json:"company_id"`
	EffectiveDate                 *time.Time `gorm:"column:effective_date;not null" json:"effective_date"`
	SettingTechnicianSystemNumber int        `gorm:"column:setting_technician_system_number;size:30;not null;primaryKey;autoincrement" json:"setting_technician_system_number"`
}

func (*SettingTechnician) TableName() string {
	return TableNameSettingTechnician
}
