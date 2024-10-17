package transactionsparepartentities

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const TableNameItemClaim = "trx_item_claim"

type ItemClaim struct {
	CompanyId                  int       `gorm:"column:company_id;null;size:30" json:"company_id"`
	ClaimSystemNumber          int       `gorm:"column:claim_system_number;not null;primaryKey;size:30" json:"claim_system_number"`
	ClaimStatus                string    `gorm:"column:claim_status;null;size:2" json:"claim_status"`
	ClaimDocumentNumber        string    `gorm:"column:claim_document_number;null;size:25" json:"claim_document_number"`
	ClaimDate                  time.Time `gorm:"column:claim_date;null" json:"claim_date"`
	ClaimTypeId                int       `gorm:"column:claim_type_id;null;size:30" json:"claim_type_id"`
	ClaimType                  *masterentities.ItemClaimType
	GoodsReceiveSystemNumber   int    `gorm:"column:goods_receive_system_number;null;size:30" json:"goods_receive_system_number"`
	GoodsReceiveDocumentNumber string `gorm:"column:goods_receive_document_number;null;size:25" json:"goods_receive_document_number"`
	VehicleBrandId             int    `gorm:"column:vehicle_brand_id;null;size:30" json:"vehicle_brand_id"`
	CostCenterId               int    `gorm:"column:cost_center_id;null;size:30" json:"cost_center_id"`
	ProfitCenterId             int    `gorm:"column:profit_center_id;null;size:30" json:"profit_center_id"`
	TransactionTypeId          int    `gorm:"column:transaction_type_id;null;size:30" json:"transaction_type_id"`
	EventId                    int    `gorm:"column:event_id;null;size:30" json:"event_id"`
	SupplierId                 int    `gorm:"column:supplier_id;null;size:30" json:"supplier_id"`
	SuppplierDoNumber          string `gorm:"column:suppplier_do_number;null;size:25" json:"suppplier_do_number"`
	ReferenceTypeGoodReceiveId int    `gorm:"column:reference_type_good_receive_id;null;size:30" json:"reference_type_good_receive_id"`
	ReferenceSystemNumber      int    `gorm:"column:reference_system_number;not null;size:30" json:"reference_system_number"`
	ReferenceDocumentNumber    string `gorm:"column:reference_document_number;null;size:25" json:"reference_document_number"`
	WarehouseGroupId           int    `gorm:"column:warehouse_group_id;null;size:30" json:"warehouse_group_id"`
	WarehouseGroup             *masterwarehouseentities.WarehouseGroup
	//WarehouseId                int `gorm:"column:warehouse_id;null;size:30" json:"warehouse_id"`
	//Warehouse                  masterwarehouseentities.WarehouseMaster
	ItemGroupId              int `gorm:"column:item_group_id;null;size:30" json:"item_group_id"`
	ItemGroup                *masteritementities.ItemGroup
	ViaBinning               bool      `gorm:"column:via_binning;null" json:"via_binning"`
	CurrencyId               int       `gorm:"column:currency_id;null;size:30" json:"currency_id"`
	CurrencyExchangeRateDate time.Time `gorm:"column:currency_exchange_rate_date;null" json:"currency_exchange_rate_date"`
	CurrencyExchangeRate     float64   `gorm:"column:currency_exchange_rate;null" json:"currency_exchange_rate"`
	CurrencyRateType         string    `gorm:"column:currency_rate_type;null;size:25" json:"currency_rate_type"`
	PrintingNumber           int       `gorm:"column:printing_number;null;size:30" json:"printing_number"`
	LastPrintedBy            string    `gorm:"column:last_printed_by;null;size:10" json:"last_printed_by"`
	JournalSystemNumber      int       `gorm:"column:journal_system_number;not null;size:30" json:"journal_system_number"`
}
