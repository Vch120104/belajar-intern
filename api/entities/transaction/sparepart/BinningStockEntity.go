package transactionsparepartentities

import (
	masterentities "after-sales/api/entities/master"
	"time"
)

const TableNameBinningStock = "trx_binning_list_stock"

type BinningStock struct {
	CompanyId               int        `gorm:"column:company_id;not null"  json:"company_id"`
	BinningSystemNumber     int        `gorm:"column:binning_system_number;not null;primaryKey;size:30"  json:"binning_system_number"`
	BinningDocumentStatusId int        `gorm:"column:binning_document_status_id;not null"  json:"binning_document_status_id"`
	BinningDocumentNumber   string     `gorm:"column:binning_document_number;not null"  json:"binning_document_number"`
	BinningDocumentDate     *time.Time `gorm:"column:binning_document_date;not null"  json:"binning_document_date"`
	BinningReferenceTypeId  int        `gorm:"column:binning_reference_type_id;not null"  json:"binning_reference_type_id"`
	BinningReferenceType    *masterentities.BinningReferenceTypeMaster
	ReferenceSystemNumber   int    `gorm:"column:reference_system_number;not null"  json:"reference_system_number"`
	ReferenceDocumentNumber string `gorm:"column:reference_document_number;not null"  json:"reference_document_number"`
	WarehouseGroupId        int    `gorm:"column:warehouse_group_id;not null;size:30"  json:"warehouse_group_id"`
	//WarehouseGroup          *masterwarehouseentities.WarehouseGroup //`gorm:"foreignKey:WarehouseGroupId;references:WarehouseGroupId"`
	WarehouseId int `gorm:"column:warehouse_id;not null;size:30"  json:"warehouse_id"`
	//Warehouse                   *masterwarehouseentities.WarehouseMaster //`gorm:"foreignKey:WarehouseId;references:WarehouseId"`
	SupplierId                  int        `gorm:"column:supplier_id;not null"  json:"supplier_id"`
	SupplierDeliveryOrderNumber string     `gorm:"column:supplier_delivery_order_number;not null"  json:"supplier_delivery_order_number"`
	SupplierInvoiceNumber       string     `gorm:"column:supplier_invoice_number;null"  json:"supplier_invoice_number"`
	SupplierInvoiceDate         *time.Time `gorm:"column:supplier_invoice_date;null"  json:"supplier_invoice_date"`
	SupplierFakturPajakNumber   string     `gorm:"column:supplier_faktur_pajak_number;null"  json:"supplier_faktur_pajak_number"`
	SupplierFakturPajakDate     *time.Time `gorm:"column:supplier_faktur_pajak_date;null"  json:"supplier_faktur_pajak_date"`
	SupplierDeliveryPerson      string     `gorm:"column:supplier_delivery_person;not null"  json:"supplier_delivery_person"`
	SupplierCaseNumber          string     `gorm:"column:supplier_case_number;not null"  json:"supplier_case_number"`
	BinningTypeId               int        `gorm:"column:binning_type_id ;not null"  json:"binning_type_id "`
	BinningType                 *masterentities.BinningTypeMaster
	CurrencyId                  int `gorm:"column:currency_id ;not null"  json:"currency_id "`
	ExchangeId                  int `gorm:"column:exchange_id ;not null"  json:"exchange_id "`
	//BinningStockDetail          []BinningStockDetail                     `gorm:"foreignKey:BinningSystemNumber;references:BinningSystemNumber" json:"binning_stock_detail"`
	BinningStockDetail []BinningStockDetail `gorm:"foreignKey:BinningSystemNumber;references:BinningSystemNumber" json:"binning_stock_detail"`
	CreatedByUserId    int                  `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate        *time.Time           `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId    int                  `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate        *time.Time           `gorm:"column:updated_date" json:"updated_date"`
	ChangeNo           int                  `gorm:"column:change_no;size:30;" json:"change_no"`
}

func (*BinningStock) TableName() string {
	return TableNameBinningStock
}
