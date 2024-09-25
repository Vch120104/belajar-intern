package transactionworkshopentities

import "time"

const TableNameBookingEstimationDetail = "trx_booking_estimation_detail"

type BookingEstimationDetail struct {
	EstimationLineID               int        `gorm:"column:estimation_line_id;size:30;primaryKey" json:"estimation_line_id"`
	EstimationLineCode             int        `gorm:"column:estimation_line_code;size:30;default:null" json:"estimation_line_code"`
	EstimationSystemNumber         int        `gorm:"column:estimation_system_number;size:30;default:null" json:"estimation_system_number"`
	BillID                         int        `gorm:"column:bill_id;size:30;default:null" json:"bill_id"`
	EstimationLineDiscountApproval int        `gorm:"column:estimation_line_discount_approval_status;size:30;default:null" json:"estimation_line_discount_approval_status"`
	ItemOperationID                int        `gorm:"column:item_operation_id;size:30;default:null" json:"item_operation_id"`
	LineTypeID                     int        `gorm:"column:line_type_id;size:30;default:null" json:"line_type_id"`
	PackageID                      int        `gorm:"column:package_id;size:30;default:null" json:"package_id"`
	JobTypeID                      int        `gorm:"column:job_type_id;size:30;default:null" json:"job_type_id"`
	FieldActionSystemNumber        int        `gorm:"column:field_action_system_number;size:30;default:null" json:"field_action_system_number"`
	ApprovalRequestNumber          int        `gorm:"column:approval_request_number;size:30;default:null" json:"approval_request_number"`
	UOMID                          int        `gorm:"column:uom_id;size:30;default:null" json:"uom_id"`
	RequestDescription             string     `gorm:"column:request_description;default:null" json:"request_description"`
	FRTQuantity                    float64    `gorm:"column:frt_quantity;default:null" json:"frt_quantity"`
	ItemOperationPrice             float64    `gorm:"column:item_operation_price;default:null" json:"item_operation_price"`
	DiscountItemOperationAmount    float64    `gorm:"column:discount_item_operation_amount;default:null" json:"discount_item_operation_amount"`
	DiscountItemOperationPercent   float64    `gorm:"column:discount_item_operation_percent;default:null" json:"discount_item_operation_percent"`
	DiscountRequestPercent         float64    `gorm:"column:discount_request_percent;default:null" json:"discount_request_percent"`
	DiscountRequestAmount          float64    `gorm:"column:discount_request_amount;default:null" json:"discount_request_amount"`
	DiscountApprovalBy             string     `gorm:"column:discount_approval_by;size:10;default:null" json:"discount_approval_by"`
	DiscountApprovalDate           *time.Time `gorm:"column:discount_approval_date;default:null" json:"discount_approval_date"`
}

func (*BookingEstimationDetail) TableName() string {
	return TableNameBookingEstimationDetail
}
