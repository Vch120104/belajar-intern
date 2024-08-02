package masteritementities

var CreateItemImportTable = "mtr_item_import"

type ItemImport struct {
	ItemImportId       int     `gorm:"column:item_import_id;size:30;not null;primaryKey"        json:"item_import_id"`
	SupplierId         int     `gorm:"column:supplier_id;not null;size:30;uniqueindex:idx_item_import"        json:"supplier_id"`
	ItemId             int     `gorm:"column:item_id;not null;size:30;uniqueindex:idx_item_import"        json:"item_id"`
	OrderQtyMultiplier float64 `gorm:"column:order_qty_multiplier;not null"        json:"order_qty_multiplier"`
	ItemAliasCode      string  `gorm:"column:item_alias_code;not null;size:30"        json:"item_alias_code"`
	RoyaltyFlag        string  `gorm:"column:royalty_flag;not null;size:30"        json:"royalty_flag"`
	ItemAliasName      string  `gorm:"column:item_alias_name;not null;size:100"        json:"item_alias_name"`
	OrderConversion    float64 `gorm:"column:order_conversion;not null"        json:"order_conversion"`
}

func (*ItemImport) TableName() string {

	return CreateItemImportTable
}
