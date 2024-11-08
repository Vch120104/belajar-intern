package masteritementities

var CreateDiscountTable = "mtr_discount"

type Discount struct {
	DiscountCodeId      int               `gorm:"column:discount_code_id;size:30;not null;primaryKey"        json:"discount_code_id"`
	IsActive            bool              `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	DiscountCode        string            `gorm:"column:discount_code;size:20;not null"        json:"discount_code"`
	DiscountDescription string            `gorm:"column:discount_description;size:50;not null"      json:"discount_description"`
	Discounts           []DiscountPercent `gorm:"foreignKey:DiscountCodeId;" json:"discount_percent"`
}

func (*Discount) TableName() string {
	return CreateDiscountTable
}
