package transactionsparepartentities

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const TableNameItemWarehouseTransferIn = "trx_item_warehouse_transfer_in"

type ItemWarehouseTransferIn struct {
	CompanyId                 int       `gorm:"column:company_id;size:30;not null" json:"company_id"`
	TransferInSystemNumber    int       `gorm:"column:transfer_in_system_number;size:30;not null;primaryKey" json:"transfer_in_system_number"`
	TransferInDocumentNumber  string    `gorm:"column:transfer_in_document_number;size:25;not null" json:"transfer_in_document_number"`
	TransferInStatusId        int       `gorm:"column:transfer_in_status;size:30;not null" json:"transfer_in_status"`
	TransferInDate            time.Time `gorm:"column:transfer_in_date;not null" json:"transfer_in_date"`
	TransferOutSystemNumberId int       `gorm:"column:transfer_out_system_number;size:30;not null" json:"transfer_out_system_number"`
	EventId                   int       `gorm:"event_id;size:30;null" json:"event_id"`
	JournalSystemNumber       int       `gorm:"journal_system_number;size:30;null" json:"journal_system_number"`
	WarehouseId               *int      `gorm:"column:warehouse_id;size:30;null" json:"warehouse_id"`
	ProfitCenterId            int       `gorm:"column:profit_center_id;size:30;null" json:"profir_center_id"`

	TransferOut ItemWarehouseTransferOut                `gorm:"foreignKey:TransferOutSystemNumberId;references:transfer_out_system_number" json:"trx_item_warehouse_transfer_out"`
	Warehouse   masterwarehouseentities.WarehouseMaster `gorm:"foreignKey:WarehouseId;references:warehouse_id" json:"warehouse"`
}

func (*ItemWarehouseTransferIn) TableName() string {
	return TableNameItemWarehouseTransferIn
}
