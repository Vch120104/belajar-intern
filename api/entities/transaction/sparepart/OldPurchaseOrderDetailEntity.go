package transactionsparepartentities

import masteritementities "after-sales/api/entities/master/item"

const TableNameChangesItemPurchaseOrderDetail = "trx_item_purchase_order_detail_changed_item"

type PurchaseOrderDetailChangedItem struct {
	ChangedItemPurchaseOrderDetailSystemNumber int                          `gorm:"column:changed_item_purchase_order_detail_system_number;size:30;not null;primaryKey" json:"changed_item_purchase_order_detail_system_number"`
	ChangedItemPurchaseOrderSystemNumber       int                          `gorm:"column:changed_item_purchase_order_system_number;not null" json:"changed_item_purchase_order_system_number"`
	ItemPurchaseOrderDetail                    *PurchaseOrderDetailEntities `gorm:"foreignKey:ChangedItemPurchaseOrderDetailSystemNumber;references:PurchaseOrderDetailSystemNumber"`
	ChangedItemPurchaseOrderLineNumber         int                          `gorm:"column:changed_item_purchase_order_line_number;not null" json:"changed_item_purchase_order_line_number"`
	ItemId                                     int                          `gorm:"column:item_id;null" json:"item_id"`
	Item                                       *masteritementities.Item     `gorm:"foreignKey:ItemId;references:ItemId"`
	UomItemId                                  int                          `gorm:"column:uom_item_id;null" json:"uom_item_id"`
	UomItem                                    *masteritementities.UomItem  `gorm:"foreignKey:UomItemId;references:UomItemId"`
	UnitOfMeasureRate                          float64                      `gorm:"column:unit_of_measure_rate;null" json:"unit_of_measure_rate"`
	ItemQuantity                               float64                      `gorm:"column:item_quantity;null" json:"item_quantity"`
	ItemPrice                                  float64                      `gorm:"column:item_price;null" json:"item_price"`
	ItemDiscountPercent                        float64                      `gorm:"column:item_discount_percent;null" json:"item_discount_percent"`
	ItemDiscountAmount                         float64                      `gorm:"column:item_discount_amount;null" json:"item_discount_amount"`
	ItemTotal                                  float64                      `gorm:"column:item_total;null" json:"item_total"`
	SubstituteId                               int                          `gorm:"column:substitute_id;null" json:"substitute_id"`
	PurchaseRequestSystemNumber                int                          `gorm:"column:purchase_request_system_number;null" json:"purchase_request_system_number"`
	PurchaseRequestLineNumber                  int                          `gorm:"column:purchase_request_line_number;null" json:"purchase_request_line_number"`
}

func (*PurchaseOrderDetailChangedItem) TableName() string {
	return TableNameChangesItemPurchaseOrderDetail
}
