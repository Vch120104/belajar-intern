package transactionsparepartentities

import "time"

var CreateSupplySlipTable = "trx_supply_slip"

type SupplySlip struct {
	IsActive              bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SupplySystemNumber    int32     `gorm:"column:supply_system_number;not null;primaryKey"        json:"supply_system_number"`
	SupplyDocumentNumber  string    `gorm:"column:supply_document_number;null"        json:"supply_document_number"`
	SupplyStatusId        int32     `gorm:"column:supply_status_id;null"        json:"supply_status_id"`
	SupplyDate            time.Time `gorm:"column:supply_date;null"        json:"supply_date"`
	SupplyTypeId          int32     `gorm:"column:supply_type_id;not null"        json:"supply_type_id"`
	CompanyId             int32     `gorm:"column:company_id;null"        json:"company_id"`
	WorkOrderSystemNumber int32     `gorm:"column:work_order_system_number;null"        json:"work_order_system_number"`
	ProfitCenterId        int32     `gorm:"column:profit_center_id;null"        json:"profit_center_id"`
	BrandId               int32     `gorm:"column:brand_id;null"        json:"brand_id"`
	ModelId               int32     `gorm:"column:model_id;null"        json:"model_id"`
	VariantId             int32     `gorm:"column:variant_id;null"        json:"variant_id"`
	VehicleId             int32     `gorm:"column:vehicle_id;null"        json:"vehicle_id"`
	CustomerId            int32     `gorm:"column:customer_id;null"        json:"customer_id"`
	TechnicianId          int32     `gorm:"column:technician_id ;null"        json:"technician_id "`
	CampaignId            int32     `gorm:"column:campaign_id;null"        json:"campaign_id"`
	Remark                string    `gorm:"column:remark;null"        json:"remark"`
}

func (*SupplySlip) TableName() string {
	return CreateSupplySlipTable
}
