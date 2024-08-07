package transactionsparepartentities

import "time"

const TableNamePurchaseRequest = "trx_purchase_request"

type PurchaseRequestEntities struct {
	CompanyId                   int `gorm:"column:company_id;size:30;" json:"company_id"`
	PurchaseRequestSystemNumber int `gorm:"column:purchase_request_system_number;size:50;not null;primaryKey;" json:"purchase_request_system_number"`

	//PurchaseRequestSystemNumber     int        `gorm:"column:purchase_request_system_number;size:50;not null;primaryKey;" json:"purchase_request_system_number"`
	PurchaseRequestDocumentNumber   string     `gorm:"column:purchase_request_document_number;size:50;" json:"purchase_request_document_number"`
	PurchaseRequestDocumentDate     *time.Time `gorm:"column:purchase_request_document_date" json:"purchase_request_document_date"`
	PurchaseRequestDocumentStatusId int        `gorm:"column:purchase_request_document_status_id;size:30;;not null" json:"purchase_request_document_status_id"`
	ItemGroupId                     int        `gorm:"column:item_group_id;size:30;" json:"item_group_id"`
	BrandId                         int        `gorm:"column:brand_id;size:30;" json:"brand_id"`
	ReferenceTypeId                 int        `gorm:"column:reference_type_id;size:30;" json:"reference_type_id"`
	ReferenceSystemNumber           int        `gorm:"column:reference_system_number;size:30;" json:"reference_system_number"`
	ReferenceDocumentNumber         string     `gorm:"column:reference_document_number;size:50;" json:"reference_document_number"`

	//ReferenceWorkOrderSystemNumber   int                     `gorm:"column:reference_work_order_system_number;size:30;" json:"reference_work_order_system_number"`
	//ReferenceInvoiceUnitSystemNumber int                     `gorm:"column:reference_invoice_unit_system_number;size:30;" json:"reference_invoice_unit_system_number"`
	//ReferencePickingListSystemNumber int                     `gorm:"column:reference_picking_list_system_number;size:30;" json:"reference_picking_list_system_number"`
	//ReferenceSalesOrderSystemNumber  int                     `gorm:"column:reference_sales_order_system_number;size:30;" json:"reference_sales_order_system_number"`
	//ReferenceSuggorSystemNumber      int                     `gorm:"column:reference_suggor_system_number;size:30;" json:"reference_suggor_system_number"`
	//ReferenceAutoKfSystemNumber      int                     `gorm:"column:reference_auto_kf_system_number;size:30;" json:"reference_auto_kf_system_number"`
	//ReferenceSupplySlipSystemNumber  int                     `gorm:"column:reference_supply_slip_system_number;size:30;" json:"reference_supply_slip_system_number"`
	OrderTypeId                int        `gorm:"column:order_type_id;size:30;" json:"order_type_id"`
	BudgetCode                 string     `gorm:"column:budget_code;size:50;" json:"budget_code"`
	ProjectNo                  string     `gorm:"column:project_no;size:50;" json:"project_no"`
	DivisionId                 int        `gorm:"column:division_id;size:30;" json:"division_id"`
	PurchaseRequestRemark      string     `gorm:"column:purchase_request_remark;size:256;" json:"purchase_request_remark"`
	PurchaseRequestTotalAmount *float64   `gorm:"column:purchase_request_total_amount;" json:"purchase_request_total_amount"`
	ExpectedArrivalDate        *time.Time `gorm:"column:expected_arrival_date" json:"expected_arrival_date"`
	ExpectedArrivalTime        *time.Time `gorm:"column:expected_arrival_time;size:5;" json:"expected_arrival_time"`
	CostCenterId               int        `gorm:"column:cost_center_id;size:30;" json:"cost_center_id"`
	ProfitCenterId             int        `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
	WarehouseGroupId           int        `gorm:"column:warehouse_group_id;size:30;" json:"warehouse_group_id"`
	WarehouseId                int        `gorm:"column:warehouse_id;size:30;" json:"warehouse_id"`
	BackOrder                  bool       `gorm:"column:back_order;" json:"back_order"`
	SetOrder                   bool       `gorm:"set_order;" json:"set_order"`
	CurrencyId                 int        `gorm:"column:currency_id;size:30;" json:"currency_id"`
	//ReferenceSystemNumber      int                     `gorm:"column:reference_system_number;size:30;" json:"reference_system_number"`
	ItemClassId           int                     `gorm:"column:item_class_id;size:30;" json:"item_class_id"`
	PurchaseRequestDetail []PurchaseRequestDetail `gorm:"foreignKey:PurchaseRequestSystemNumber;references:PurchaseRequestSystemNumber" json:"purchase_request_detail"`
	CreatedByUserId       int                     `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate           *time.Time              `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId       int                     `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate           *time.Time              `gorm:"column:updated_date" json:"updated_date"`
	ChangeNo              int                     `gorm:"column:change_no;size:30;" json:"change_no"`
}

func (*PurchaseRequestEntities) TableName() string {
	return TableNamePurchaseRequest
}
