package transactionworkshopentities

import "time"

var CreateWorkOrderHistoryDetailTable = "trx_work_order_history_detail"

type WorkOrderHistoryDetail struct {
	WorkOrderHistoryDetailId            int       `gorm:"column:work_order_history_detail_system_number;size:30;not null;primaryKey" json:"work_order_history_detail_system_number"`
	WorkOrderSystemNumber               int       `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	WorkOrderOperationItemLine          int       `gorm:"column:work_order_operation_item_line;size:30" json:"work_order_operation_item_line"`
	WorkOrderLineStatusID               int       `gorm:"column:work_order_line_status_id;size:30" json:"work_order_line_status_id"`
	LineTypeID                          int       `gorm:"column:line_type_id;size:30" json:"line_type_id"`
	ServiceStatusCode                   int       `gorm:"column:service_status_code;size:30" json:"service_status_code"`
	WorkOrderTransactionTypeID          string    `gorm:"column:work_order_transaction_type_id;size:100" json:"work_order_transaction_type_id"`
	JobTypeID                           int       `gorm:"column:job_type_id;size:30" json:"job_type_id"`
	WorkOrderLineDiscountStatusID       int       `gorm:"column:work_order_line_discount_status_id;size:30" json:"work_order_line_discount_status_id"`
	OperationItemID                     int       `gorm:"column:operation_item_id;size:30" json:"operation_item_id"`
	Description                         string    `gorm:"column:description;size:100" json:"description"`
	ItemUOMID                           int       `gorm:"column:item_uom_id;size:30" json:"item_uom_id"`
	FRTQuantity                         float32   `gorm:"column:frt_quantity" json:"frt_quantity"`
	OperationItemPrice                  float32   `gorm:"column:operation_item_price" json:"operation_item_price"`
	OperationItemDiscountAmount         float32   `gorm:"column:operation_item_discount_amount" json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount  float32   `gorm:"column:operation_item_discount_request_amount" json:"operation_item_discount_request_amount"`
	OperationItemDiscountPercent        float32   `gorm:"column:operation_item_discount_percent" json:"operation_item_discount_percent"`
	OperationItemDiscountRequestPercent float32   `gorm:"column:operation_item_discount_request_percent" json:"operation_item_discount_request_percent"`
	PackageID                           string    `gorm:"column:package_id;size:100" json:"package_id"`
	TotalCOGS                           float32   `gorm:"column:total_cogs" json:"total_cogs"`
	PPHAmount                           float32   `gorm:"column:pph_amount" json:"pph_amount"`
	TaxID                               string    `gorm:"column:tax_id;size:100" json:"tax_id"`
	PPHTaxRate                          float32   `gorm:"column:pph_tax_rate" json:"pph_tax_rate"`
	LastApprovalBy                      string    `gorm:"column:last_approval_by;size:100" json:"last_approval_by"`
	LastApprovalDate                    time.Time `gorm:"column:last_approval_date;type:datetime" json:"last_approval_date"`
	QualityControlStatusID              int       `gorm:"column:quality_control_status_id;size:30" json:"quality_control_status_id"`
	QualityControlExtraFRT              float32   `gorm:"column:quality_control_extra_frt" json:"quality_control_extra_frt"`
	QualityControlExtraReason           string    `gorm:"column:quality_control_extra_reason;size:100" json:"quality_control_extra_reason"`
	SupplyQuantity                      float32   `gorm:"column:supply_quantity" json:"supply_quantity"`
	SubstituteID                        int       `gorm:"column:substitute_id;size:30" json:"substitute_id"`
	SubstituteItemID                    string    `gorm:"column:substitute_item_id;size:100" json:"substitute_item_id"`
	WarehouseGroupID                    int       `gorm:"column:warehouse_group_id;size:30" json:"warehouse_group_id"`
	ATPMWarrantyClaimFormDocumentNumber string    `gorm:"column:atpm_warranty_claim_form_document_number;size:100" json:"atpm_warranty_claim_form_document_number"`
	ATPMWarrantyClaimFormDate           time.Time `gorm:"column:atpm_warranty_claim_form_date;type:datetime" json:"atpm_warranty_claim_form_date"`
	ATPMWarrantyClaimFormTypeID         string    `gorm:"column:atpm_warranty_claim_form_type_id;size:100" json:"atpm_warranty_claim_form_type_id"`
	PurchaseRequestSystemNumber         int       `gorm:"column:purchase_request_system_number;size:30" json:"purchase_request_system_number"`
	PurchaseRequestLineID               int       `gorm:"column:purchase_request_line_id;size:30" json:"purchase_request_line_id"`
	PurchaseOrderSystemNumber           int       `gorm:"column:purchase_order_system_number;size:30" json:"purchase_order_system_number"`
	PurchaseOrderLineID                 int       `gorm:"column:purchase_order_line_id;size:30" json:"purchase_order_line_id"`
	InvoiceSystemNumber                 int       `gorm:"column:invoice_system_number;size:30" json:"invoice_system_number"`
	GRPOQuantity                        float32   `gorm:"column:grpo_quantity" json:"grpo_quantity"`
	QualityControlTotalExtraFRT         float32   `gorm:"column:quality_control_total_extra_frt" json:"quality_control_total_extra_frt"`
	ReorderNumber                       float32   `gorm:"column:reorder_number" json:"reorder_number"`
	BinningQuantity                     float32   `gorm:"column:binning_quantity" json:"binning_quantity"`
	IncentiveSystemNumber               float32   `gorm:"column:incentive_system_number" json:"incentive_system_number"`
	Bypass                              bool      `gorm:"column:bypass" json:"bypass"`
	TechnicianID                        int       `gorm:"column:technician_id;size:30" json:"technician_id"`
}

func (WorkOrderHistoryDetail) TableName() string {
	return CreateWorkOrderHistoryDetailTable
}
