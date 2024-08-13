package transactionsparepartpayloads

import (
	"time"
)

// type SupplySlipResponse struct {
// 	IsActive              bool      `json:"is_active"`
// 	SupplySystemNumber    int32     `json:"supply_system_number"`
// 	SupplyDocumentNumber  string    `json:"supply_document_number"`
// 	SupplyStatusId        int32     `json:"supply_status_id"`
// 	SupplyDate            time.Time `json:"supply_date"`
// 	SupplyTypeId          int32     `json:"supply_type_id"`
// 	CompanyId             int32     `json:"company_id"`
// 	WorkOrderSystemNumber int32     `json:"work_order_system_number"`
// 	ProfitCenterId        int32     `json:"profit_center_id"`
// 	BrandId               int32     `json:"brand_id"`
// 	ModelId               int32     `json:"model_id"`
// 	VariantId             int32     `json:"variant_id"`
// 	VehicleId             int32     `json:"vehicle_id"`
// 	CustomerId            int32     `json:"customer_id"`
// 	TechnicianId          int32     `json:"technician_id "`
// 	CampaignId            int32     `json:"campaign_id"`
// 	Remark                string    `json:"remark"`
// }

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

type SupplySlipSearchResponse struct {
	SupplySystemNumber      int    `json:"supply_system_number" parent_entity:"trx_supply_slip" main_table:"trx_supply_slip"`
	SupplyDocumentNumber    string `json:"supply_document_number" parent_entity:"trx_supply_slip"`
	SupplyDate              string `json:"supply_date" parent_entity:"trx_supply_slip"`
	SupplyTypeId            int    `json:"supply_type_id" parent_entity:"trx_supply_slip"` //external - general common
	WorkOrderSystemNumber   int    `json:"work_order_system_number" parent_entity:"trx_work_order" references:"trx_work_order"`
	WorkOrderDocumentNumber string `json:"work_order_document_number" parent_entity:"trx_work_order"`
	CustomerId              int    `json:"customer_id" parent_entity:"trx_work_order" `      //external - general (attribute from trx_work_order)
	SupplyStatusId          int    `json:"supply_status_id" parent_entity:"trx_supply_slip"` //external - general
}

type SupplySlipResponse struct {
	SupplySystemNumber      int    `json:"supply_system_number" parent_entity:"trx_supply_slip" main_table:"trx_supply_slip"`
	SupplyStatusId          int    `json:"supply_status_id" parent_entity:"trx_supply_slip"` //external - general
	SupplyDocumentNumber    string `json:"supply_document_number" parent_entity:"trx_supply_slip"`
	SupplyDate              string `json:"supply_date" parent_entity:"trx_supply_slip"`
	SupplyTypeId            int    `json:"supply_type_id" parent_entity:"trx_supply_slip"` //external - general common
	WorkOrderSystemNumber   int    `json:"work_order_system_number" parent_entity:"trx_work_order" references:"trx_work_order"`
	WorkOrderDocumentNumber string `json:"work_order_document_number" parent_entity:"trx_work_order"`
	WorkOrderDate           string `json:"work_order_date" parent_entity:"trx_work_order"`
	CustomerId              int    `json:"customer_id" parent_entity:"trx_work_order" `   //external - general (attribute from trx_work_order)
	VehicleId               int    `json:"vehicle_id" parent_entity:"trx_work_order" `    //external - sales (attribute from trx_work_order)
	BrandId                 int    `json:"brand_id" parent_entity:"trx_work_order" `      //external - sales (attribute from trx_work_order)
	ModelId                 int    `json:"model_id" parent_entity:"trx_work_order" `      //external - sales (attribute from trx_work_order)
	VariantId               int    `json:"variant_id" parent_entity:"trx_work_order" `    //external - sales (attribute from trx_work_order)
	TechnicianId            int    `json:"technician_id" parent_entity:"trx_supply_slip"` //external - general
	CampaignId              int    `json:"campaign_id" parent_entity:"mtr_campaign" references:"mtr_campaign"`
	CampaignCode            string `json:"campaign_code" parent_entity:"mtr_campaign"`
}

type SupplyTypeResponse struct {
	SupplyTypeId          int    `json:"supply_type_id"`
	SupplyTypeDescription string `json:"supply_type_description"`
}

type ApprovalStatusResponse struct {
	SupplyStatusId          int    `json:"approval_status_id"`
	SupplyStatusDescription string `json:"approval_status_description"`
}

type CustomerResponse struct {
	CustomerId   int    `json:"customer_id"`
	CustomerCode string `json:"customer_code"`
	CustomerName string `json:"customer_name"`
}

type TechnicianResponse struct {
	TechnicianId   int    `json:"user_employee_id"`
	TechnicianName string `json:"employee_name"`
}

// type VehicleResponse struct {
// 	TechnicianId   int    `json:"user_employee_id"`
// 	TechnicianName string `json:"employee_name"`
// }

type BrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
}

type ModelResponse struct {
	ModelId          int    `json:"model_id"`
	ModelCode        string `json:"model_code"`
	ModelDescription string `json:"model_description"`
}

type VariantResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantDescription string `json:"variant_description"`
	ProductionYear     string `json:"production_year"`
}
