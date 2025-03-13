package transactionsparepartentities

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

var CreateStockOpnameTable = "trx_stock_opname"

type StockOpname struct {
	CompanyID                   *int       `gorm:"column:company_id;size:30;null" json:"company_id"`
	StockOpnameSystemNumber     int        `gorm:"column:stock_opname_system_number;primaryKey;size:30;not null" json:"stock_opname_system_number"`
	StockOpnameDocumentNumber   string     `gorm:"column:stock_opname_document_number;size:25;null" json:"stock_opname_document_number"`
	StockOpnameStatusId         *int       `gorm:"column:stock_opname_status_id;size:30;null" json:"stock_opname_status_id"`
	WarehouseGroupId            *int       `gorm:"column:warehouse_group_id;size:30;null" json:"warehouse_group_id"`
	WarehouseId                 *int       `gorm:"column:warehouse_id;size:30;null" json:"warehouse_id"`
	LocationRangeFromId         *int       `gorm:"column:location_range_from_id;size:30;null" json:"location_range_from_id"`
	LocationRangeToId           *int       `gorm:"column:location_range_to_id;size:30;null" json:"location_range_to_id"`
	ProfitCenterId              *int       `gorm:"column:profit_center_id;size:30;null" json:"profit_center_id"`
	TransactionTypeId           string     `gorm:"column:transaction_type_id;size:10;null" json:"transaction_type_id"`
	ShowDetail                  bool       `gorm:"column:show_detail;null" json:"show_detail"`
	PersonInChargeId            *int       `gorm:"column:person_in_charge_id;size:30;null" json:"person_in_charge_id"`
	Remark                      string     `gorm:"column:remark;size:256;null" json:"remark"`
	ItemGroupId                 *int       `gorm:"column:item_group_id;size:30;null" json:"item_group_id"`
	ExecutionDateFrom           * time.Time `gorm:"column:execution_date_from;null" json:"execution_date_from"`
	ExecutionDateTo             *time.Time `gorm:"column:execution_date_to;null" json:"execution_date_to"`
	BrokenWarehouseId           *int       `gorm:"column:broken_warehouse_id;size:30;null" json:"broken_warehouse_id"`
	AdjustmentDate              time.Time  `gorm:"column:adjustment_date;null" json:"adjustment_date"`
	StockOpnameApprovalStatusId *int       `gorm:"column:stock_opname_approval_status_id;size:30;null" json:"stock_opname_approval_status_id"`
	ApprovalRequestedById       *int       `gorm:"column:approval_requested_by_id;size:30;null" json:"approval_requested_by_id"`
	ApprovalRequestedDate       *time.Time `gorm:"column:approval_requested_date;null" json:"approval_requested_date"`
	ApprovalById                *int       `gorm:"column:approval_by_id;size:30;null" json:"approval_by_id"`
	ApprovalDate                *time.Time `gorm:"column:approval_date;null" json:"approval_date"`
	TotalAdjustmentCost         float64    `gorm:"column:total_adjustment_cost;type:decimal(17,4);null" json:"total_adjustment_cost"`
	IncludeZeroOnhand           bool       `gorm:"column:include_zero_onhand;null" json:"include_zero_onhand"`

	StockOpnameStatus masterentities.StockOpnameStatus          `gorm:"foreignKey:stock_opname_status_id;references:stock_opname_status_id" json:"stock_opname_status"`
	WarehouseGroup    masterwarehouseentities.WarehouseGroup    `gorm:"foreignKey:warehouse_group_id;references:warehouse_group_id" json:"warehouse_group"`
	Warehouse         masterwarehouseentities.WarehouseMaster   `gorm:"foreignKey:warehouse_id;references:warehouse_id" json:"warehouse"`
	LocationRangeFrom masterwarehouseentities.WarehouseLocation `gorm:"foreignKey:location_range_from_id;references:warehouse_location_code" json:"location_range_from"`
	LocationRangeTo   masterwarehouseentities.WarehouseLocation `gorm:"foreignKey:location_range_to_id;references:warehouse_location_code" json:"location_range_to"`
	ItemGroup         masteritementities.ItemGroup              `gorm:"foreignKey:item_group_id;references:item_group_id" json:"item_group"`
	BrokenWarehouse   masterwarehouseentities.WarehouseMaster   `gorm:"foreignKey:broken_warehouse_id;references:warehouse_id" json:"broken_warehouse"`
}

func (*StockOpname) TableName() string {
	return CreateStockOpnameTable
}
