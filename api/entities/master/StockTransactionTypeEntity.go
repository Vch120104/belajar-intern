package masterentities

var StockTransactionEntity = "mtr_stock_transaction_type"

type StockTransactionType struct {
	StockTransactionTypeId      int    `gorm:"column:stock_transaction_type_id;not null;primaryKey;size:30" json:"stock_transaction_type_id"`
	StockTransactionTypeCode    string `gorm:"column:stock_transaction_type_code;not null" json:"stock_transaction_type_code"`
	StockTransactionDescription string `gorm:"column:stock_transaction_description;not null" json:"stock_transaction_description"`
}

func (*StockTransactionType) TableName() string {
	return StockTransactionEntity
}
