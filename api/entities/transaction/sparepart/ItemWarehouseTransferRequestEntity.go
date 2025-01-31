package transactionsparepartentities

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const TableNameItemWarehouseTransferRequest = "trx_item_warehouse_transfer_request"

type ItemWarehouseTransferRequest struct {
	CompanyId                     int        `gorm:"column:company_id;size:30;not null" json:"company_id"`
	TransferRequestSystemNumber   int        `gorm:"column:transfer_request_system_number;size:30;not null;primaryKey" json:"transfer_request_system_number"`
	TransferRequestDocumentNumber string     `gorm:"column:transfer_request_document_number;size:25;not null" json:"transfer_request_document_number"`
	TransferRequestStatusId       int        `gorm:"column:transfer_request_status_id;size:30;not null" json:"transfer_request_status_id"`
	TransferRequestDate           *time.Time `gorm:"column:transfer_request_date;null" json:"transfer_request_date"`
	TransferRequestById           *int       `gorm:"column:transfer_request_by_id;size:30;null" json:"transfer_request_by_id"`
	RequestFromWarehouseId        int        `gorm:"column:request_from_warehouse_id;size:30;not null" json:"request_from_warehouse_id"`
	RequestToWarehouseId          int        `gorm:"column:request_to_warehouse_id;size:30;not null" json:"request_to_warehouse_id"`
	Purpose                       string     `gorm:"column:purpose;size:512;null" json:"purpose"`
	TransferInSystemNumber        *int       `gorm:"column:transfer_in_system_number;size:30;null" json:"transfer_in_system_number"`
	TransferOutSystemNumber       *int       `gorm:"column:transfer_out_system_number;size:30;null" json:"transfer_out_system_number"`
	ApprovalById                  *int       `gorm:"column:approval_by_id;size:30;null" json:"approval_by_id"`
	ApprovalDate                  *time.Time `gorm:"column:approval_date;null" json:"approval_date"`
	ApprovalRemark                string     `gorm:"column:approval_remark;size:256;null" json:"approval_remark"`
	ModifiedById                  int        `gorm:"column:modified_by_id;size:30;null" json:"modified_by_id"`

	TransferRequestStatus masteritementities.ItemTransferStatus   `gorm:"foreignKey:TransferRequestStatusId;references:item_transfer_status_id" json:"transfer_request_status"`
	RequestFromWarehouse  masterwarehouseentities.WarehouseMaster `gorm:"foreignKey:RequestFromWarehouseId;references:warehouse_id" json:"request_from_warehouse"`
	RequestToWarehouse    masterwarehouseentities.WarehouseMaster `gorm:"foreignKey:RequestToWarehouseId;references:warehouse_id" json:"request_to_warehouse"`
	// TransferIn
	// TransferOut
}

func (*ItemWarehouseTransferRequest) TableName() string {
	return TableNameItemWarehouseTransferRequest
}
