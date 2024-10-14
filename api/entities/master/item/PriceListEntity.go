package masteritementities

import "time"

var CreatePriceListTable = "mtr_item_price_list"

type PriceList struct {
	IsActive            bool      `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	PriceListId         int       `gorm:"column:price_list_id;not null;primaryKey;size:30"        json:"price_list_id"`
	PriceListCodeId     int       `gorm:"column:price_list_code_id;size:30;null"        json:"price_list_code_id"`
	CompanyId           int       `gorm:"column:company_id;not null;size:30"        json:"company_id"`
	BrandId             int       `gorm:"column:bran2d_id;not null;size:30;uniqueindex:idx_price_list"        json:"brand_id"`
	CurrencyId          int       `gorm:"column:currency_id;not null;size:30;uniqueindex:idx_price_list"        json:"currency_id"`
	EffectiveDate       time.Time `gorm:"column:effective_date;not null;uniqueindex:idx_price_list"        json:"effective_date"`
	ItemId              int       `gorm:"column:item_id;not null;size:30;uniqueindex:idx_price_list"        json:"item_id"`
	ItemGroupId         int       `gorm:"column:item_group_id;not null;size:30;uniqueindex:idx_price_list"        json:"item_group_id"`
	ItemClassId         int       `gorm:"column:item_class_id;not null;size:30"        json:"item_class_3id"`
	ItemClass           *ItemClass
	PriceListAmount     float64   `gorm:"column:price_list_amount;size:17,4;not null"        json:"price_list_amount"`
	PriceListModifiable bool      `gorm:"column:price_list_modifiable;null"        json:"price_list_modifiable"`
	AtpmSyncronize      bool      `gorm:"column:atpm_syncronize;null"        json:"atpm_syncronize"`
	AtpmSyncronizeTime  time.Time `gorm:"column:atpm_syncronize_time;null"        json:"atpm_syncronize_time"`
}

func (*PriceList) TableName() string {
	return CreatePriceListTable
}
