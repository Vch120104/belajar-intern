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

type GetItemLocationTransferByIdResponse struct {
	CompanyId                        int        `json:"company_id"`
	TransferRequestSystemNumber      int        `json:"transfer_request_system_number"`
	TransferRequestDocumentNumber    string     `json:"transfer_request_document_number"`
	TransferRequestStatusId          int        `json:"transfer_request_status_id"`
	TransferRequestStatusCode        string     `json:"transfer_request_status_code"`
	TransferRequestStatusDescription string     `json:"transfer_request_status_description"`
	TransferRequestDate              *time.Time `json:"transfer_request_date"`
	TransferRequestById              *int       `json:"transfer_request_by_id"`
	TransferRequestByName            *string    `json:"transfer_request_by_name"`
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
	Purpose                          *string    `json:"purpose"`
	TransferInSystemNumber           *int       `json:"transfer_in_system_number"`
	TransferOutSystemNumber          *int       `json:"transfer_out_system_number"`
	ApprovalById                     *int       `json:"approval_by_id"`
	ApprovalByName                   *string    `json:"approval_by_name"`
	ApprovalDate                     *time.Time `json:"approval_date"`
	ApprovalRemark                   *string    `json:"approval_remark"`
}

type InsertItemLocationTransferRequest struct {
	CompanyId               int     `json:"company_id" validate:"required"`
	TransferRequestById     *int    `json:"transfer_request_by_id"`
	RequestFromWarehouseId  int     `json:"request_from_warehouse_id" validate:"required"`
	RequestToWarehouseId    int     `json:"request_to_warehouse_id" validate:"required"`
	Purpose                 *string `json:"purpose"`
	TransferInSystemNumber  *int    `json:"transfer_in_system_number"`
	TransferOutSystemNumber *int    `json:"transfer_out_system_number"`
}

type UpdateItemLocationTransferRequest struct {
	RequestFromWarehouseId int     `json:"request_from_warehouse_id" validate:"required"`
	RequestToWarehouseId   int     `json:"request_to_warehouse_id" validate:"required"`
	Purpose                *string `json:"purpose"`
}

type AcceptItemLocationTransferRequest struct {
	ApprovalById   *int    `json:"approval_by_id"`
	ApprovalRemark *string `json:"approval_remark"`
}

type RejectItemLocationTransferRequest struct {
	ApprovalById   *int    `json:"approval_by_id"`
	ApprovalRemark *string `json:"approval_remark"`
}

type GetAllItemLocationTransferDetailResponse struct {
	TransferRequestDetailSystemNumber int    `json:"transfer_request_detail_system_number"`
	TransferRequestSystemNumber       int    `json:"transfer_request_system_number"`
	ItemId                            int    `json:"item_id"`
	ItemCode                          string `json:"item_code"`
	ItemName                          string `json:"item_name"`
	UnitOfMeasurementStockId          int    `json:"unit_of_measurement_stock_id"`
	UnitOfMeasurementStockCode        string `json:"unit_of_measurement_stock_code"`
	UnitOfMeasurementStockDescription string `json:"unit_of_measurement_stock_description"`
	RequestQuantity                   int    `json:"request_quantity"`
	LocationIdFrom                    int    `json:"location_id_from"`
	LocationCodeFrom                  string `json:"location_code_from"`
	LocationIdTo                      int    `json:"location_id_to"`
	LocationCodeTo                    string `json:"location_code_to"`
}

type GetItemLocationTransferDetailByIdResponse struct {
	TransferRequestDetailSystemNumber int    `json:"transfer_request_detail_system_number"`
	TransferRequestSystemNumber       int    `json:"transfer_request_system_number"`
	ItemId                            *int   `json:"item_id"`
	ItemCode                          string `json:"item_code"`
	ItemName                          string `json:"item_name"`
	UnitOfMeasurementStockId          int    `json:"unit_of_measurement_stock_id"`
	UnitOfMeasurementStockCode        string `json:"unit_of_measurement_stock_code"`
	UnitOfMeasurementStockDescription string `json:"unit_of_measurement_stock_description"`
	RequestQuantity                   *int   `json:"request_quantity"`
	LocationIdFrom                    *int   `json:"location_id_from"`
	LocationCodeFrom                  string `json:"location_code_from"`
	LocationIdTo                      *int   `json:"location_id_to"`
	LocationCodeTo                    string `json:"location_code_to"`
}

type InsertItemLocationTransferDetailRequest struct {
	TransferRequestSystemNumber int  `json:"transfer_request_system_number" validate:"required"`
	ItemId                      *int `json:"item_id" validate:"required"`
	RequestQuantity             *int `json:"request_quantity" validate:"required"`
	LocationIdFrom              *int `json:"location_id_from" validate:"required"`
	LocationIdTo                *int `json:"location_id_to" validate:"required"`
}

type UpdateItemLocationTransferDetailRequest struct {
	ItemId          *int `json:"item_id" validate:"required"`
	RequestQuantity *int `json:"request_quantity" validate:"required"`
	LocationIdFrom  *int `json:"location_id_from" validate:"required"`
	LocationIdTo    *int `json:"location_id_to" validate:"required"`
}
