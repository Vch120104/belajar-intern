package masterentities

var CreateStockOpnameStatusTable = "mtr_stock_opname_status"

type StockOpnameStatus struct {
	StockOpnameStatusId   int    `gorm:"column:stock_opname_status_id;primaryKey;size:30;not null" json:"stock_opname_status_id"`
	StockOpnameStatusCode string `gorm:"column:stock_opname_status_code;size:25;not null" json:"stock_opname_status_code"`
	StockOpnameStatusName string `gorm:"column:stock_opname_status_name;size:100;null" json:"stock_opname_status_name"`
}

func (*StockOpnameStatus) TableName() string {
	return CreateStockOpnameStatusTable
}
