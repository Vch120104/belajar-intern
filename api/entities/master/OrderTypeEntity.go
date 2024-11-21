package masterentities

const TableNameOrderType = "mtr_order_type"

type OrderType struct {
	IsActive      bool   `gorm:"column:is_active;default:false;not null" json:"is_active"`
	OrderTypeId   int    `gorm:"column:order_type_id;size:30;primaryKey" json:"order_type_id"`
	OrderTypeCode string `gorm:"column:order_type_code;size:5;unique;not null" json:"order_type_code"`
	OrderTypeName string `gorm:"column:order_type_name;size:100;not null" json:"order_type_name"`
}

func (*OrderType) TableName() string {
	return TableNameOrderType
}
