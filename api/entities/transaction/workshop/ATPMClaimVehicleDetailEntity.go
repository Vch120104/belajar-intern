package transactionworkshopentities

const TableNameAtpmClaimVehicleDetail = "trx_atpm_claim_vehicle_detail"

type AtpmClaimVehicleDetail struct {
	ClaimDetailSystemNumber int     `gorm:"column:claim_detail_system_number;size:30;primaryKey" json:"claim_detail_system_number"`
	ClaimSystemNumber       int     `gorm:"column:claim_system_number;size:30" json:"claim_system_number"`
	CompanyId               int     `gorm:"column:company_id;size:30" json:"company_id"`
	ClaimLineNumber         int     `gorm:"column:claim_line_number;size:30" json:"claim_line_number"`
	WorkOrderSystemNumber   int     `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	WorkOrderLineNumber     int     `gorm:"column:work_order_line_number;size:30" json:"work_order_line_number"`
	LineTypeId              int     `gorm:"column:line_type_id;size:30" json:"line_type_id"`
	ItemId                  int     `gorm:"column:item_id;size:30" json:"item_id"`
	FrtQuantity             float64 `gorm:"column:frt_quantity" json:"frt_quantity"`
	ItemPrice               float64 `gorm:"column:item_price" json:"item_price"`
	DiscountPercent         float64 `gorm:"column:discount_percent" json:"discount_percent"`
	DiscountAmount          float64 `gorm:"column:discount_amount" json:"discount_amount"`
	TotalAfterDiscount      float64 `gorm:"column:total_after_discount" json:"total_after_discount"`
	PartRequest             int     `gorm:"column:part_request;size:30" json:"part_request"`
	RecallNumber            int     `gorm:"column:recall_number;size:30" json:"recall_number"`
	IncidentPartReceived    int     `gorm:"column:incident_part_received;size:30" json:"incident_part_received"`
}

func (*AtpmClaimVehicleDetail) TableName() string {
	return TableNameAtpmClaimVehicleDetail
}
