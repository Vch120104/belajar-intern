package transactionsparepartpayloads

import "time"

type InsertItemWarehouseHeaderTransferOutRequest struct {
	CompanyId                   int       `json:"company_id"`
	TransferOutDate             time.Time `json:"transfer_out_date"`
	TransferRequestSystemNumber int       `json:"transfer_request_system_number"`
	WarehouseId                 *int      `json:"warehouse_id"`
	WarehouseGroupId            *int      `json:"warehouse_group_id"`
}

type InsertItemWarehouseTransferOutDetailRequest struct {
	TransferOutSystemNumber           int       `json:"transfer_out_system_number"`
	TransferOutDate                   time.Time `json:"transfer_out_date"`
	TransferRequestDetailSystemNumber int       `json:"transfer_request_detail_system_number"`
	QuantityOut                       float64   `json:"quantity_out"`
	WarehouseId                       *int      `json:"warehouse_id"`
	WarehouseGroupId                  *int      `json:"warehouse_group_id"`
}

type GetTransferOutByIdResponse struct {
	CompanyId                       int       `json:"company_id"`
	TransferOutSystemNumber         int       `json:"transfer_out_system_number"`
	TransferOutDocumentSystemNumber string    `json:"transfer_out_document_number"`
	TransferStatusId                int       `json:"transfer_out_status_id"`
	TransferOutDate                 time.Time `json:"transfer_out_date"`
	TransferRequestSystemNumber     int       `json:"transfer_request_system_number"`
	WarehouseId                     int       `json:"warehouse_id"`
	WarehouseGroupId                int       `json:"warehouse_group_id"`
	ProfitCenterId                  int       `json:"profit_center_id"`
}

type GetAllTransferOutResponse struct {
	CompanyId                       int       `json:"company_id"`
	TransferOutSystemNumber         int       `json:"transfer_out_system_number"`
	TransferOutDocumentSystemNumber string    `json:"transfer_out_document_number"`
	TransferStatusId                int       `json:"transfer_out_status_id"`
	TransferOutStatusCode           string    `json:"transfer_out_status_code"`
	TransferOutstatusDescription    string    `json:"transfer_out_status_description"`
	TransferOutDate                 time.Time `json:"transfer_out_date"`
	TransferRequestSystemNumber     int       `json:"transfer_request_system_number"`
	WarehouseId                     int       `json:"warehouse_id"`
	WarehouseName                   string    `json:"warehouse_name"`
	WarehouseGroupId                int       `json:"warehouse_group_id"`
	WarehouseGroupName              string    `json:"warehouse_group_name"`
	ProfitCenterId                  int       `json:"profit_center_id"`
}

type GetAllDetailTransferOutResponse struct {
	TransferOutSystemNumber       int     `json:"transfer_out_system_number"`
	TransferOutDetailSystemNumber int     `json:"transfer_out_detail_system_number"`
	ItemId                        int     `json:"item_id"`
	ItemName                      string  `json:"item_name"`
	LocationIdFrom                int     `json:"location_id_from"`
	LocationCodeFrom              string  `json:"location_code_from"`
	LocationIdTo                  int     `json:"location_id_to"`
	LocationCodeTo                string  `json:"location_code_to"`
	QuantityAvailable             float64 `json:"quantity_available"`
	RequestQuantity               float64 `json:"request_quantity"`
	QuantityOut                   float64 `json:"quantity_out"`
	CostOfGoodsSold               float64 `json:"cost_of_goods_sold"`
	UnitOfMeasurement             string  `json:"unit_of_measurement"`
	WarehouseGroupId              int     `json:"warehouse_group_id"`
}

type InsertItemWarehouseTransferOutDetailCopyReceiptRequest struct {
	TransferOutSystemNumber     int `json:"transfer_out_system_number"`
	TransferRequestSystemNumber int `json:"transfer_request_system_number"`
}

type UpdateItemWarehouseTransferOutDetailRequest struct {
	LocationId   int     `json:"location_id"`
	LocationToId int     `json:"location_to_id"`
	QuatityOut   float64 `json:"quantity_out"`
}
