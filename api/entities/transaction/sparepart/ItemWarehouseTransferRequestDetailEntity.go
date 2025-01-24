package transactionsparepartentities

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
)

const TableNameItemWarehouseTransferRequestDetail = "trx_item_warehouse_transfer_request_detail"

type ItemWarehouseTransferRequestDetail struct {
	TransferRequestDetailSystemNumber int     `gorm:"column:transfer_request_detail_system_number;size:30;not null;primaryKey" json:"transfer_request_detail_system_number"`
	TransferRequestSystemNumberId     int     `gorm:"column:transfer_request_system_number;size:30;not null" json:"transfer_request_system_number"`
	ItemId                            *int    `gorm:"column:item_id;size:30;null" json:"item_id"`
	RequestQuantity                   float64 `gorm:"column:request_quantity;null" json:"request_quantity"`
	LocationIdFrom                    *int    `gorm:"column:location_id_from;size:30;null" json:"location_id_from"`
	LocationIdTo                      *int    `gorm:"column:location_id_to;size:30;null" json:"location_id_to"`

	TransferRequest ItemWarehouseTransferRequest              `gorm:"foreignKey:TransferRequestSystemNumberId;references:transfer_request_system_number" json:"transfer_request"`
	Item            masteritementities.Item                   `gorm:"foreignKey:ItemId;references:item_id" json:"item"`
	LocationFrom    masterwarehouseentities.WarehouseLocation `gorm:"foreignKey:LocationIdFrom;references:warehouse_location_id" json:"location_from"`
	LocationTo      masterwarehouseentities.WarehouseLocation `gorm:"foreignKey:LocationIdTo;references:warehouse_location_id" json:"location_to"`
}

func (*ItemWarehouseTransferRequestDetail) TableName() string {
	return TableNameItemWarehouseTransferRequestDetail
}
