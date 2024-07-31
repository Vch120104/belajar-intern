package transactionsparepartpayloads

import "time"

type PurchaseRequestResponses struct {
	PurchaseRequestSystemNumber     int       `json:"purchase_request_system_number" parent_entity:"trx_purchase_request" main_table:"trx_purchase_request"`
	PurchaseRequestDocumentNumber   string    `json:"purchase_request_no" parent_entity:"trx_purchase_request" main_table:"trx_purchase_request"`
	PurchaseRequestDocumentDate     time.Time `json:"purchase_request_document_date" parent_entity:"trx_purchase_request" `
	ItemGroupId                     int       `json:"item_group_id" parent_entity:"mtr_work_order_status"`
	OrderTypeId                     int       `json:"order_type_id" parent_entity:"trx_purchase_request"`
	ReferenceDocumentNumber         string    `json:"reference_document_number" parent_entity:"trx_purchase_request"`
	ExpectedArrivalDate             time.Time `json:"expected_arrival_date" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentStatusId int       `json:"purchase_request_document_status_id" parent_entity:"trx_purchase_request"`
	CreatedByUserId                 int       `json:"created_by_user_id" parent_entity:"trx_purchase_request"`
	//BillingCustomer         int       `json:"billing_customer" gorm:"column:billable_to_id"`
}

type PurchaseRequestGetAllListResponses struct {
	PurchaseRequestDocumentNumber string    `json:"purchase_request_no" parent_entity:"trx_purchase_request" main_table:"trx_work_order"`
	PurchaseRequestDocumentDate   time.Time `json:"purchase_request_date" parent_entity:"trx_purchase_request" `
	ItemGroup                     string    `json:"item_group" parent_entity:"mtr_work_order_status"`
	OrderType                     string    `json:"order_type" parent_entity:"trx_purchase_request"`
	ReferenceNo                   string    `json:"reference_no" parent_entity:"trx_purchase_request"`
	ExpectedArrivalDate           time.Time `json:"expected_arrival_date" parent_entity:"trx_purchase_request"`
	Status                        string    `json:"status" parent_entity:"trx_purchase_request"`
	RequestBy                     string    `json:"request_by" parent_entity:"trx_purchase_request"`

	//BillingCustomer         int       `json:"billing_customer" gorm:"column:billable_to_id"`
}

type PurchaseRequestStatusResponse struct {
	PurchaseRequestStatuId           int    `json:"document_status_id"`
	PurchaseRequestStatusCode        string `json:"document_status_code"`
	PurchaseRequestStatusDescription string `json:"document_status_description"`
}
type DivisionResponse struct {
	DivisionId   int    `json:"division_id"`
	DivisionCode string `json:"division_code"`
	DivisionName string `json:"division_name"`
}

type CostCenterResponses struct {
	CostCenterId   int    `json:"cost_center_id"`
	CostCenterCode string `json:"cost_center_code"`
	CostCenterName string `json:"cost_center_name"`
}
type ProfitCenterResponses struct {
	ProfitCenterId   int    `json:"profit_center_id"`
	ProfitCenterCode string `json:"profit_center_code"`
	ProfitCenterName string `json:"profit_center_name"`
}

type WarehouseGroupResponses struct {
	WarehouseGroupId   int    `json:"warehouse_group_id"`
	WarehouseGroupCode string `json:"warehouse_group_code"`
	WarehouseGroupName string `json:"warehouse_group_name"`
}
type WarehouseResponses struct {
	WarehouseId   int    `json:"warehouse_id"`
	WarehouseCode string `json:"warehouse_code"`
	WarehouseName string `json:"warehouse_name"`
}

type CurrencyCodeResponse struct {
	CurrencyId   int    `json:"currency_id"`
	CurrencyCode string `json:"currency_code"`
	CurrencyName string `json:"currency_name"`
}
type PurchaseRequestBrandName struct {
	BrandId   int    `json:"brand_id"`
	BrandName string `json:"brand_name"`
	BrandCode string `json:"brand_code"`
}
type PurchaseRequestRequestedByResponse struct {
	UserEmployeeId   int    `json:"user_employee_id"`
	UserEmployeeName string `json:"employee_name"`
	UserId           int    `json:"user_id"`
}
type WorkOrderDocNoResponses struct {
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
}

type PurchaseRequestReferenceResponses struct {
	ReferenceTypeId   int    `json:"reference_type_id" parent_entity:"mtr_reference_purchase_request"`
	ReferenceTypeCode string `json:"reference_type_code"  parent_entity:"mtr_reference_purchase_request"`
	ReferenceTypeName string `json:"reference_type_name" parent_entity:"mtr_reference_type_purchase_request"`
}
type PurchaseRequestItemGroupResponse struct {
	ItemGroupId   int    `json:"item_group_id"`
	ItemGroupCode string `json:"item_group_code"`
	ItemGroupName string `json:"item_group_name"`
}
type PurchaseRequestItemResponse struct {
	ItemId   int    `json:"item_id"`
	ItemCode string `json:"item_code"`
	ItemName string `json:"item_name"`
}
type UomItemResponses struct {
	SourceConvertion *float64 `json:"source_convertion"`
	TargetConvertion *float64 `json:"target_convertion"`
	SourceUomId      *float64 `json:"source_uom_id"`
	TargetUomId      *float64 `json:"target_uom_id"`
}
type PurchaseRequestOrderTypeResponse struct {
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeCode string `json:"order_type_code"`
	OrderTypeName string `json:"order_type_name"`
}

type PurchaseRequestCompanyResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
}
type PurchaseRequestGetByIdResponses struct {
	CompanyId                       int       `json:"company_id" parent_entity:"trx_purchase_request"`
	PurchaseRequestSystemNumber     int       `json:"purchase_request_system_number" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentNumber   string    `json:"purchase_request_document_number" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentDate     time.Time `json:"purchase_request_document_date" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentStatusId int       `json:"purchase_request_document_status_id" parent_entity:"trx_purchase_request"`
	ItemGroupId                     int       `json:"item_group_id" parent_entity:"trx_purchase_request"`
	BrandId                         int       `json:"brand_id" parent_entity:"trx_purchase_request"`
	ReferenceTypeId                 int       `json:"reference_type_id" parent_entity:"trx_purchase_request"`
	ReferenceSystemNumber           int       `json:"reference_system_number" parent_entity:"trx_purchase_request"`
	ReferenceDocumentNumber         string    `json:"reference_document_number" parent_entity:"trx_purchase_request"`
	//ReferenceWorkOrderSystemNumber   int       `json:"reference_work_order_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceInvoiceUnitSystemNumber int       `json:"reference_invoice_unit_system_number" parent_entity:"trx_purchase_request"`
	//ReferencePickingListSystemNumber int       `json:"reference_picking_list_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceSalesOrderSystemNumber  int       `json:"reference_sales_order_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceSuggorSystemNumber      int       `json:"reference_suggor_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceAutoKfSystemNumber      int       `json:"reference_auto_kf_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceSupplySlipSystemNumber  int       `json:"reference_supply_slip_system_number" parent_entity:"trx_purchase_request"`
	OrderTypeId                int       `json:"order_type_id" parent_entity:"trx_purchase_request"`
	BudgetCode                 string    `json:"budget_code" parent_entity:"trx_purchase_request"`
	ProjectNo                  string    `json:"project_no" parent_entity:"trx_purchase_request"`
	DivisionId                 int       `json:"division_id" parent_entity:"trx_purchase_request"`
	PurchaseRequestRemark      string    `json:"purchase_request_remark" parent_entity:"trx_purchase_request"`
	PurchaseRequestTotalAmount *float64  `json:"purchase_request_total_amount" parent_entity:"trx_purchase_request"`
	ExpectedArrivalDate        time.Time `json:"expected_arrival_date" parent_entity:"trx_purchase_request"`
	ExpectedArrivalTime        time.Time `json:"expected_arrival_time" parent_entity:"trx_purchase_request"`
	CostCenterId               int       `json:"cost_center_id" parent_entity:"trx_purchase_request"`
	ProfitCenterId             int       `json:"profit_center_id" parent_entity:"trx_purchase_request"`
	WarehouseGroupId           int       `json:"warehouse_group_id" parent_entity:"trx_purchase_request"`
	WarehouseId                int       `json:"warehouse_id" parent_entity:"trx_purchase_request"`
	BackOrder                  bool      `parent_entity:"trx_purchase_request" json:"back_order"`
	SetOrder                   bool      `json:"set_order" parent_entity:"trx_purchase_request"`
	CurrencyId                 int       `json:"currency_id" parent_entity:"trx_purchase_request"`
	ItemClassId                int       `json:"column:item_class_id;size:30;" parent_entity:"trx_purchase_request"`
	ChangeNo                   int       `json:"change_no" parent_entity:"trx_purchase_request"`
	CreatedByUserId            int       `json:"created_by_user_id" parent_entity:"trx_purchase_request"`
	CreatedDate                time.Time `json:"created_date" parent_entity:"trx_purchase_request"`
	UpdatedByUserId            int       `json:"updated_by_user_id" parent_entity:"trx_purchase_request"`
	UpdatedDate                time.Time `json:"updated_date" parent_entity:"trx_purchase_request"`
}

type PurchaseRequestGetByIdNormalizeResponses struct {
	Company                       string    `json:"company_id" parent_entity:"trx_purchase_request"`
	PurchaseRequestSystemNumber   int       `json:"purchase_request_system_number" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentNumber string    `json:"purchase_request_document_number" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentDate   time.Time `json:"purchase_request_document_date" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentStatus string    `json:"purchase_request_document_status_id" parent_entity:"trx_purchase_request"`
	ItemGroup                     string    `json:"item_group_id" parent_entity:"trx_purchase_request"`
	Brand                         string    `json:"brand_id" parent_entity:"trx_purchase_request"`
	ReferenceType                 string    `json:"reference_type" parent_entity:"trx_purchase_request"`
	ReferenceSystemNumber         int       `json:"reference_system_number" parent_entity:"trx_purchase_request"`
	ReferenceDocumentNumber       string    `json:"reference_document_number" parent_entity:"trx_purchase_request"`
	//ReferenceInvoiceUnitDocNo     string    `json:"reference_invoice_unit_system_number_doc_no" parent_entity:"trx_purchase_request"`
	//ReferencePickingListDocNo     string    `json:"reference_picking_list_system_number_doc_no" parent_entity:"trx_purchase_request"`
	//ReferenceSalesOrderDocNo      string    `json:"reference_sales_order_system_number_doc_no" parent_entity:"trx_purchase_request"`
	//ReferenceSuggorDocNo          string    `json:"reference_suggor_system_number_doc_no" parent_entity:"trx_purchase_request"`
	//ReferenceAutoKfDocNo          string    `json:"reference_auto_kf_system_number_doc_no" parent_entity:"trx_purchase_request"`
	//ReferenceSupplySlipDocNo      string    `json:"reference_supply_slip_system_number_doc_no" parent_entity:"trx_purchase_request"`
	OrderType                  string    `json:"order_type" parent_entity:"trx_purchase_request"`
	BudgetCode                 string    `json:"budget_code" parent_entity:"trx_purchase_request"`
	ProjectNo                  string    `json:"project_no" parent_entity:"trx_purchase_request"`
	Division                   string    `json:"division_id" parent_entity:"trx_purchase_request"`
	PurchaseRequestRemark      string    `json:"purchase_request_remark" parent_entity:"trx_purchase_request"`
	PurchaseRequestTotalAmount *float64  `json:"purchase_request_total_amount" parent_entity:"trx_purchase_request"`
	ExpectedArrivalDate        time.Time `json:"expected_arrival_date" parent_entity:"trx_purchase_request"`
	ExpectedArrivalTime        time.Time `json:"expected_arrival_time" parent_entity:"trx_purchase_request"`
	CostCenter                 string    `json:"cost_center" parent_entity:"trx_purchase_request"`
	ProfitCenter               string    `json:"profit_center" parent_entity:"trx_purchase_request"`
	WarehouseGroup             string    `json:"warehouse_group_" parent_entity:"trx_purchase_request"`
	Warehouse                  string    `json:"warehouse_" parent_entity:"trx_purchase_request"`
	BackOrder                  bool      `parent_entity:"trx_purchase_request" json:"back_order"`
	SetOrder                   bool      `json:"set_order" parent_entity:"trx_purchase_request"`
	Currency                   string    `json:"currency" parent_entity:"trx_purchase_request"`
	ChangeNo                   int       `json:"change_no" parent_entity:"trx_purchase_request"`
	CreatedByUser              string    `json:"created_by_user_id" parent_entity:"trx_purchase_request"`
	CreatedDate                time.Time `json:"created_date" parent_entity:"trx_purchase_request"`
	UpdatedByUser              string    `json:"updated_by_user_id" parent_entity:"trx_purchase_request"`
	UpdatedDate                time.Time `json:"updated_date" parent_entity:"trx_purchase_request"`
}

type PurchaseRequestDetailRequestPayloads struct {
	ItemCode          string   `json:"item_code" parent_entity:"trx_purchase_request_detail"`
	ItemQuantity      *float64 `json:"item_quantity" parent_entity:"trx_purchase_request_detail"`
	ItemUnitOfMeasure string   `json:"item_unit_of_measures" parent_entity:"trx_purchase_request_detail"`
	ItemRemark        string   `json:"item_remark" parent_entity:"trx_purchase_request_detail"`
}
type PurchaseRequestDetailResponsesPayloads struct {
	ItemCode              string   `json:"item_code"`
	ItemName              string   `json:"item_name"`
	ItemQuantity          *float64 `json:"item_quantity"`
	ItemUnitOfMeasure     string   `json:"item_unit_of_measures"`
	ItemUnitOfMeasureRate float64  `json:"item_unit_of_measure_rate"`
	ItemRemark            string   `json:"item_remark"`
}

type PurchaseRequestHeaderSaveRequest struct {
	CompanyId                       int       `json:"company_id" parent_entity:"trx_purchase_request"`
	PurchaseRequestSystemNumber     int       `json:"purchase_request_system_number" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentNumber   string    `json:"purchase_request_document_number" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentDate     time.Time `json:"purchase_request_document_date" parent_entity:"trx_purchase_request"`
	PurchaseRequestDocumentStatusId int       `json:"purchase_request_document_status_id" parent_entity:"trx_purchase_request"`
	ItemGroupId                     int       `json:"item_group_id" parent_entity:"trx_purchase_request"`
	BrandId                         int       `json:"brand_id" parent_entity:"trx_purchase_request"`
	ReferenceTypeId                 int       `json:"reference_type_id" parent_entity:"trx_purchase_request"`
	ReferenceSystemNumber           int       `json:"reference_system_number" parent_entity:"trx_purchase_request"`
	ReferenceDocumentNumber         string    `json:"reference_document_number" parent_entity:"trx_purchase_request"`
	//ReferenceWorkOrderSystemNumber   int       `json:"reference_work_order_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceInvoiceUnitSystemNumber int       `json:"reference_invoice_unit_system_number" parent_entity:"trx_purchase_request"`
	//ReferencePickingListSystemNumber int       `json:"reference_picking_list_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceSalesOrderSystemNumber  int       `json:"reference_sales_order_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceSuggorSystemNumber      int       `json:"reference_suggor_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceAutoKfSystemNumber      int       `json:"reference_auto_kf_system_number" parent_entity:"trx_purchase_request"`
	//ReferenceSupplySlipSystemNumber  int       `json:"reference_supply_slip_system_number" parent_entity:"trx_purchase_request"`
	OrderTypeId                int       `json:"order_type_id" parent_entity:"trx_purchase_request"`
	BudgetCode                 string    `json:"budget_code" parent_entity:"trx_purchase_request"`
	ProjectNo                  string    `json:"project_no" parent_entity:"trx_purchase_request"`
	DivisionId                 int       `json:"division_id" parent_entity:"trx_purchase_request"`
	PurchaseRequestRemark      string    `json:"purchase_request_remark" parent_entity:"trx_purchase_request"`
	PurchaseRequestTotalAmount *float64  `json:"purchase_request_total_amount" parent_entity:"trx_purchase_request"`
	ExpectedArrivalDate        time.Time `json:"expected_arrival_date" parent_entity:"trx_purchase_request"`
	ExpectedArrivalTime        time.Time `json:"expected_arrival_time" parent_entity:"trx_purchase_request"`
	CostCenterId               int       `json:"cost_center_id" parent_entity:"trx_purchase_request"`
	ProfitCenterId             int       `json:"profit_center_id" parent_entity:"trx_purchase_request"`
	WarehouseGroupId           int       `json:"warehouse_group_id" parent_entity:"trx_purchase_request"`
	WarehouseId                int       `json:"warehouse_id" parent_entity:"trx_purchase_request"`
	BackOrder                  bool      `parent_entity:"trx_purchase_request" json:"back_order"`
	SetOrder                   bool      `json:"set_order" parent_entity:"trx_purchase_request"`
	CurrencyId                 int       `json:"currency_id" parent_entity:"trx_purchase_request"`
	ItemClassId                int       `json:"column:item_class_id;size:30;" parent_entity:"trx_purchase_request"`
	ChangeNo                   int       `json:"change_no" parent_entity:"trx_purchase_request"`
	CreatedByUserId            int       `json:"created_by_user_id" parent_entity:"trx_purchase_request"`
	CreatedDate                time.Time `json:"created_date" parent_entity:"trx_purchase_request"`
	UpdatedByUserId            int       `json:"updated_by_user_id" parent_entity:"trx_purchase_request"`
	UpdatedDate                time.Time `json:"updated_date" parent_entity:"trx_purchase_request"`
}
