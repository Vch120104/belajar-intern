package masteritementities

var CreateItemPriceCodeTable = "mtr_item_price_code"

type ItemPriceCode struct {
	ItemPriceCodeId   int    `gorm:"column:item_price_code_id;not null;primaryKey;size:30;" json:"item_price_code__id"`
	IsActive          bool   `gorm:"column:is_active;not null"        json:"is_active"`
	ItemPriceCode     string `gorm:"column:item_price_code;not null;unique;size:20"        json:"item_price_code"`
	ItemPriceCodeName string `gorm:"column:item_price_code_name;not null;size:256"        json:"item_price_code_name"`
}
