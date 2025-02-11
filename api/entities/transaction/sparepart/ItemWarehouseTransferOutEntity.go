package transactionsparepartentities

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const TableNameItemWarehouseTransferOut = "trx_item_warehouse_transfer_out"

type ItemWarehouseTransferOut struct {
	CompanyId                    int       `gorm:"column:company_id;size:30;not null" json:"company_id"`
	TransferOutSystemNumber      int       `gorm:"column:transfer_out_system_number;size:30;not null;primaryKey" json:"transfer_out_system_number"`
	TransferOutDocumentNumber    string    `gorm:"column:transfer_out_document_number;size:25;null" json:"transfer_out_document_number"`
	TransferRequestSystemNumbers int       `gorm:"column:transfer_request_system_number;size:30;not null" json:"transfer_request_system_number"`
	TransferOutDate              time.Time `gorm:"column:transfer_out_date; null" json:"transfer_out_date"`
	TransferOutStatusId          int       `gorm:"column:transfer_out_status_id;size:30;not null" json:"transfer_out_status_id"`
	WarehouseId                  *int      `gorm:"column:warehouse_id;size:30;null" json:"warehouse_id"`
	ProfitCenterId               int       `gorm:"column:profit_center_id;size:30;null" json:"profir_center_id"`

	TransferRequest ItemWarehouseTransferRequest            `gorm:"foreignKey:TransferRequestSystemNumbers;references:transfer_request_system_number" json:"transfer_request"`
	Warehouse       masterwarehouseentities.WarehouseMaster `gorm:"foreignKey:WarehouseId;references:warehouse_id" json:"warehouse"`
}

func (*ItemWarehouseTransferOut) TableName() string {
	return TableNameItemWarehouseTransferOut
}
