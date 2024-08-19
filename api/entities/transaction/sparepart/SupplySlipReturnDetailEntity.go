package transactionsparepartentities

const TableNameSupplySlipReturnDetail = "trx_supply_slip_return_detail"

type SupplySlipReturnDetail struct {
	SupplyReturnDetailSystemNumber int `gorm:"column:supply_return_detail_system_number;size:30;not null;primaryKey;" json:"supply_return_detail_system_number"`
	SupplyReturnID                 int `gorm:"column:supply_return_system_number;size:30;not null" json:"supply_return_system_number"`
	SupplyReturn                   *SupplySlipReturn
	SupplyReturnLineNumber         int `gorm:"column:supply_return_line_number;null;size:30;" json:"supply_return_line_number"`
	SupplyDetailID                 int `gorm:"column:supply_detail_system_number;size:30;not null" json:"supply_detail_system_number"`
	SupplyDetail                   *SupplySlipDetail
	QuantityReturn                 float32 `gorm:"column:quantity_return;not null" json:"quantity_return"`
	SupplyReturnReasonID           int     `gorm:"column:supply_return_reason_id;null;size:30;" json:"supply_return_reason_id"`
	CostOfGoodsSoldReturn          float32 `gorm:"column:cost_of_good_sold_return;null" json:"cost_of_good_sold_return"`
}

func (*SupplySlipReturnDetail) TableName() string {
	return TableNameSupplySlipReturnDetail
}
