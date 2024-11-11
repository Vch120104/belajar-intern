package masterentities

var StockTransactionReasonTableName = "mtr_stock_transaction_reason"

type StockTransactionReason struct {
	StockTransactionReasonId    int    `gorm:"column:stock_transaction_reason_id;not null;primaryKey;size:30" json:"stock_transaction_reason_id"`
	StockTransactionReasonCode  string `gorm:"column:stock_transaction_reason_code;not null" json:"stock_transaction_reason_code"`
	StockTransactionDescription string `gorm:"column:stock_transaction_reason_description;not null" json:"stock_transaction_reason_description"`
}

func (*StockTransactionReason) TableName() string {
	return StockTransactionReasonTableName
}
