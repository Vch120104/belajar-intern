package masterentities

import masteritementities "after-sales/api/entities/master/item"

const TableNameGmmDiscountSetting = "mtr_gmm_discount_setting"

type GmmDiscountSetting struct {
	IsActive             bool                          `gorm:"column:is_active;default:false;not null" json:"is_active"`
	GmmDiscountSettingId int                           `gorm:"column:gmm_discount_setting_id;size:30;primaryKey" json:"gmm_discount_setting_id"`
	GmmPriceCodeId       int                           `gorm:"column:gmm_price_code_id;size:30;not null;uniqueindex:idx_gmm_discount_setting" json:"gmm_price_code_id"`
	GmmPriceCode         GmmPriceCode                  `gorm:"foreignKey:GmmPriceCodeId;references:GmmPriceCodeId"`
	ItemLevel1Id         int                           `gorm:"column:item_level_1_id;size:30;not null;uniqueindex:idx_gmm_discount_setting" json:"item_level_1_id"`
	ItemLevel1           masteritementities.ItemLevel1 `gorm:"foreignKey:ItemLevel1Id;references:ItemLevel1Id"`
	OrderTypeId          int                           `gorm:"column:order_type_id;size:30;not null;uniqueindex:idx_gmm_discount_setting" json:"order_type_id"`
	OrderType            OrderType                     `gorm:"foreignKey:OrderTypeId;references:OrderTypeId"`
	DiscountPercentage   float64                       `gorm:"column:discount_percentage;not null" json:"discount_percentage"`
}

func (*GmmDiscountSetting) TableName() string {
	return TableNameGmmDiscountSetting
}
