package masterentities

var CreateDiscountTable = "mtr_discount"

type Discount struct {
	IsActive                bool   `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DiscountCodeId          int    `gorm:"column:discount_code_id;size:30;not null;primaryKey"        json:"discount_code_id"`
	DiscountCodeValue       string `gorm:"column:discount_code_value;size:20;not null"        json:"discount_code_value"`
	DiscountCodeDescription string `gorm:"column:discount_code_description;size:50;not null"      json:"discount_code_description"`
}

func (*Discount) TableName() string {
	return CreateDiscountTable
}
