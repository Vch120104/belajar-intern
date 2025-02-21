package transactionworkshopentities

const TableNameBookingEstimationDetail = "trx_booking_estimation_detail"

type BookingEstimationDetail struct {
	EstimationDetailId                  int     `gorm:"column:estimation_detail_id;size:30;primaryKey" json:"estimation_detail_id"`
	EstimationLine                      int     `gorm:"column:estimation_line;size:30;default:null" json:"estimation_line"`
	EstimationSystemNumber              int     `gorm:"column:estimation_system_number;size:30;default:null" json:"estimation_system_number"`
	EstimationDocumentNumber            string  `gorm:"column:estimation_document_number;" json:"estimation_document_number"`
	TransactionTypeId                   int     `gorm:"column:transaction_type_id;size:30" json:"transaction_type_id"`
	JobTypeId                           int     `gorm:"column:job_type_id;size:30;default:null" json:"job_type_id"`
	EstimationLineDiscountStatus        int     `gorm:"column:estimation_line_discount_status;size:30;" json:"estimation_line_discount_status"`
	OperationItemId                     int     `gorm:"column:operation_item_id;size:30;" json:"operation_item_id"`
	OperationItemCode                   string  `gorm:"column:operation_item_code;size:30;" json:"operation_item_code"`
	LineTypeId                          int     `gorm:"column:line_type_id;size:30;default:null" json:"line_type_id"`
	PackageId                           int     `gorm:"column:package_id;size:30;default:null" json:"package_id"`
	FRTQuantity                         float64 `gorm:"column:frt_quantity;" json:"frt_quantity"`
	OperationItemPrice                  float64 `gorm:"column:operation_item_price;" json:"operation_item_price"`
	OperationItemDiscountAmount         float64 `gorm:"column:operation_item_discount_amount;" json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount  float64 `gorm:"column:operation_item_discount_request_amount;default:null" json:"operation_item_discount_request_amount"`
	OperationItemDiscountPercent        float64 `gorm:"column:operation_item_discount_percent;default:null" json:"operation_item_discount_percent"`
	OperationItemDiscountRequestPercent float64 `gorm:"column:operation_item_discount_request_percent;default:null" json:"operation_item_discount_request_percent"`
	FieldActionSystemNumber             int     `gorm:"column:field_action_system_number;size:30;default:null" json:"field_action_system_number"`
}

func (*BookingEstimationDetail) TableName() string {
	return TableNameBookingEstimationDetail
}
