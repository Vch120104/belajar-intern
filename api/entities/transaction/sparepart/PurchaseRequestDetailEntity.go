package transactionsparepartentities

import "time"

const TableNamePurchaseRequestDetail = "trx_purchase_request_detail"

type PurchaseRequestDetail struct {
	PurchaseRequestDetailSystemNumber int        `gorm:"column:purchase_request_detail_system_number;size:50;not null;primaryKey;" json:"purchase_request_detail_system_number"`
	PurchaseRequestSystemNumber       int        `gorm:"column:purchase_request_system_number;size:50" json:"purchase_request_system_number"`
	PurchaseRequestLineNumber         int        `gorm:"column:purchase_request_line_number;size:30" json:"purchase_request_line_number"`
	PurchaseRequestLineStatus         string     `gorm:"column:purchase_request_line_status;size:2;" json:"purchase_request_line_status"`
	ItemCode                          string     `gorm:"column:item_code;size:30;" json:"item_id"`
	ItemQuantity                      *float64   `gorm:"column:item_quantity;" json:"item_quantity"`
	ItemUnitOfMeasure                 string     `gorm:"column:item_unit_of_measure;size:5;" json:"item_unit_of_measures"`
	ItemPrice                         *float64   `gorm:"column:item_price;" json:"item_price"`
	ItemTotal                         *float64   `gorm:"column:item_total;" json:"item_total"`
	ItemRemark                        string     `gorm:"column:item_remark;size:256;" json:"item_remark"`
	PurchaseOrderSystemNumber         int        `gorm:"column:purchase_order_system_number;size:30;" json:"purchase_order_system_number"`
	PurchaseOrderLine                 int        `gorm:"column:purchase_order_line;size:30;" json:"purchase_order_line"`
	ReferenceTypeId                   string     `gorm:"column:reference_type_id;size:10;" json:"reference_type_id"`
	ReferenceSystemNumber             int        `gorm:"column:reference_system_number;size:30;" json:"reference_system_number"`
	ReferenceLine                     int        `gorm:"column:reference_line;size:30;" json:"reference_line"`
	VehicleId                         int        `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
	CreatedByUserId                   int        `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate                       *time.Time `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId                   int        `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate                       *time.Time `gorm:"column:updated_date" json:"updated_date"`
}

func (*PurchaseRequestDetail) TableName() string {
	return TableNamePurchaseRequestDetail
}
