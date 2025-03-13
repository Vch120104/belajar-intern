package transactionsparepartpayloads

import "time"

type StockOpnameInsertRequest struct {
	CompanyId                 int       `json:"company_id"`
	StockOpnameSystemNumber   int       `json:"stock_opname_system_number"`
	StockOpnameDocumentNumber string    `json:"stock_opname_document_number"`
	WarehouseGroup            int       `json:"warehouse_group" validate:"required"`
	WarehouseCode             int       `json:"warehouse_code" validate:"required"`
	FromLocation              int       `json:"from_location" validate:"required"`
	ToLocation                int       `json:"to_location" validate:"required"`
	ItemGroup                 int       `json:"item_group" validate:"required"`
	ShowDetail                bool      `json:"show_details"`
	PersonInCharge            int       `json:"person_in_charge" validate:"required"`
	Remark                    string    `json:"remark"`
	ExecutionDateFrom         time.Time `json:"execution_date_from" validate:"required"`
	ExecutionDateTo           time.Time `json:"execution_date_to" validate:"required"`
	IncludeZeroOnhand         bool      `json:"include_zero_onhand"`
	ApprovalRequestedById     int       `json:"approval_requested_by_id"`
	ApprovalRequestedDate     time.Time `json:"approval_requested_date"`
	ApprovalById              int       `json:"approval_by_id"`
	ApprovalDate              time.Time `json:"approval_date"`
}

type StockOpnameUpdateRequest struct {
	CompanyId                 int       `json:"company_id"`
	StockOpnameSystemNumber   int       `json:"stock_opname_system_number"`
	StockOpnameDocumentNumber string    `json:"stock_opname_document_number"`
	WarehouseGroup            int       `json:"warehouse_group"`
	WarehouseCode             int       `json:"warehouse_code"`
	FromLocation              int       `json:"from_location"`
	ToLocation                int       `json:"to_location"`
	ItemGroup                 int       `json:"item_group"`
	ShowDetail                bool      `json:"show_details"`
	PersonInCharge            int       `json:"person_in_charge"`
	Remark                    string    `json:"remark"`
	ExecutionDateFrom         time.Time `json:"execution_date_from"`
	ExecutionDateTo           time.Time `json:"execution_date_to"`
	IncludeZeroOnhand         bool      `json:"include_zero_onhand"`
	ApprovalRequestedById     int       `json:"approval_requested_by_id"`
	ApprovalRequestedDate     time.Time `json:"approval_requested_date"`
	ApprovalById              int       `json:"approval_by_id"`
	ApprovalDate              time.Time `json:"approval_date"`
}

type GetAllStockOpnameResponse struct {
	StockOpnameDocumentNumber  string    `json:"stock_opname_document_number"`
	ExecutionDateFrom          time.Time `json:"execution_date_from"`
	ExecutionDateTo            time.Time `json:"execution_date_to"`
	WarehouseLocationGroupName string    `json:"warehouse_location_group_name"`
	WarehouseName              string    `json:"warehouse_name"`
	StockOpnameStatus          string    `json:"stock_opname_status"`
}

type GetAllStockOpnameDetailResponse struct {
	ItemId         int     `json:"item_id"`
	ItemName       string  `json:"item_name"`
	Location       string  `json:"location"`
	OnHandQuantity float64 `json:"on_hand_quantity"`
}

type GetStockOpnameByStockOpnameSystemNumberResponse struct {
	StockOpnameApprovalStatus string    `json:"stock_opname_approval_status"`
	StockOpnameDocumentNumber string    `json:"stock_opname_document_number"`
	WarehouseGroup            string    `json:"warehouse_group"`
	WarehouseCode             string    `json:"warehouse_code"`
	FromLocation              string    `json:"from_location"`
	ToLocation                string    `json:"to_location"`
	ItemGroup                 string    `json:"item_group"`
	ShowDetail                bool      `json:"show_detail"`
	PersonInCharge            string    `json:"person_in_charge"`
	Remark                    string    `json:"remark"`
	ExecutionDateFrom         time.Time `json:"execution_date_from"`
	ExecutionDateTo           time.Time `json:"execution_date_to"`
	IncludeZeroOnhand         bool      `json:"include_zero_onhand"`
}

type StockOpnameSubmitRequest struct {
	StockOpnameApprovalRequestId int `json:"stock_opname_approval_request_id" validate:"required"`
}

type StockOpnameInsertDetailRequest struct {
	ItemId []int `json:"item_id"`
}

type StockOpnameUpdateDetailRequest struct {
	ItemId          int `json:"item_id"`
	StockOpnameLine int `json:"stock_opname_line"`
}

// type GetAllStockOpnameResponse struct {
// 	StockOpnameNo   string    `json:"stock_opname_no"`
// 	StockOpnameFrom time.Time `json:"stock_opname_from"`
// 	StockOpnameTo   time.Time `json:"stock_opname_to"`
// 	WarehouseGroup  string    `json:"warehouse_group"`
// 	WarehouseName   string    `json:"warehouse_name"`
// 	Status          string    `json:"status"`
// }

// type GetAllLocationList struct {
// 	LocationCode string `json:"location_code"`
// 	LocationName string `json:"location_name"`
// 	Status       string `json:"status"`
// }

// type GetPersonInChargeResponse struct {
// 	EmployeeNo   string `json:"employee_no"`
// 	EmployeeName string `json:"employee_name"`
// 	Position     string `json:"position"`
// 	Status       string `json:"status"`
// }

// type GetItemListResponse struct {
// 	ItemCode        string  `json:"item_code"`
// 	ItemName        string  `json:"item_name"`
// 	Location        string  `json:"location"`
// 	OnHandQuantity  float64 `json:"on_hand_quantity"`
// 	StockOpnameLine float64 `json:"stock_opname_line"`
// }

// type GetOnGoingStockOpnameResponse struct {
// 	Status              string      `json:"status"`
// 	StockOpnameSysNo    string      `json:"stock_opname_sys_no"`
// 	WarehouseGroup      string      `json:"warehouse_group"`
// 	WarehouseCode       string      `json:"warehouse_code"`
// 	FromLocation        string      `json:"from_location"`
// 	ToLocation          string      `json:"to_location"`
// 	ItemGroup           string      `json:"item_group"`
// 	PersonInCharge      string      `json:"person_in_charge"`
// 	StockOpnameDateFrom time.Time   `json:"stock_opname_date_from"`
// 	StockOpnameDateTo   time.Time   `json:"stock_opname_date_to"`
// 	GetItemListResponse interface{} `json:"get_item_list_response"`
// }

// type InsertNewStockOpnameRequest struct {
// 	Status string `json:"record_status"`
// 	// StockOpnameSysNo    string    `json:"stock_opname_sys_no"`
// 	StockOpnameDocNo    string    `json:"stock_opname_doc_no"`
// 	CompanyCode         float64   `json:"company_code"`
// 	StockOpnameStatus   string    `json:"stock_opname_status"`
// 	WarehouseGroup      string    `json:"warehouse_group" validate:"required"`
// 	WarehouseCode       string    `json:"warehouse_code" validate:"required"`
// 	FromLocation        string    `json:"from_location" validate:"required"`
// 	ToLocation          string    `json:"to_location" validate:"required"`
// 	ItemGroup           string    `json:"item_group" validate:"required"`
// 	PersonInCharge      string    `json:"person_in_charge" validate:"required"`
// 	StockOpnameDateFrom time.Time `json:"stock_opname_date_from" validate:"required"`
// 	StockOpnameDateTo   time.Time `json:"stock_opname_date_to" validate:"required"`
// 	Remark              string    `json:"remark"`
// 	UserIdCreated       string    `json:"user_id_created"`
// 	TotalAdjCost        float64   `json:"total_adj_cost"`
// }
