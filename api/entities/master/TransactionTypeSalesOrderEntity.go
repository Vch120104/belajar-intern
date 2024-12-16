package masterentities

const TransactionTypeSalesOrderTableName = "mtr_transaction_type_sales_order"

type TransactionTypeSalesOrder struct {
	TransactionTypeSalesOrderId          int    `gorm:"column:transaction_type_sales_order_id;not null;primaryKey;size:30" json:"transaction_type_sales_order_id"`
	TransactionTypeSalesOrderCode        string `gorm:"column:transaction_type_sales_order_code;not null" json:"transaction_type_sales_order_code"`
	TransactionTypeSalesOrderDescription string `gorm:"column:transaction_type_sales_order_description;not null" json:"transaction_type_sales_order_description"`
}

func (*TransactionTypeSalesOrder) TableName() string {
	return TransactionTypeSalesOrderTableName
}
