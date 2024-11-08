package masterentities

var CreateGroupStockTable = "mtr_group_stock"

type GroupStock struct {
	GroupStockId     int     `gorm:"column:group_stock_id;primaryKey;size:30" json:"group_stock_id"`
	CompanyId        int     `gorm:"column:company_id;size:30" json:"company_id"`
	PeriodYear       string  `gorm:"column:period_year;size:30" json:"period_year"`
	PeriodMonth      string  `gorm:"column:period_month;size:30" json:"period_month"`
	WhsGroup         string  `gorm:"column:whs_group;size:30" json:"whs_group"`
	WarehouseGroupId int     `gorm:"column:warehouse_group_id;size:30" json:"warehouse_group_id"`
	ItemId           int     `gorm:"column:item_id;size:30" json:"item_id"`
	ItemCode         string  `gorm:"column:item_code;size:100" json:"item_code"`
	PriceBegin       float64 `gorm:"column:price_begin" json:"price_begin"`
	PriceCurrent     float64 `gorm:"column:price_current" json:"price_current"`
}

func (*GroupStock) TableName() string {
	return CreateGroupStockTable
}
