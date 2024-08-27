package masterentities

const itemCycletableName = "mtr_item_cycle"

type ItemCycle struct {
	IsActive          bool    `gorm:"column:is_active;not null;default:true" json:"is_active"`
	ItemCycleId       int     `gorm:"column:item_cycle_id;not null;primaryKey" json:"item_cycle_id"`
	CompanyId         int     `gorm:"column:company_id;not null;uniqueIndex:item_cycle" json:"company_id"`
	PeriodYear        string  `gorm:"column:period_year;not null;uniqueIndex:item_cycle" json:"period_year"`
	PeriodMonth       string  `gorm:"column:period_month;not null;uniqueIndex:item_cycle" json:"period_month"`
	ItemId            int     `gorm:"column:item_id;not null;uniqueIndex:item_cycle" json:"item_id"`
	OrderCycle        int     `gorm:"column:order_cycle" json:"order_cycle"`
	QuantityOnOrder   float64 `gorm:"column:quantity_on_order" json:"quantity_on_order"`
	QuantityBackOrder float64 `gorm:"column:quantity_back_order" json:"quantity_back_order"`
}

func (*ItemCycle) TableName() string {
	return itemCycletableName
}
