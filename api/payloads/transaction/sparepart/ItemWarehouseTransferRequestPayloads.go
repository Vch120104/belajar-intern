package transactionsparepartpayloads

import "time"

type InsertItemWarehouseTransferRequest struct {
	CompanyId              int        `json:"company_id" parent_entity:"trx_item_warehouse_transfer_request"`
	TransferRequestDate    *time.Time `json:"transfer_request_date" parent_entity:"trx_item_warehouse_transfer_request"`
	TransferRequestById    *int       `json:"transfer_request_by_id" parent_entity:"trx_item_warehouse_transfer_request"`
	RequestFromWarehouseId int        `json:"request_from_warehouse_id" parent_entity:"trx_item_warehouse_transfer_request"`
	RequestToWarehouseId   int        `json:"request_to_warehouse_id" parent_entity:"trx_item_warehouse_transfer_request"`
	Purpose                string     `json:"purpose" parent_entity:"trx_item_warehouse_transfer_request"`
}

type InsertItemWarehouseTransferDetailRequest struct {
	TransferRequestSystemNumberId int     `json:"transfer_request_system_number" parent_entity:"trx_item_warehouse_transfer_request_detail"`
	ItemId                        *int    `json:"item_id" parent_entity:"trx_item_warehouse_transfer_request_detail"`
	RequestQuantity               float64 `json:"request_quantity" parent_entity:"trx_item_warehouse_transfer_request_detail"`
	ModifiedById                  int     `json:"modified_by_id"`
}

type UpdateItemWarehouseTransferRequest struct {
	ModifiedById           int    `json:"modified_by_id"`
	TransferRequestById    *int   `json:"transfer_request_by_id" parent_entity:"trx_item_warehouse_transfer_request"`
	RequestFromWarehouseId *int   `json:"request_from_warehouse_id" parent_entity:"trx_item_warehouse_transfer_request"`
	RequestToWarehouseId   *int   `json:"request_to_warehouse_id" parent_entity:"trx_item_warehouse_transfer_request"`
	Purpose                string `json:"purpose" parent_entity:"trx_item_warehouse_transfer_request"`
}

type GetAllDetailItemWarehouseTransferRequestResponse struct {
	TransferRequestDetailSystemNumber int     `json:"transfer_request_detail_system_number"`
	TransferRequestSystemNumber       int     `json:"transfer_request_system_number"`
	ItemId                            int     `json:"item_id"`
	ItemCode                          string  `json:"item_code"`
	ItemName                          string  `json:"item_name"`
	UnitOfMeasurement                 string  `json:"unit_of_measurement"`
	RequestQuantity                   float64 `json:"request_quantity"`
	LocationIdFrom                    int     `json:"location_id_from"`
	LocationIdTo                      int     `json:"location_id_to"`
}

type GetByIdItemWarehouseTransferRequestResponse struct {
	TransferRequestSystemNumber      int        `json:"transfer_request_system_number"`
	TransferRequestDocumentNumber    string     `json:"transfer_request_document_number"`
	TransferRequestStatusId          int        `json:"transfer_request_status_id"`
	TransferRequestStatusCode        string     `json:"transfer_request_status_code"`
	TransferRequestStatusDescription string     `json:"transfer_request_status_description"`
	TransferRequestDate              time.Time  `json:"transfer_request_date"`
	RequestFromWarehouseId           int        `json:"request_from_warehouse_id"`
	RequestFromWarehouseCode         string     `json:"request_from_warehouse_code"`
	RequestFromWarehouseName         string     `json:"request_from_warehouse_name"`
	RequestFromWarehouseGroupId      int        `json:"request_from_warehouse_group_id"`
	RequestFromWarehouseGroupCode    string     `json:"request_from_warehouse_group_code"`
	RequestFromWarehouseGroupName    string     `json:"request_from_warehouse_group_name"`
	RequestToWarehouseId             int        `json:"request_to_warehouse_id"`
	RequestToWarehouseCode           string     `json:"request_to_warehouse_code"`
	RequestToWarehouseName           string     `json:"request_to_warehouse_name"`
	RequestToWarehouseGroupId        int        `json:"request_to_warehouse_group_id"`
	RequestToWarehouseGroupCode      string     `json:"request_to_warehouse_group_code"`
	RequestToWarehouseGroupName      string     `json:"request_to_warehouse_group_name"`
	Purpose                          string     `json:"purpose"`
	ApprovalById                     *int       `json:"approval_by_id"`
	ApprovalDate                     *time.Time `json:"approval_date"`
	ApprovalRemark                   string     `json:"approval_remark"`
	CreatedById                      int        `json:"created_by_id"`
	CreatedByName                    string     `json:"created_by_name"`
	ModifiedById                     int        `json:"modified_by_id"`
	ModifiedByName                   string     `json:"modified_by_name"`
}

// type GetAllItemWarehouseTransferRequestRequest struct {
// 	TransferRequestDocumentNumber   string    `json:"transfer_request_document_number"`
// 	TransferRequestDateFrom         time.Time `json:"transfer_request_date_from"`
// 	TransferRequestDateTo           time.Time `json:"transfer_request_date_to"`
// 	TransferRequestWarehouseGroupId *int      `json:"transfer_request_warehouse_group_id"`
// 	TransferRequestStatusId         *int      `json:"transfer_request_status_id"`
// }

type GetAllItemWarehouseTransferRequestResponse struct {
	TransferRequestSystemNumber      int       `json:"transfer_request_system_number"`
	TransferRequestDocumentNumber    string    `json:"transfer_request_document_number"`
	TransferRequestDate              time.Time `json:"transfer_request_date"`
	TransferRequestStatusId          int       `json:"transfer_request_status_id"`
	TransferRequestStatusCode        string    `json:"transfer_request_status_code"`
	TransferRequestStatusDescription string    `json:"transfer_request_status_description"`
	TransferRequestById              int       `json:"transfer_request_by_id"`
	TransferRequestByName            string    `json:"transfer_request_by_name"`
	RequestFromWarehouseId           int       `json:"request_from_warehouse_id"`
	RequestFromWarehouseName         string    `json:"request_from_warehouse_name"`
	RequestFromWarehouseGroupId      int       `json:"request_from_warehouse_group_id"`
	RequestFromWarehouseGroupName    string    `json:"request_from_warehouse_group_name"`
	RequestToWarehouseId             int       `json:"request_to_warehouse_id"`
	RequestToWarehouseName           string    `json:"request_to_warehouse_name"`
	RequestToWarehouseGroupId        int       `json:"request_to_warehouse_group_id"`
	RequestToWarehouseGroupName      string    `json:"request_to_warehouse_group_name"`
}

type UpdateItemWarehouseTransferRequestDetailRequest struct {
	RequestQuantity float64 `json:"request_quantity"`
	ModifiedById    int     `json:"modified_by_id"`
}

type GetByIdItemWarehouseTransferRequestDetailResponse struct {
	ItemId   int     `json:"item_id"`
	ItemCode string  `json:"item_code"`
	StockUom string  `json:"stock_uom"`
	Quantity float64 `json:"quantity"`
}

type UploadPreviewItemWarehouseTransferRequestPayloads struct {
	ItemCode          string  `json:"item_code"`
	ItemName          string  `json:"item_name"`
	RequestQuantity   float64 `json:"request_quantity"`
	UnitOfMeasurement string  `json:"unit_of_measurement"`
}

type UploadProcessItemWarehouseTransferRequestPayloads struct {
	TransferRequestSystemNumber int                                                 `json:"transfer_request_system_number"`
	ModifiedById                int                                                 `json:"modified_by_id"`
	TransferRequestDetails      []UploadPreviewItemWarehouseTransferRequestPayloads `json:"transfer_request_details"`
}

type SubmitItemWarehouseTransferRequest struct {
	ModifiedById int `json:"modified_by_id"`
}

type DeleteDetailItemWarehouseTransferRequest struct {
	ModifiedById int `json:"modified_by_id"`
}

type AcceptWarehouseTransferRequestRequest struct {
	ApprovalById   int       `json:"approval_by_id"`
	ApprovalDate   time.Time `json:"approval_date"`
	ApprovalRemark string    `json:"approval_remark"`
}

type RejectWarehouseTransferRequestRequest struct {
	ApprovalById   int       `json:"approval_by_id"`
	ApprovalDate   time.Time `json:"approval_date"`
	ApprovalRemark string    `json:"approval_remark"`
}

type GetAllItemWarehouseLookUp struct {
	TransferRequestSystemNumber   int       `json:"transfer_request_system_number"`
	TransferRequestDocumentNumber string    `json:"transfer_request_document_number"`
	TransferRequestDate           time.Time `json:"transfer_request_date"`
	TransferRequestById           int       `json:"transfer_request_by_id"`
	TransferRequestByName         string    `json:"transfer_request_by_name"`
	RequestFromWarehouseId        int       `json:"request_from_warehouse_id"`
	RequestFromWarehouseName      string    `json:"request_from_warehouse_name"`
	RequestFromWarehouseCode      string    `json:"request_from_warehouse_code"`
	RequestFromWarehouseGroupId   int       `json:"request_from_warehouse_group_id"`
	RequestFromWarehouseGroupName string    `json:"request_from_warehouse_group_name"`
	RequestFromWarehouseGroupCode string    `json:"request_from_warehouse_group_code"`
}

type GetAllItemWarehouseDetailLookUp struct {
	ItemId            int     `json:"transfer_request_system_number"`
	ItemCode          string  `json:"transfer_request_document_number"`
	ItemDescription   string  `json:"item_description"`
	RequestQuantity   float64 `json:"request_quantity"`
	UnitOfMeasurement string  `json:"unit_of_measurement"`
}
