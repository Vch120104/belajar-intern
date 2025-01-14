package transactionworkshopentities

import "time"

const TableNameAtpmClaimVehicle = "trx_atpm_claim_vehicle"

type AtpmClaimVehicle struct {
	ClaimSystemNumber       int       `gorm:"column:claim_system_number;size:30;primaryKey" json:"claim_system_number"`
	CompanyId               int       `gorm:"column:company_id;size:30" json:"company_id"`
	ClaimNumber             string    `gorm:"column:claim_number;size:30" json:"claim_number"`
	ClaimDate               time.Time `gorm:"column:claim_date" json:"claim_date"`
	ClaimStatusId           int       `gorm:"column:claim_status_id;size:30" json:"claim_status_id"`
	BrandId                 int       `gorm:"column:brand_id;size:30" json:"brand_id"`
	ClaimFrom               time.Time `gorm:"column:claim_from" json:"claim_from"`
	ClaimTo                 time.Time `gorm:"column:claim_to" json:"claim_to"`
	ClaimTypeId             int       `gorm:"column:claim_type_id;size:30" json:"claim_type_id"`
	Countermeasure          string    `gorm:"column:countermeasure" json:"countermeasure"`
	WorkOrderSystemNumber   int       `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `gorm:"column:work_order_document_number;size:30" json:"work_order_document_number"`
	WorkOrderDate           time.Time `gorm:"column:work_order_date" json:"work_order_date"`
	VehicleId               int       `gorm:"column:vehicle_id;size:30" json:"vehicle_id"`
	ModelId                 int       `gorm:"column:model_id;size:30" json:"model_id"`
	VariantId               int       `gorm:"column:variant_id;size:30" json:"variant_id"`
	BpkBookNumber           int       `gorm:"column:bpk_book_number;size:30" json:"bpk_book_number"`
	BpkDocumentNumber       string    `gorm:"column:bpk_document_number;size:30" json:"bpk_document_number"`
	TotalAfterDiscount      float64   `gorm:"column:total_after_discount" json:"total_after_discount"`
	FspAmountClaimStandard  float64   `gorm:"column:fsp_amount_claim_standard" json:"fsp_amount_claim_standard"`
	ShippingCostPartRequest float64   `gorm:"column:shipping_cost_part_request" json:"shipping_cost_part_request"`
	ClaimHeader             string    `gorm:"column:claim_header" json:"claim_header"`
	SymtomCode              string    `gorm:"column:symtom_code" json:"symtom_code"`
	TroubleCode             string    `gorm:"column:trouble_code" json:"trouble_code"`
	Pfp                     string    `gorm:"column:pfp" json:"pfp"`
	AfsArea                 string    `gorm:"column:afs_area" json:"afs_area"`
	LabourSellingPrice      float64   `gorm:"column:labour_selling_price" json:"labour_selling_price"`
	TotalFrtQty             float64   `gorm:"column:total_frt_qty" json:"total_frt_qty"`
	TotalLabour             float64   `gorm:"column:total_labour" json:"total_labour"`
	TotalPart               float64   `gorm:"column:total_part" json:"total_part"`
	CustomerComplaint       string    `gorm:"column:customer_complaint" json:"customer_complaint"`
	TechnicianDiagnostic    string    `gorm:"column:technician_diagnostic" json:"technician_diagnostic"`
	Fuel                    float64   `gorm:"column:fuel" json:"fuel"`
	CustomerId              int       `gorm:"column:customer_id;size:30" json:"customer_id"`
	Vdn                     string    `gorm:"column:vdn" json:"vdn"`
	DealerCodeKia           string    `gorm:"column:dealer_code_kia" json:"dealer_code_kia"`
	RepairEndDate           time.Time `gorm:"column:repair_end_date" json:"repair_end_date"`
}

func (*AtpmClaimVehicle) TableName() string {
	return TableNameAtpmClaimVehicle
}
