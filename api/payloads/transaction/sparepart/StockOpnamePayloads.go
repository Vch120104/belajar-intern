package transactionsparepartpayloads

import "time"

// type GetAllStockOpnameResponse struct {
// 	StockOpnameNo     string    `json:"stock_opname_no"`
// 	StockOpnameStatus string    `json:"stock_opname_status"`
// 	StockOpnameFrom   time.Time `json:"stock_opname_from"`
// 	StockOpnameTo     time.Time `json:"stock_opname_to"`
// 	WarehouseGroup    string    `json:"warehouse_group"`
// 	WarehouseCode     string    `json:"warehouse_code"`
// }

type GetAllStockOpnameResponse struct {
	StockOpnameNo   string    `json:"stock_opname_no"`
	StockOpnameFrom time.Time `json:"stock_opname_from"`
	StockOpnameTo   time.Time `json:"stock_opname_to"`
	WarehouseGroup  string    `json:"warehouse_group"`
	WarehouseName   string    `json:"warehouse_name"`
	Status          string    `json:"status"`
}

type GetAllLocationList struct {
	LocationCode string `json:"location_code"`
	LocationName string `json:"location_name"`
	Status       string `json:"status"`
}

type GetPersonInChargeResponse struct {
	EmployeeNo   string `json:"employee_no"`
	EmployeeName string `json:"employee_name"`
	Position     string `json:"position"`
	Status       string `json:"status"`
}

type GetItemListResponse struct {
	ItemCode        string  `json:"item_code"`
	ItemName        string  `json:"item_name"`
	Location        string  `json:"location"`
	OnHandQuantity  float64 `json:"on_hand_quantity"`
	StockOpnameLine float64 `json:"stock_opname_line"`
}

type GetOnGoingStockOpnameResponse struct {
	Status              string      `json:"status"`
	StockOpnameSysNo    string      `json:"stock_opname_sys_no"`
	WarehouseGroup      string      `json:"warehouse_group"`
	WarehouseCode       string      `json:"warehouse_code"`
	FromLocation        string      `json:"from_location"`
	ToLocation          string      `json:"to_location"`
	ItemGroup           string      `json:"item_group"`
	PersonInCharge      string      `json:"person_in_charge"`
	StockOpnameDateFrom time.Time   `json:"stock_opname_date_from"`
	StockOpnameDateTo   time.Time   `json:"stock_opname_date_to"`
	GetItemListResponse interface{} `json:"get_item_list_response"`
}

type InsertNewStockOpnameRequest struct {
	Status string `json:"record_status"`
	// StockOpnameSysNo    string    `json:"stock_opname_sys_no"`
	StockOpnameDocNo    string    `json:"stock_opname_doc_no"`
	CompanyCode         float64   `json:"company_code"`
	StockOpnameStatus   string    `json:"stock_opname_status"`
	WarehouseGroup      string    `json:"warehouse_group" validate:"required"`
	WarehouseCode       string    `json:"warehouse_code" validate:"required"`
	FromLocation        string    `json:"from_location" validate:"required"`
	ToLocation          string    `json:"to_location" validate:"required"`
	ItemGroup           string    `json:"item_group" validate:"required"`
	PersonInCharge      string    `json:"person_in_charge" validate:"required"`
	StockOpnameDateFrom time.Time `json:"stock_opname_date_from" validate:"required"`
	StockOpnameDateTo   time.Time `json:"stock_opname_date_to" validate:"required"`
	Remark              string    `json:"remark"`
	UserIdCreated       string    `json:"user_id_created"`
	TotalAdjCost        float64   `json:"total_adj_cost"`
}
