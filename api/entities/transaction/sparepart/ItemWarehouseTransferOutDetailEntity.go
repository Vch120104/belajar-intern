package transactionsparepartentities

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
)

const TableNameItemWarehouseTransferOutDetail = "trx_item_warehouse_transfer_out_detail"

type ItemWarehouseTransferOutDetail struct {
	TransferOutDetailSystemNumber     int     `gorm:"column:transfer_out_detail_system_number;size:30;not null;primaryKey" json:"transfer_out_detail_system_number"`
	TransferOutSystemNumber           int     `gorm:"column:transfer_out_system_number;size:30;not null" json:"transfer_out_system_number"`
	TransferRequestDetailSystemNumber int     `gorm:"column:transfer_request_detail_system_number;size:30;null" json:"transfer_request_detail_system_number"`
	ItemId                            *int    `gorm:"column:item_id;size:30;null" json:"item_id"`
	QuantityOut                       float64 `gorm:"column:quantity_out;null" json:"quantity_out"`
	LocationIdFrom                    *int    `gorm:"column:location_id_from;size:30;null" json:"location_id_from"`
	LocationIdTo                      *int    `gorm:"column:location_id_to;size:30;null" json:"location_id_to"`
	CostOfGoodsSold                   float64 `gorm:"column:cost_of_goods_sold;null" json:"cost_of_goods_sold"`
	TotalTransferCost                 float64 `gorm:"total_transfer_cost;null" json:"total_transfer_cost"`

	TransferOut           ItemWarehouseTransferOut                  `gorm:"foreignKey:TransferOutSystemNumber;references:transfer_out_system_number" json:"transfer_out"`
	TransferRequestDetail ItemWarehouseTransferRequestDetail        `gorm:"foreignKey:TransferRequestDetailSystemNumber;references:transfer_request_detail_system_number" json:"transfer_request_detail"`
	Item                  masteritementities.Item                   `gorm:"foreignKey:ItemId;references:item_id" json:"item"`
	LocationFrom          masterwarehouseentities.WarehouseLocation `gorm:"foreignKey:LocationIdFrom;references:warehouse_location_id" json:"location_from"`
	LocationTo            masterwarehouseentities.WarehouseLocation `gorm:"foreignKey:LocationIdTo;references:warehouse_location_id" json:"location_to"`
}

func (*ItemWarehouseTransferOutDetail) TableName() string {
	return TableNameItemWarehouseTransferOutDetail
}
