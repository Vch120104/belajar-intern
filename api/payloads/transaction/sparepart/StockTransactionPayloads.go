package transactionsparepartpayloads

import (
	"time"
)

type StockTransactionInsertPayloads struct {
	CompanyId                    int        `gorm:"column:company_id;not null;size:30"        json:"company_id" :"company_id"`
	StockTransactionId           int        `gorm:"column:stock_transaction_id;not null;primaryKey;size:30"        json:"stock_transaction_id" :"stock_transaction_id"`
	TransactionTypeId            int        `gorm:"column:transaction_type_id;not null;size:30"        json:"transaction_type_id" :"transaction_type_id"`
	TransactionReasonId          int        `gorm:"column:transaction_reason_id;null;size:30"        json:"transaction_reason_id" :"transaction_reason_id"`
	ReferenceId                  int        `gorm:"column:reference_id;not null;size:30"        json:"reference_id" :"reference_id"`
	ReferenceDocumentNumber      string     `gorm:"column:reference_document_number;null;size:25"        json:"reference_document_number" :"reference_document_number"`
	ReferenceDate                *time.Time `gorm:"column:reference_date;null"        json:"reference_date" :"reference_date"`
	ReferenceWarehouseId         int        `gorm:"column:reference_warehouse_id;null;size:30"        json:"reference_warehouse_id" :"reference_warehouse_id"`
	ReferenceWarehouseGroupId    int        `gorm:"column:reference_warehouse_group_id;null;size:30"        json:"reference_warehouse_group_id" :"reference_warehouse_group_id"`
	ReferenceLocationId          int        `gorm:"column:reference_location_id;null;size:30"        json:"reference_location_id" :"reference_location_id"`
	ReferenceItemId              int        `gorm:"column:reference_item_id;null;size:30"        json:"reference_item_id" :"reference_item_id"`
	ReferenceQuantity            float64    `gorm:"column:reference_quantity;null"        json:"reference_quantity" :"reference_quantity"`
	ReferenceUnitOfMeasurementId int        `gorm:"column:reference_unit_of_measurement_id;null;size:30"        json:"reference_unit_of_measurement_id" :"reference_unit_of_measurement_id"`
	ReferencePrice               float64    `gorm:"column:reference_price;null"        json:"reference_price" :"reference_price"`
	ReferenceCurrencyId          int        `gorm:"column:reference_currency_id;null;size:30"        json:"reference_currency_id" :"reference_currency_id"`
	TransactionCogs              float64    `gorm:"column:transaction_cogs;null"        json:"transaction_cogs" :"transaction_cogs"`
	ChangeNo                     int        `gorm:"column:change_no;not null;size:30"        json:"change_no" :"change_no"`
	CreatedByUserId              int        `gorm:"column:created_by_user_id;size:30"        json:"created_by_user_id" :"created_by_user_id"`
	CreatedDate                  time.Time  `gorm:"column:created_date;"        json:"created_date" :"created_date"`
	UpdatedByUserId              int        `gorm:"column:updated_by_user_id;size:30"        json:"updated_by_user_id" :"updated_by_user_id"`
	UpdatedDate                  time.Time  `gorm:"column:updated_date;"        json:"updated_date" :"updated_date"`
	VehicleId                    int        `gorm:"column:vehicle_id;null;size:30"        json:"vehicle_id" :"vehicle_id"`
	ItemClassId                  int        `gorm:"column:item_class_id;null;size:30"        json:"item_class_id" :"item_class_id"`
}

type LocationUpdateResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
