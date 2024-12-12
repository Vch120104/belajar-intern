package transactionsparepartpayloads

import "time"

type GetAllItemLocationTransferResponse struct {
	TransferRequestSystemNumber      int    `json:"transfer_request_system_number"`
	TransferRequestDocumentNumber    string `json:"transfer_request_document_number"`
	TransferRequestStatusId          int    `json:"transfer_request_status_id"`
	TransferRequestStatusCode        string `json:"transfer_request_status_code"`
	TransferRequestStatusDescription string `json:"transfer_request_status_description"`
	TransferRequestDate              string `json:"transfer_request_date"`
	TransferRequestById              int    `json:"transfer_request_by_id"`
	TransferRequestByName            string `json:"transfer_request_by_name"`
	RequestFromWarehouseId           int    `json:"request_from_warehouse_id"`
	RequestFromWarehouseCode         string `json:"request_from_warehouse_code"`
	RequestFromWarehouseName         string `json:"request_from_warehouse_name"`
	RequestFromWarehouseGroupId      int    `json:"request_from_warehouse_group_id"`
	RequestFromWarehouseGroupCode    string `json:"request_from_warehouse_group_code"`
	RequestFromWarehouseGroupName    string `json:"request_from_warehouse_group_name"`
	RequestToWarehouseId             int    `json:"request_to_warehouse_id"`
	RequestToWarehouseCode           string `json:"request_to_warehouse_code"`
	RequestToWarehouseName           string `json:"request_to_warehouse_name"`
	RequestToWarehouseGroupId        int    `json:"request_to_warehouse_group_id"`
	RequestToWarehouseGroupCode      string `json:"request_to_warehouse_group_code"`
	RequestToWarehouseGroupName      string `json:"request_to_warehouse_group_name"`
}

type InsertItemLocationTransferRequest struct {
	CompanyId               int        `json:"company_id"`
	TransferRequestStatusId int        `json:"transfer_request_status_id"`
	TransferRequestDate     *time.Time `json:"transfer_request_date"`
	TransferRequestById     *int       `json:"transfer_request_by_id"`
	RequestFromWarehouseId  int        `json:"request_from_warehouse_id"`
	RequestToWarehouseId    int        `json:"request_to_warehouse_id"`
	Purpose                 *string    `json:"purpose"`
	TransferInSystemNumber  *int       `json:"transfer_in_system_number"`
	TransferOutSystemNumber *int       `json:"transfer_out_system_number"`
}
