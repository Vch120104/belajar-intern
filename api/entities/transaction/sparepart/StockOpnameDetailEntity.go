package transactionsparepartentities

var CreateStockOpnameDetailTable = "trx_stock_opname_detail"

type StockOpnameDetail struct {
	StockOpnameDetailSystemNumber int     `gorm:"column:stock_opname_detail_system_number;primaryKey;size:30;not null" json:"stock_opname_detail_number"`
	StockOpnameSystemNumber       int     `gorm:"column:stock_opname_system_number;size:30;index:idx_stock_opname_detail,unique;not null" json:"stock_opname_system_number"`
	StockOpnameLine               int     `gorm:"column:stock_opname_line;size:30;index:idx_stock_opname_detail,unique;not null" json:"stock_opname_line"`
	WarehouseId                   *int    `gorm:"column:warehouse_id;size:30;null" json:"warehouse_id"`
	LocationId                    *int    `gorm:"column:location_id;size:30;null" json:"location_id"`
	ItemId                        *int    `gorm:"column:item_id;size:30;null" json:"item_id"`
	SystemQuantity                float64 `gorm:"column:system_quantity;type:decimal(10,2);null" json:"system_quantity"`
	FoundQuantity                 float64 `gorm:"column:found_quantity;type:decimal(10,2);null" json:"found_quantity"`
	BrokenQuantity                float64 `gorm:"column:broken_quantity;type:decimal(10,2);null" json:"broken_quantity"`
	NeedAdjustment                bool    `gorm:"column:need_adjustment;null" json:"need_adjustment"`
	Remark                        string  `gorm:"column:remark;size:256;null" json:"remark"`
	Cogs                          float64 `gorm:"column:COGS;type:decimal(17,4);null" json:"cogs"`
	AdjustmentCost                float64 `gorm:"column:adjustment_cost;type:decimal(17,4);null" json:"adjustment_cost"`
	AllocatedQuantity             float64 `gorm:"column:allocated_quantity;type:decimal(10,2);null" json:"allocated_quantity"`

	StockOpnameSystem *StockOpname `gorm:"foreignKey:StockOpnameSystemNumber;references:stock_opname_system_number" json:"stock_opname_system"`
}

func (*StockOpnameDetail) TableName() string {
	return CreateStockOpnameDetailTable
}
