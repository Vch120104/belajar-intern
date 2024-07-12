package transactionworkshopentities

import "time"

const TableNameBookingEstimationDetail = "trx_booking_estimation_detail"

type BookingEstimationDetail struct {
	EstimationLineID               int        `gorm:"column:estimation_line_id;size:30;primaryKey" json:"estimation_line_id"`
	EstimationLineCode             int        `gorm:"column:estimation_line_code;size:30;default:null" json:"estimation_line_code"`
	EstimationSystemNumber         int        `gorm:"column:estimation_system_number;size:30;default:null" json:"estimation_system_number"`
	BillID                         string     `gorm:"column:bill_id;default:null" json:"bill_id"`
	EstimationLineDiscountApproval int        `gorm:"column:estimation_line_discount_approval_status;size:30;default:null" json:"estimation_line_discount_approval_status"`
	ItemID                         int        `gorm:"column:item_id;size:30;default:null" json:"item_id"`
	LineTypeID                     int        `gorm:"column:line_type_id;size:30;default:null" json:"line_type_id"`
	PackageID                      int        `gorm:"column:package_id;size:30;default:null" json:"package_id"`
	JobTypeID                      int        `gorm:"column:job_type_id;size:30;default:null" json:"job_type_id"`
	FieldActionSystemNumber        int        `gorm:"column:field_action_system_number;size:30;default:null" json:"field_action_system_number"`
	ApprovalRequestNumber          int        `gorm:"column:approval_request_number;size:30;default:null" json:"approval_request_number"`
	UOMID                          int        `gorm:"column:uom_id;size:30;default:null" json:"uom_id"`
	RequestDescription             string     `gorm:"column:request_description;default:null" json:"request_description"`
	FRTQuantity                    float32    `gorm:"column:frt_quantity;default:null" json:"frt_quantity"`
	OperationItemPrice             float32    `gorm:"column:operation_item_price;default:null" json:"operation_item_price"`
	DiscountItemAmount             float32    `gorm:"column:discount_item_amount;default:null" json:"discount_item_amount"`
	DiscountItemPercent            float32    `gorm:"column:discount_item_percent;default:null" json:"discount_item_percent"`
	DiscountRequestPercent         float32    `gorm:"column:discount_request_percent;default:null" json:"discount_request_percent"`
	DiscountRequestAmount          float32    `gorm:"column:discount_request_amount;default:null" json:"discount_request_amount"`
	DiscountApprovalBy             string     `gorm:"column:discount_approval_by;size:10;default:null" json:"discount_approval_by"`
	DiscountApprovalDate           *time.Time `gorm:"column:discount_approval_date;default:null" json:"discount_approval_date"`
}

func (*BookingEstimationDetail) TableName() string {
	return TableNameBookingEstimationDetail
}
