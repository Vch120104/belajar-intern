package transactionsparepartentities

type PurchaseOrderDetailEntities struct {
	PurchaseOrderDetailSystemNumber int      `gorm:"column:purchase_order_detail_system_number;size:30;primaryKey" json:"purchase_order_detail_system_number"`
	PurchaseOrderSystemNumber       int      `gorm:"column:purchase_order_system_number;size:30;" json:"purchase_order_system_number"`
	PurchaseOrderLineNumber         int      `gorm:"column:purchase_order_line_number;" json:"purchase_order_line_number"`
	ItemId                          int      `gorm:"column:item_id;" json:"item_id"`
	ItemUnitOfMeasurement           string   `gorm:"column:item_unit_of_measurement;size:1;" json:"item_unit_of_measurement"`
	UnitOfMeasurementRate           *float64 `gorm:"column:unit_of_measurement_rate;" json:"unit_of_measurement_rate"`
	ItemQuantity                    *float64 `gorm:"column:item_quantity;" json:"item_quantity"`
	ItemPrice                       *float64 `gorm:"column:item_price;" json:"item_price"`
	ItemDiscountPercentage          *float64 `gorm:"column:item_discount_percentage;" json:"item_discount_percentage"`
	ItemDiscountAmount              *float64 `gorm:"column:item_discount_amount;" json:"item_discount_amount"`
	ItemTotal                       *float64 `gorm:"column:item_total;" json:"item_total"`
	SubstituteTypeId                int      `gorm:"column:substitute_type_id;" json:"substitute_type_id"`
	PurchaseRequestSystemNumber     int      `gorm:"column:purchase_request_system_number;" json:"purchase_request_system_number"`
	PurchaseRequestLineNumber       int      `gorm:"column:purchase_request_line_number;" json:"purchase_request_line_number"`
	GoodsReceiveQuantity            *float64 `gorm:"column:goods_receive_quantity;" json:"goods_receive_quantity"`
	QuantityInvoiceAccountPayable   *float64 `gorm:"column:quantity_invoice_account_payable;" json:"quantity_invoice_account_payable"`
	OldPurchaseOrderSystemNo        int      `gorm:"column:old_purchase_order_system_no;" json:"old_purchase_order_system_no"`
	OldPurchaseOrderLineNumber      int      `gorm:"column:old_purchase_order_line_number;" json:"old_purchase_order_line_number"`
	BinningQuantity                 *float64 `gorm:"column:binning_quantity;" json:"binning_quantity"`
	VehicleChassisNumber            string   `gorm:"column:vehicle_chassis_number;size:255;" json:"vehicle_chassis_number"`
	StockOnHand                     *float64 `gorm:"column:stock_on_hand;" json:"stock_on_hand"`
	SalesOrderSystemNumber          int      `gorm:"column:sales_order_system_number;" json:"sales_order_system_number"`
	SalesOrderLineNumber            int      `gorm:"column:sales_order_line_number;" json:"sales_order_line_number"`
	ItemRemark                      string   `gorm:"column:item_remark;size:255;" json:"item_remark"`
	GoodsReceiveSystemNumber        int      `gorm:"column:goods_receive_system_number;" json:"goods_receive_system_number"`
	GoodsReceiveLineNumber          int      `gorm:"column:goods_receive_line_number;" json:"goods_receive_line_number"`
	Snp                             *float64 `gorm:"column:snp;" json:"snp"`
}

func (*PurchaseOrderDetailEntities) TableName() string {
	return "trx_item_purchase_order_detail"
}
