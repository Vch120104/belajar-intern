package transactionsparepartentities

import (
	masterentities "after-sales/api/entities/master"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const tableNameStockTransaction = "trx_stock_transaction"

type StockTransaction struct {
	IsActive                  bool `gorm:"column:is_active;not null"        json:"is_active"`
	CompanyId                 int  `gorm:"column:company_id;not null;size:30"        json:"company_id"`
	StockTransactionId        int  `gorm:"column:stock_transaction_id;not null;primaryKey;size:30"        json:"stock_transaction_id"`
	TransactionLine           int  `gorm:"column:transaction_line;null;size:30"        json:"transaction_line"`
	TransactionTypeId         int  `gorm:"column:transaction_type_id;not null;size:30"        json:"transaction_type_id"`
	TransactionType           *masterentities.StockTransactionType
	TransactionReasonId       int `gorm:"column:transaction_reason_id;null;size:30"        json:"transaction_reason_id"`
	TransactionReason         *masterentities.StockTransactionReason
	ReferenceId               int       `gorm:"column:reference_id;not null;size:30"        json:"reference_id"`
	ReferenceDocumentNumber   string    `gorm:"column:reference_document_number;null;size:25"        json:"reference_document_number"`
	ReferenceDate             time.Time `gorm:"column:reference_date;null"        json:"reference_date"`
	ReferenceWarehouseId      int       `gorm:"column:reference_warehouse_id;null;size:30"        json:"reference_warehouse_id"`
	ReferenceWarehouse        *masterwarehouseentities.WarehouseMaster
	ReferenceWarehouseGroupId int `gorm:"column:reference_warehouse_group_id;null;size:30"        json:"reference_warehouse_group_id"`
	ReferenceWarehouseGroup   *masterwarehouseentities.WarehouseGroup
	ReferenceLocationId       int `gorm:"column:reference_location_id;null;size:30"        json:"reference_location_id"`
	//ReferenceLocation         masterwarehouseentities.WarehouseLocation
	ReferenceItemId int `gorm:"column:reference_item_id;null;size:30"        json:"reference_item_id"`
	//ReferenceItem                masteritementities.Item
	ReferenceQuantity            float64   `gorm:"column:reference_quantity;null"        json:"reference_quantity"`
	ReferenceUnitOfMeasurementId int       `gorm:"column:reference_unit_of_measurement_id;null;size:30"        json:"reference_unit_of_measurement_id"`
	ReferencePrice               float64   `gorm:"column:reference_price;null"        json:"reference_price"`
	ReferenceCurrencyId          int       `gorm:"column:reference_currency_id;null;size:30"        json:"reference_currency_id"`
	TransactionCogs              float64   `gorm:"column:transaction_cogs;null"        json:"transaction_cogs"`
	ChangeNo                     int       `gorm:"column:change_no;not null;size:30"        json:"change_no"`
	CreatedByUserId              int       `gorm:"column:created_by_user_id;size:30"        json:"created_by_user_id"`
	CreatedDate                  time.Time `gorm:"column:created_date;"        json:"created_date"`
	UpdatedByUserId              int       `gorm:"column:updated_by_user_id;size:30"        json:"updated_by_user_id"`
	UpdatedDate                  time.Time `gorm:"column:updated_date;"        json:"updated_date"`
	VehicleId                    int       `gorm:"column:vehicle_id;null;size:30"        json:"vehicle_id"`
	ItemClassId                  int       `gorm:"column:item_class_id;null;size:30"        json:"item_class_id"`
	//ItemClass                    *masteritementities.ItemClass
}

func (*StockTransaction) TableName() string {
	return tableNameStockTransaction
}
