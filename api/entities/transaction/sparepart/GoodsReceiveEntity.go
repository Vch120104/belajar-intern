package transactionsparepartentities

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const GoodReceiveTableName = "trx_goods_receive"

type GoodsReceive struct {
	CompanyId                  int       `gorm:"column:company_id;not null;size:30" json:"company_id"`
	GoodsReceiveSystemNumber   int       `gorm:"column:goods_receive_system_number;size:30;not null;primaryKey"        json:"goods_receive_system_number"`
	GoodsReceiveDocumentNumber string    `gorm:"column:goods_receive_document_number;not null;size:25"        json:"goods_receive_document_number"`
	GoodsReceiveDocumentDate   time.Time `gorm:"column:goods_receive_document_date;not null"        json:"goods_receive_document_date"`
	GoodsReceiveStatusId       int       `gorm:"column:goods_receive_status_id;size:30;not null"        json:"goods_receive_status_id"`
	ReferenceTypeGoodReceiveId int       `gorm:"column:reference_type_good_receive_id;size:30;not null"        json:"reference_type_good_receive_id"`
	//ReferenceTypeGoodReceive    *masterentities.GoodsReceiveReferenceType // `gorm:"foreignKey:reference_type_good_receive_id;references:reference_type_good_receive_id"`
	ReferenceSystemNumber   int    `gorm:"column:reference_system_number;size:30;null"        json:"reference_system_number"`
	ReferenceDocumentNumber string `gorm:"column:reference_document_number;not null;size:25"        json:"reference_document_number"`

	AffiliatedPurchaseOrder bool `gorm:"column:affiliated_purchase_order;not null"        json:"affiliated_purchase_order"`
	ViaBinning              bool `gorm:"column:via_binning;not null"        json:"via_binning"`
	SetOrder                bool `gorm:"column:set_order;not null"        json:"set_order"`
	BackOrder               bool `gorm:"column:back_order;not null"        json:"back_order"`
	BrandId                 int  `gorm:"column:brand_id;not null"        json:"brand_id"`
	CostCenterId            int  `gorm:"column:cost_center_id;not null"        json:"cost_center_id"`
	ProfitCenterId          int  `gorm:"column:profit_center_id;not null"        json:"profit_center_id"`

	TransactionTypeId int `gorm:"column:transaction_type_id;not null"        json:"transaction_type_id"`

	EventId    int `gorm:"column:event_id;not null"        json:"event_id"`
	SupplierId int `gorm:"column:supplier_id;not null"        json:"supplier_id"`

	SupplierDeliveryOrderNumber string `gorm:"column:supplier_delivery_order_number;null;size:25"        json:"supplier_delivery_order_number"`

	SupplierInvoiceNumber      string    `gorm:"column:supplier_invoice_number;null;size:25"        json:"supplier_invoice_number"`
	SupplierInvoiceDate        time.Time `gorm:"column:supplier_invoice_date;null"        json:"supplier_invoice_date"`
	SupplierTaxInvoiceNumber   string    `gorm:"column:supplier_tax_invoice_number;null"        json:"supplier_tax_invoice_number"`
	SupplierTaxInvoiceDate     time.Time `gorm:"column:supplier_tax_invoice_date;null"        json:"supplier_tax_invoice_date"`
	WarehouseGroupId           int       `gorm:"column:warehouse_group_id;not null"        json:"warehouse_group_id"`
	WarehouseId                int       `gorm:"column:warehouse_id;not null;size:30;"        json:"warehouse_id"`
	Warehouse                  *masterwarehouseentities.WarehouseMaster
	WarehouseClaimId           int `gorm:"column:warehouse_claim_id;not null;size:30;"        json:"warehouse_claim_id"`
	ItemGroupId                int `gorm:"column:item_group_id;not null;size:30;"        json:"item_group_id"`
	ItemGroup                  *masteritementities.ItemGroup
	CurrencyId                 int       `gorm:"column:currency_id;not null;size:30;"        json:"currency_id"`
	CurrencyExchangeRateDate   time.Time `gorm:"column:currency_exchange_rate_date; null"        json:"currency_exchange_rate_date"`
	CurrencyExchangeRate       float64   `gorm:"column:currency_exchange_rate;not null"        json:"currency_exchange_rate"`
	CurrencyExchangeRateTypeId int       `gorm:"column:currency_exchange_rate_type_id;not null"        json:"currency_exchange_rate_type_id"`
	PrintingNumber             int       `gorm:"column:printing_number;null"        json:"printing_number"`
	LastPrintedById            int       `gorm:"column:last_printed_by_id;null"        json:"last_printed_by_id"`
	JournalSystemNumber        int       `gorm:"column:journal_system_number;not null"        json:"journal_system_number"`
	ItemClassId                int       `gorm:"column:item_class_id;not null;size:30;"        json:"item_class_id"`
	itemClass                  *masteritementities.ItemClass
	UseInTransitWarehouse      bool                 `gorm:"column:use_in_transit_warehouse;null"        json:"use_in_transit_warehouse"`
	InTransitWarehouseId       int                  `gorm:"column:in_transit_warehouse_id;null"        json:"in_transit_warehouse_id"`
	FreightCost                float64              `gorm:"column:freight_cost;null"        json:"freight_cost"`
	InsuranceCost              float64              `gorm:"column:insurance_cost;null"        json:"insurance_cost"`
	OthersCost                 float64              `gorm:"column:others_cost;null"        json:"others_cost"`
	Rounding                   float64              `gorm:"column:rounding;null"        json:"rounding"`
	ActualPartArrivalDate      time.Time            `gorm:"column:actual_part_arrival_date;null"        json:"actual_part_arrival_date"`
	GoodsReceiveDetail         []GoodsReceiveDetail `gorm:"foreignKey:GoodsReceiveSystemNumber;references:GoodsReceiveSystemNumber" json:"binning_stock_detail"`
	CreatedByUserId            int                  `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate                time.Time            `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId            int                  `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate                time.Time            `gorm:"column:updated_date" json:"updated_date"`
	ChangeNo                   int                  `gorm:"column:change_no;size:30;" json:"change_no"`
}

func (*GoodsReceive) TableName() string {
	return GoodReceiveTableName
}
