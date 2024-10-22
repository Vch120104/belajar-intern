package transactionworkshopentities

const TableNameRecallDetailType = "trx_recall_detail_type"

type RecallDetailType struct {
	RecallDetailTypeId     int     `gorm:"column:recall_detail_type_id;size:30;primaryKey" json:"recall_detail_type_id"`
	RecallSystemNumber     int     `gorm:"column:recall_system_number;size:30;" json:"recall_system_number"`
	RecallLineNumber       int     `gorm:"column:recall_line_number;size:30;" json:"recall_line_number"`
	RecallTypeId           int     `gorm:"column:recall_type_id;size:30;" json:"recall_type_id"`
	RecallDetailTypeNumber int     `gorm:"column:recall_detail_type_number;size:30;" json:"recall_detail_type_number"`
	OperationItemId        int     `gorm:"column:operation_item_id;size:30;" json:"operation_item_id"`
	OperationItemCode      string  `gorm:"column:operation_item_code;" json:"operation_item_code"`
	FrtQty                 float64 `gorm:"column:frt_qty;" json:"frt_qty"`
	HasRecall              bool    `gorm:"column:has_recall;" json:"has_recall"`
	IsActive               bool    `gorm:"column:is_active;" json:"is_active"`
}

func (*RecallDetailType) TableName() string {
	return TableNameRecallDetailType
}
