package masterentities

const POTypeSoTableName = "mtr_purchase_order_type_sales_order"

type PurchaseOrderTypeSalesOrderEntity struct {
	PurchaseOrderTypeSalesOrderId          int    `gorm:"column:purchase_order_type_sales_order_id;not null;primaryKey;size:30" json:"purchase_order_type_sales_order_id"`
	PurchaseOrderTypeSalesOrderCode        string `gorm:"column:purchase_order_type_sales_order_code;not null" json:"purchase_order_type_sales_order_code"`
	PurchaseOrderTypeSalesOrderDescription string `gorm:"column:purchase_order_type_sales_order_description;not null" json:"purchase_order_type_sales_order_description"`
}

func (*PurchaseOrderTypeSalesOrderEntity) TableName() string {
	return POTypeSoTableName
}
