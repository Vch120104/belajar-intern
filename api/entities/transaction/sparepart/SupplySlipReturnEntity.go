package transactionsparepartentities

import "time"

const TableNameSupplySlipReturn = "trx_supply_slip_return"

type SupplySlipReturn struct {
	IsActive                   bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SupplyReturnSystemNumber   int    `gorm:"column:supply_return_system_number;size:30;not null;primaryKey;" json:"supply_return_system_number"`
	SupplyReturnDocumentNumber string `gorm:"column:supply_return_document_number;size:50;" json:"supply_return_document_number"`
	SupplyID                   int    `gorm:"column:supply_system_number;size:30;not null" json:"supply_system_number"`
	Supply                     *SupplySlip
	SupplyReturnStatusId       int        `gorm:"column:supply_return_status_id;size:30;not null" json:"supply_return_status_id"`
	SupplyReturnDate           *time.Time `gorm:"column:supply_return_date;not null" json:"supply_return_date"`
	Remark                     string     `gorm:"column:remark;size:50;null" json:"remark"`
}

func (*SupplySlipReturn) TableName() string {
	return TableNameSupplySlipReturn
}
