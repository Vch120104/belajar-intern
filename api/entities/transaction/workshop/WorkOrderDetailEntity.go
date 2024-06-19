package transactionworkshopentities

import "time"

const TableNameWorkOrderDetail = "trx_work_order_detail"

type WorkOrderDetail struct {
	WorkOrderDetailId                   int       `gorm:"column:work_order_detail_id;size:30;primaryKey" json:"work_order_detail_id"`
	WorkOrderSystemNumber               int       `gorm:"column:work_order_system_number;size:30;" json:"work_order_system_number"`
	LineType                            string    `gorm:"column:line_type;size:50;" json:"line_type"`
	OperationId                         int       `gorm:"column:operation_id;size:30;" json:"operation_id"`
	ItemId                              int       `gorm:"column:item_id;size:30;" json:"item_id"`
	OperationItemCode                   string    `gorm:"column:operation_item_code;size:50;" json:"operation_item_code"`
	WorkOrderOperationItemLine          int       `gorm:"column:work_order_operation_item_line;size:30;" json:"work_order_operation_item_line"`
	WorkorderStatusId                   int       `gorm:"column:workorder_status_id;size:30;" json:"workorder_status_id"`
	LineTypeId                          int       `gorm:"column:line_type_id;size:30;" json:"line_type_id"`
	ServiceStatusId                     int       `gorm:"column:service_status_id;size:30;" json:"service_status_id"`
	TransactionTypeId                   int       `gorm:"column:transaction_type_id;size:30;" json:"transaction_type_id"`
	JobTypeId                           int       `gorm:"column:job_type_id;size:30;" json:"job_type_id"`
	ApprovalId                          int       `gorm:"column:approval_id;size:30;" json:"approval_id"`
	Description                         string    `gorm:"column:description;size:50;" json:"description"`
	FrtQuantity                         float32   `gorm:"column:frt_quantity" json:"frt_quantity"`
	OperationItemPrice                  float32   `gorm:"column:operation_item_price" json:"operation_item_price"`
	OperationItemDiscountAmount         float32   `gorm:"column:operation_item_discount_amount" json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount  float32   `gorm:"column:operation_item_discount_request_amount" json:"operation_item_discount_request_amount"`
	OperationItemDiscountPercent        float32   `gorm:"column:operation_item_discount_percent" json:"operation_item_discount_percent"`
	OperationItemDiscountRequestPercent float32   `gorm:"column:operation_item_discount_request_percent" json:"operation_item_discount_request_percent"`
	PackageId                           int       `gorm:"column:package_id;size:30;" json:"package_id"`
	TotalCostOfGoodsSold                float32   `gorm:"column:total_cost_of_goods_sold" json:"total_cost_of_goods_sold"`
	PphAmount                           float32   `gorm:"column:pph_amount" json:"pph_amount"`
	TaxId                               int       `gorm:"column:tax_id;size:30;" json:"tax_id"`
	PphTaxRate                          float32   `gorm:"column:pph_tax_rate" json:"pph_tax_rate"`
	LastApprovalBy                      string    `gorm:"column:last_approval_by;size:50;" json:"last_approval_by"`
	LastApprovalDate                    time.Time `gorm:"column:last_approval_date;" json:"last_approval_date"`
	QualityControlStatus                string    `gorm:"column:quality_control_status;size:50;" json:"quality_control_status"`
	QualityControlExtraFrt              float32   `gorm:"column:quality_control_extra_frt" json:"quality_control_extra_frt"`
	QualityControlExtraReason           string    `gorm:"column:quality_control_extra_reason;size:50;" json:"quality_control_extra_reason"`
	SupplyQuantity                      float32   `gorm:"column:supply_quantity" json:"supply_quantity"`
	SubstituteId                        int       `gorm:"column:substitute_id;size:30;" json:"substitute_id"`
	SubstrituteItemCode                 string    `gorm:"column:substritute_item_code;size:50;" json:"substritute_item_code"`
	WarehouseId                         int       `gorm:"column:warehouse_id;size:30;" json:"warehouse_id"`
	AtpmClaimNumber                     string    `gorm:"column:atpm_claim_number;size:50;" json:"atpm_claim_number"`
	AtpmClaimDate                       time.Time `gorm:"column:atpm_claim_date;" json:"atpm_claim_date"`
	WarrantyClaimTypeId                 int       `gorm:"column:warranty_claim_type_id;size:30;" json:"warranty_claim_type_id"`
	PurchaseRequestSystemNumber         int       `gorm:"column:purchase_request_system_number;size:30;" json:"purchase_request_system_number"`
	PurchaseRequestDetailId             int       `gorm:"column:purchase_request_detail_id;size:30;" json:"purchase_request_detail_id"`
	PurchaseOrderSystemNumber           int       `gorm:"column:purchase_order_system_number;size:30;" json:"purchase_order_system_number"`
	PurchaseOrderLine                   int       `gorm:"column:purchase_order_line;size:30;" json:"purchase_order_line"`
	InvoiceSystemNumber                 int       `gorm:"column:invoice_system_number;size:30;" json:"invoice_system_number"`
	GoodsReceiveQuantity                float32   `gorm:"column:goods_receive_quantity" json:"goods_receive_quantity"`
	QualityControlTotalExtraFrt         float32   `gorm:"column:quality_control_total_extra_frt" json:"quality_control_total_extra_frt"`
	ReorderNumber                       float32   `gorm:"column:reorder_number" json:"reorder_number"`
	BinningQuantity                     float32   `gorm:"column:binning_quantity" json:"binning_quantity"`
	IncentiveSystemNumber               int       `gorm:"column:incentive_system_number;size:30;" json:"incentive_system_number"`
	Bypass                              bool      `gorm:"column:bypass" json:"bypass"`
	TechnicianId                        int       `gorm:"column:technician_id;size:30;" json:"technician_id"`
	UserEmployeeId                      int       `gorm:"column:user_employee_id;size:30;" json:"user_employee_id"`
	FieldActionSystemNumber             int       `gorm:"column:field_action_system_number;size:30;" json:"field_action_system_number"`
	Request                             string    `gorm:"column:request;size:50;" json:"request"`
	FrtQuantityExpress                  float32   `gorm:"column:frt_quantity_express" json:"frt_quantity_express"`
	PriceListId                         int       `gorm:"column:price_list_id;size:30;" json:"price_list_id"`
	ServiceCategoryId                   int       `gorm:"column:service_category_id;size:30;" json:"service_category_id"`
	ClaimSystemNumber                   int       `gorm:"column:claim_system_number;size:30;" json:"claim_system_number"`
	QualityControlPassDatetime          time.Time `gorm:"column:quality_control_pass_datetime;" json:"quality_control_pass_datetime"`
	ExtendedWarranty                    bool      `gorm:"column:extended_warranty" json:"extended_warranty"`
	RemarkExtendedWarranty              string    `gorm:"column:remark_extended_warranty;size:50;" json:"remark_extended_warranty"`
}

func (*WorkOrderDetail) TableName() string {
	return TableNameWorkOrderDetail
}
