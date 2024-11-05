package transactionsparepartpayloads

import "time"

type GoodsReceivesGetAllPayloads struct {
	GoodsReceiveSystemNumber    int       `json:"goods_receive_system_number"`
	GoodsReceiveDocumentNumber  string    `json:"goods_receive_document_number"`
	GoodsReceiveDocumentDate    time.Time `json:"goods_receive_document_date"`
	ItemGroupName               string    `json:"item_group_name"`
	ReferenceDocumentNumber     string    `json:"reference_document_number"`
	SupplierId                  int       `json:"supplier_id"`
	SupplierName                string    `json:"supplier_name"`
	GoodsReceiveStatusId        int       `json:"goods_receive_status_id"`
	JournalSystemNumber         int       `json:"journal_system_number"`
	SupplierDeliveryOrderNumber string    `json:"supplier_delivery_order_number"`
	QuantityGoodsReceive        float64   `json:"quantity_goods_receive"`
	TotalAmount                 float64   `json:"total_amount"`
}

type GoodsReceivesGetByIdResponses struct {
	GoodsReceiveSystemNumber    int       `json:"goods_receive_system_number"`
	GoodsReceiveDocumentNumber  string    `json:"goods_receive_document_number"`
	ItemGroupId                 int       `json:"item_group_id"`
	GoodsReceiveDocumentDate    time.Time `json:"goods_receive_document_date"`
	SupplierId                  int       `json:"supplier_id"`
	GoodsReceiveStatusId        int       `json:"goods_receive_status_id"`
	ReferenceTypeGoodReceiveId  int       `json:"reference_type_good_receive_id"`
	ReferenceSystemNumber       int       `json:"reference_system_number"`
	ReferenceDocumentNumber     string    `json:"reference_document_number"`
	AffiliatedPurchaseOrder     bool      `json:"affiliated_purchase_order"`
	ViaBinning                  bool      `json:"via_binning"`
	SetOrder                    bool      `json:"set_order"`
	BackOrder                   bool      `json:"back_order"`
	BrandId                     int       `json:"brand_id"`
	CostCenterId                int       `json:"cost_center_id"`
	ProfitCenterId              int       `json:"profit_center_id"`
	TransactionTypeId           int       `json:"transaction_type_id"`
	EventId                     int       `json:"event_id"`
	SupplierDeliveryOrderNumber string    `json:"supplier_delivery_order_number"`
	SupplierInvoiceNumber       string    `json:"supplier_invoice_number"`
	SupplierTaxInvoiceNumber    string    `json:"supplier_tax_invoice_number"`
	WarehouseId                 int       `json:"warehouse_id"`
	WarehouseCode               string    `json:"warehouse_code"`
	WarehouseName               string    `json:"warehouse_name"`
	WarehouseClaimId            int       `json:"warehouse_claim_id"`
	WarehouseClaimCode          string    `json:"warehouse_claim_code"`
	WarehouseClaimName          string    `json:"warehouse_claim_name"`
	ItemClassId                 int       `json:"item_class_id"`
}

type GoodsReceiveInsertPayloads struct {
	GoodsReceiveStatusId        int       `gorm:"column:goods_receive_status_id;size:30;not null"        json:"goods_receive_status_id"`
	CompanyId                   int       `gorm:"column:company_id;not null;size:30" json:"company_id"`
	GoodsReceiveDocumentNumber  string    `gorm:"column:goods_receive_document_number;not null;size:25"        json:"goods_receive_document_number"`
	GoodsReceiveDocumentDate    time.Time `gorm:"column:goods_receive_document_date;not null"        json:"goods_receive_document_date"`
	ReferenceTypeGoodReceiveId  int       `gorm:"column:reference_type_good_receive_id;size:30;not null"        json:"reference_type_good_receive_id"`
	ReferenceSystemNumber       int       `gorm:"column:reference_system_number;size:30;null"        json:"reference_system_number"`
	ReferenceDocumentNumber     string    `gorm:"column:reference_document_number;not null;size:25"        json:"reference_document_number"`
	AffiliatedPurchaseOrder     bool      `gorm:"column:affiliated_purchase_order;not null"        json:"affiliated_purchase_order"`
	ViaBinning                  bool      `gorm:"column:via_binning;not null"        json:"via_binning"`
	SetOrder                    bool      `gorm:"column:set_order;not null"        json:"set_order"`
	BackOrder                   bool      `gorm:"column:back_order;not null"        json:"back_order"`
	BrandId                     int       `gorm:"column:brand_id;not null"        json:"brand_id"`
	CostCenterId                int       `gorm:"column:cost_center_id;not null"        json:"cost_center_id"`
	ProfitCenterId              int       `gorm:"column:profit_center_id;not null"        json:"profit_center_id"`
	TransactionTypeId           int       `gorm:"column:transaction_type_id;not null"        json:"transaction_type_id"`
	EventId                     int       `gorm:"column:event_id;not null"        json:"event_id"`
	SupplierId                  int       `gorm:"column:supplier_id;not null"        json:"supplier_id"`
	SupplierInvoiceNumber       string    `gorm:"column:supplier_invoice_number;null;size:25"        json:"supplier_invoice_number"`
	SupplierInvoiceDate         time.Time `gorm:"column:supplier_invoice_date;null"        json:"supplier_invoice_date"`
	SupplierTaxInvoiceNumber    string    `gorm:"column:supplier_tax_invoice_number;null"        json:"supplier_tax_invoice_number"`
	SupplierTaxInvoiceDate      time.Time `gorm:"column:supplier_tax_invoice_date;null"        json:"supplier_tax_invoice_date"`
	WarehouseGroupId            int       `gorm:"column:warehouse_group_id;not null"        json:"warehouse_group_id"`
	WarehouseId                 int       `gorm:"column:warehouse_id;not null;size:30;"        json:"warehouse_id"`
	WarehouseClaimId            int       `gorm:"column:warehouse_claim_id;not null;size:30;"        json:"warehouse_claim_id"`
	CurrencyId                  int       `gorm:"column:currency_id;not null;size:30;"        json:"currency_id"`
	CurrencyExchangeRate        float64   `gorm:"column:currency_exchange_rate;not null"        json:"currency_exchange_rate"`
	CurrencyExchangeRateTypeId  int       `gorm:"column:currency_exchange_rate_type_id;not null"        json:"currency_exchange_rate_type_id"`
	UseInTransitWarehouse       bool      `gorm:"column:use_in_transit_warehouse;null"        json:"use_in_transit_warehouse"`
	InTransitWarehouseId        int       `gorm:"column:in_transit_warehouse_id;null"        json:"in_transit_warehouse_id"`
	SupplierDeliveryOrderNumber string    `gorm:"column:supplier_delivery_order_number;null;size:25"        json:"supplier_delivery_order_number"`
	ItemGroupId                 int       `gorm:"column:item_group_id;not null;size:30;"        json:"item_group_id"`
	CreatedByUserId             int       `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate                 time.Time `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId             int       `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate                 time.Time `gorm:"column:updated_date" json:"updated_date"`
}
type GoodsReceiveUpdatePayloads struct {
	GoodsReceiveDetailSystemNumber int       `json:"goods_receive_detail_system_number"`
	SupplierDeliveryOrderNumber    string    `gorm:"column:supplier_delivery_order_number" json:"supplier_delivery_order_number"`
	ReferenceSystemNumber          int       `gorm:"column:reference_system_number" json:"reference_system_number"`
	ReferenceTypeGoodReceiveId     int       `gorm:"column:reference_type_good_receive_id" json:"reference_type_good_receive_id"`
	ReferenceDocumentNumber        string    `gorm:"column:reference_document_number" json:"reference_document_number"`
	ProfitCenterId                 int       `gorm:"column:profit_center_id" json:"profit_center_id"`
	TransactionTypeId              int       `gorm:"column:transaction_type_id" json:"transaction_type_id"`
	EventId                        int       `gorm:"column:event_id" json:"event_id"`
	AffiliatedPurchaseOrder        bool      `gorm:"column:affiliated_purchase_order" json:"affiliated_purchase_order"`
	ViaBinning                     bool      `gorm:"column:via_binning" json:"via_binning"`
	SupplierId                     int       `gorm:"column:supplier_id" json:"supplier_id"`
	SupplierInvoiceNumber          string    `gorm:"column:supplier_invoice_number" json:"supplier_invoice_number"`
	SupplierInvoiceDate            time.Time `gorm:"column:supplier_invoice_date" json:"supplier_invoice_date"`
	SupplierTaxInvoiceNumber       string    `gorm:"column:supplier_tax_invoice_number" json:"supplier_tax_invoice_number"`
	SupplierTaxInvoiceDate         time.Time `gorm:"column:supplier_tax_invoice_date" json:"supplier_tax_invoice_date"`
	WarehouseGroupId               int       `gorm:"column:warehouse_group_id" json:"warehouse_group_id"`
	WarehouseId                    int       `gorm:"column:warehouse_id" json:"warehouse_id"`
	WarehouseClaimId               int       `gorm:"column:warehouse_claim_id" json:"warehouse_claim_id"`
	ItemGroupId                    int       `gorm:"column:item_group_id" json:"item_group_id"`
	UpdatedByUserId                int       `gorm:"column:updated_by_user_id" json:"updated_by_user_id"`
	UpdatedDate                    time.Time `gorm:"column:updated_date" json:"updated_date"`
	UseInTransitWarehouse          bool      `gorm:"column:use_in_transit_warehouse" json:"use_in_transit_warehouse"`
	InTransitWarehouseId           int       `gorm:"column:in_transit_warehouse_id" json:"in_transit_warehouse_id"`
}

type GoodsReceiveDetailInsertPayloads struct {
	GoodsReceiveSystemNumber int        `json:"goods_receive_system_number"`
	GoodsReceiveLineNumber   int        `json:"goods_receive_line_number"`
	ItemId                   int        `json:"item_id"`
	ItemUnitOfMeasurement    string     `json:"item_unit_of_measurement"`
	UnitOfMeasurementRate    float64    `json:"unit_of_measurement_rate"`
	ItemPrice                float64    `json:"item_price"`
	ItemDiscountPercent      float64    `json:"item_discount_percent"`
	ItemDiscountAmount       float64    `json:"item_discount_amount"`
	QuantityReference        float64    `json:"quantity_reference"`
	QuantityDeliveryOrder    float64    `json:"quantity_delivery_order"`
	QuantityShort            float64    `json:"quantity_short"`
	QuantityDamage           float64    `json:"quantity_damage"`
	QuantityOver             float64    `json:"quantity_over"`
	QuantityWrong            float64    `json:"quantity_wrong"`
	QuantityGoodsReceive     float64    `json:"quantity_goods_receive"`
	WarehouseLocationId      int        `json:"warehouse_location_id"`
	WarehouseLocationClaimId int        `json:"warehouse_location_claim_id"`
	CaseNumber               string     `json:"case_number"`
	BinningDocumentNumber    string     `json:"binning_document_number"`
	BinningLineNumber        int        `json:"binning_line_number"`
	ReferenceSystemNumber    int        `json:"reference_system_number"`
	ReferenceLineNumber      int        `json:"reference_line_number"`
	ClaimSystemNumber        int        `json:"claim_system_number"`
	ClaimLineNumber          int        `json:"claim_line_number"`
	ItemTotal                float64    `json:"item_total"`
	ItemTotalBaseAmount      float64    `json:"item_total_base_amount"`
	CreatedByUserId          int        `json:"created_by_user_id"`
	CreatedDate              *time.Time `json:"created_date"`
	UpdatedByUserId          int        `json:"updated_by_user_id"`
	UpdatedDate              *time.Time `json:"updated_date"`
	ChangeNo                 int        `json:"change_no"`
	BinningId                int        `gorm:"column:binning_system_number"        json:"binning_system_number"`
}
type GoodsReceiveDetailUpdatePayloads struct {
	GoodsReceiveDetailSystemNumber int        `json:"goods_receive_detail_system_number"`
	GoodsReceiveSystemNumber       int        `json:"goods_receive_system_number"`
	GoodsReceiveLineNumber         int        `json:"goods_receive_line_number"`
	ItemId                         int        `json:"item_id"`
	ItemUnitOfMeasurement          string     `json:"item_unit_of_measurement"`
	UnitOfMeasurementRate          float64    `json:"unit_of_measurement_rate"`
	ItemPrice                      float64    `json:"item_price"`
	ItemDiscountPercent            float64    `json:"item_discount_percent"`
	ItemDiscountAmount             float64    `json:"item_discount_amount"`
	QuantityReference              float64    `json:"quantity_reference"`
	QuantityDeliveryOrder          float64    `json:"quantity_delivery_order"`
	QuantityShort                  float64    `json:"quantity_short"`
	QuantityDamage                 float64    `json:"quantity_damage"`
	QuantityOver                   float64    `json:"quantity_over"`
	QuantityWrong                  float64    `json:"quantity_wrong"`
	QuantityGoodsReceive           float64    `json:"quantity_goods_receive"`
	WarehouseLocationId            int        `json:"warehouse_location_id"`
	WarehouseLocationClaimId       int        `json:"warehouse_location_claim_id"`
	CaseNumber                     string     `json:"case_number"`
	BinningDocumentNumber          string     `json:"binning_document_number"`
	BinningLineNumber              int        `json:"binning_line_number"`
	ReferenceSystemNumber          int        `json:"reference_system_number"`
	ReferenceLineNumber            int        `json:"reference_line_number"`
	ClaimSystemNumber              int        `json:"claim_system_number"`
	ClaimLineNumber                int        `json:"claim_line_number"`
	ItemTotal                      float64    `json:"item_total"`
	ItemTotalBaseAmount            float64    `json:"item_total_base_amount"`
	CreatedByUserId                int        `json:"created_by_user_id"`
	CreatedDate                    *time.Time `json:"created_date"`
	UpdatedByUserId                int        `json:"updated_by_user_id"`
	UpdatedDate                    *time.Time `json:"updated_date"`
	ChangeNo                       int        `json:"change_no"`
	BinningId                      int        `gorm:"column:binning_system_number"        json:"binning_system_number"`
}
type ItemGoodsReceiveTemp struct {
	ItemPrice       float64
	ItemDiscPercent float64
	ItemDiscAmount  float64
}
type GetAllLocationGRPOResponse struct {
	WarehouseId           int    `gorm:"column:warehouse_id" json:"warehouse_id"`
	ItemId                int    `gorm:"column:item_id" json:"item_id"`
	ItemLocationId        int    `gorm:"column:item_location_id" json:"item_location_id"`
	WarehouseLocationName string `gorm:"column:warehouse_location_name" json:"warehouse_location_name"`
	ItemCode              string `gorm:"column:item_code" json:"item_code"`
	CompanyId             int    `gorm:"column:company_id" json:"company_id"`
	WarehouseCode         string `gorm:"column:warehouse_code" json:"warehouse_code"`
}
