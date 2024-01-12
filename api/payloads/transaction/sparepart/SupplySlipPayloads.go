package transactionsparepartpayloads

import "time"

type SupplySlipResponse struct {
	IsActive              bool      `json:"is_active"`
	SupplySystemNumber    int32     `json:"supply_system_number"`
	SupplyDocumentNumber  string    `json:"supply_document_number"`
	SupplyStatusId        int32     `json:"supply_status_id"`
	SupplyDate            time.Time `json:"supply_date"`
	SupplyTypeId          int32     `json:"supply_type_id"`
	CompanyId             int32     `json:"company_id"`
	WorkOrderSystemNumber int32     `json:"work_order_system_number"`
	ProfitCenterId        int32     `json:"profit_center_id"`
	BrandId               int32     `json:"brand_id"`
	ModelId               int32     `json:"model_id"`
	VariantId             int32     `json:"variant_id"`
	VehicleId             int32     `json:"vehicle_id"`
	CustomerId            int32     `json:"customer_id"`
	TechnicianId          int32     `json:"technician_id "`
	CampaignId            int32     `json:"campaign_id"`
	Remark                string    `json:"remark"`
}

type SupplySlipRequest struct {
	SupplyDocumentNumber  string    `json:"supply_document_number"`
	SupplyStatusId        int32     `json:"supply_status_id"`
	SupplyDate            time.Time `json:"supply_date"`
	SupplyTypeId          int32     `json:"supply_type_id"`
	CompanyId             int32     `json:"company_id"`
	WorkOrderSystemNumber int32     `json:"work_order_system_number"`
	ProfitCenterId        int32     `json:"profit_center_id"`
	BrandId               int32     `json:"brand_id"`
	ModelId               int32     `json:"model_id"`
	VariantId             int32     `json:"variant_id"`
	VehicleId             int32     `json:"vehicle_id"`
	CustomerId            int32     `json:"customer_id"`
	TechnicianId          int32     `json:"technician_id "`
	CampaignId            int32     `json:"campaign_id"`
	Remark                string    `json:"remark"`
}
