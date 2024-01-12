package masteritementities

import masterentities "after-sales/api/entities/master"

var CreateDiscountPercentTable = "mtr_discount_percent"

type DiscountPercent struct {
	IsActive          bool                    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DiscountPercentId int                     `gorm:"column:discount_percent_id;not null;primaryKey"        json:"discount_percent_id"`
	DiscountCodeId    int                     `gorm:"column:discount_code_id;not null"        json:"discount_code_id"`
	DiscountMaster    masterentities.Discount `gorm:"foreignKey:DiscountCodeId"`
	OrderTypeId       int                     `gorm:"column:order_type_id;not null"        json:"order_type_id"` //Fk dari mtr_order_type general service
	Discount          float64                 `gorm:"column:discount;not null"        json:"discount"`
}

func (*DiscountPercent) TableName() string {
	return CreateDiscountPercentTable
}
