package transactionsparepartentities

import "time"

const TableNameSupplySlip = "trx_supply_slip"

type SupplySlip struct {
	IsActive              bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SupplySystemNumber    int        `gorm:"column:supply_system_number;size:30;;not null;primaryKey;size:30" json:"supply_system_number"`
	SupplyDocumentNumber  string     `gorm:"column:supply_document_number;size:50;" json:"supply_document_number"`
	SupplyStatusId        int        `gorm:"column:supply_status_id;size:30;" json:"supply_status_id"`
	SupplyDate            *time.Time `gorm:"column:supply_date" json:"supply_date"`
	SupplyTypeId          int        `gorm:"column:supply_type_id;size:30;;not null" json:"supply_type_id"`
	CompanyId             int        `gorm:"column:company_id;size:30;" json:"company_id"`
	WorkOrderSystemNumber int        `gorm:"column:work_order_system_number;size:30;" json:"work_order_system_number"`
	ProfitCenterId        int        `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
	BrandId               int        `gorm:"column:brand_id;size:30;" json:"brand_id"`
	ModelId               int        `gorm:"column:model_id;size:30;" json:"model_id"`
	VariantId             int        `gorm:"column:variant_id;size:30;" json:"variant_id"`
	VehicleId             int        `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
	CustomerId            int        `gorm:"column:customer_id;size:30;" json:"customer_id"`
	TechnicianId          int        `gorm:"column:technician_id;size:30;" json:"technician_id"`
	CampaignId            int        `gorm:"column:campaign_id;size:30;" json:"campaign_id"`
	Remark                string     `gorm:"column:remark;size:50;" json:"remark"`
}

func (*SupplySlip) TableName() string {
	return TableNameSupplySlip
}
