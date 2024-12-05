package masterentities

const SupplySlipStatusTableName = "mtr_supply_slip_status"

type SupplySlipStatus struct {
	SupplySlipStatusId   int    `gorm:"column:supply_slip_status_id;size:30;primaryKey" json:"supply_slip_status_id"`
	SupplySlipStatusCode string `gorm:"column:supply_slip_status_code;size:20;not null" json:"supply_slip_status_code"`
	SupplySlipStatusName string `gorm:"column:supply_slip_status_name;size:256;not null" json:"supply_slip_status_name"`
}

func (*SupplySlipStatus) TableName() string {
	return SupplySlipStatusTableName
}
