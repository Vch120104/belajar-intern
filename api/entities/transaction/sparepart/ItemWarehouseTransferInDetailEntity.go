package transactionsparepartentities

import masteritementities "after-sales/api/entities/master/item"

const TableNameItemWarehouseTransferInDetail = "trx_item_warehouse_transfer_in_detail"

type ItemWarehouseTransferInDetail struct {
	TransferInDetailSystemNumber    int     `gorm:"column:transfer_in_detail_system_number;size:30;not null;primaryKey" json:"transfer_in_detail_system_number"`
	TransferInSystemNumberId        int     `gorm:"column:transfer_in_system_number;size:30;not null" json:"transfer_in_system_number"`
	TransferOutDetailSystemNumberId int     `gorm:"column:transfer_out_detail_system_number;size:30;null" json:"transfer_out_detail_system_number"`
	ItemId                          int     `gorm:"column:item_id;size:30;not null" json:"item_id"`
	LocationId                      int     `gorm:"column:location_id;size:30;null" json:"location_id"`
	QuantityReceived                float64 `gorm:"column:quantity_received;not null" json:"quantity_received"`
	ItemCogs                        float64 `gorm:"column:item_cogs;not null" json:"item_cogs"`
	TransferCost                    float64 `gorm:"column:transfer_cost;null" json:"transfer_cost"`
	ReferencePrice                  float64 `gorm:"column:reference_price;null" json:"reference_price"`
	HppVariance                     float64 `gorm:"column:hpp_variance;null" json:"hpp_variance"`

	Item       masteritementities.Item `gorm:"foreignKey:ItemId;references:item_id" json:"item"`
	TransferIn ItemWarehouseTransferIn `gorm:"foreignKey:TransferInSystemNumberId;references:transfer_in_system_number" json:"trx_item_warehouse_transfer_in"`
}

func (*ItemWarehouseTransferInDetail) TableName() string {
	return TableNameItemWarehouseTransferInDetail
}
