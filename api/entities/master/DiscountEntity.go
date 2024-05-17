package masterentities

import masteritementities "after-sales/api/entities/master/item"

var CreateDiscountTable = "mtr_discount"

type Discount struct {
	IsActive                bool                               `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DiscountCodeId          int                                `gorm:"column:discount_code_id;not null;primaryKey"        json:"discount_code_id"`
	DiscountCodeValue       string                             `gorm:"column:discount_code_value;size:20;not null"        json:"discount_code_value"`
	DiscountCodeDescription string                             `gorm:"column:discount_code_description;size:50;not null"      json:"discount_code_description"`
	DiscountPercent         masteritementities.DiscountPercent `gorm:"foreignKey:DiscountCodeId;references:DiscountCodeId"`
}

func (*Discount) TableName() string {
	return CreateDiscountTable
}
