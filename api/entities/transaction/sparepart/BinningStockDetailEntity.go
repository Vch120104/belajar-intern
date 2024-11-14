package transactionsparepartentities

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const TableNameBinningStockDetail = "trx_binning_list_stock_detail"

type BinningStockDetail struct {
	BinningDetailSystemNumber       int                      `gorm:"column:binning_detail_system_number;not null;primaryKey" json:"binning_detail_system_number"`
	BinningSystemNumber             int                      `gorm:"column:binning_system_number;not null;size:30" json:"binning_system_number"`
	Binning                         *BinningStock            `gorm:"foreignKey:BinningSystemNumber;references:BinningSystemNumber"`
	BinningLineNumber               int                      `gorm:"column:binning_line_number;not null" json:"binning_line_number"`
	ReferenceDetailSystemNumber     int                      `gorm:"column:reference_detail_system_number;not null" json:"reference_detail_system_number"`
	OriginalItemId                  int                      `gorm:"column:original_item_id;null" json:"original_item_id"`
	ItemId                          int                      `gorm:"column:item_id;not null;size:30" json:"item_id"`
	Item                            *masteritementities.Item //`gorm:"foreignKey:ItemId;references:ItemId"`
	ItemPrice                       float64                  `gorm:"column:item_price;not null" json:"item_price"`
	UomId                           int                      `gorm:"column:uom_id;not null;size:30" json:"uom_id"`
	Uom                             *masteritementities.Uom
	WarehouseLocationId             int                                        `gorm:"column:warehouse_location_id;not null;size:30" json:"warehouse_location_id"`
	WarehouseLocation               *masterwarehouseentities.WarehouseLocation //`gorm:"foreignKey:WarehouseLocationId;references:WarehouseLocationId"`
	PurchaseOrderQuantity           int                                        `gorm:"column:purchase_order_quantity;not null" json:"purchase_order_quantity"`
	DeliveryOrderQuantity           float64                                    `gorm:"column:delivery_order_quantity;not null" json:"delivery_order_quantity"`
	ReferenceSystemNumber           int                                        `gorm:"column:reference_system_number;not null" json:"reference_system_number"`
	ReferenceLineNumber             int                                        `gorm:"column:reference_line_number;not null" json:"reference_line_number"`
	PurchaseOrderDetailSystemNumber int                                        `gorm:"column:purchase_order_detail_system_number" json:"purchase_order_detail_system_number"`
	GoodsReceiveSystemNumber        int                                        `gorm:"column:goods_receive_system_number;not null" json:"goods_receive_system_number"`
	GoodsReceiveDetailSystemNumber  int                                        `gorm:"column:goods_receive_detail_system_number;size:30" json:"goods_receive_detail_system_number"`
	GoodsReceiveLineNumber          int                                        `gorm:"column:goods_receive_line_number;null" json:"goods_receive_line_number"`
	SubCaseNumber                   string                                     `gorm:"column:sub_case_number;null" json:"sub_case_number"`
	CreatedByUserId                 int                                        `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate                     *time.Time                                 `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId                 int                                        `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate                     *time.Time                                 `gorm:"column:updated_date" json:"updated_date"`
	ChangeNo                        int                                        `gorm:"column:change_no;size:30;" json:"change_no"`
}

func (*BinningStockDetail) TableName() string {
	return TableNameBinningStockDetail
}
