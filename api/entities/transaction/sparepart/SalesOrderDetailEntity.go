package transactionsparepartentities

import (
	masteritementities "after-sales/api/entities/master/item"
	"time"
)

const TableNameSalesOrderDetail = "trx_sales_order_detail"

type SalesOrderDetail struct {
	SalesOrderDetailSystemNumber                       int                                `gorm:"column:sales_order_detail_system_number;size:30;primaryKey" json:"sales_order_detail_system_number"`
	SalesOrderSystemNumber                             int                                `gorm:"column:sales_order_system_number;size:30;not null" json:"sales_order_system_number"`
	SalesOrder                                         SalesOrder                         `gorm:"foreignKey:SalesOrderSystemNumber;references:SalesOrderSystemNumber"`
	SalesOrderLineNumber                               int                                `gorm:"column:sales_order_line_number;size:30;not null" json:"sales_order_line_number"`
	SalesOrderLineStatusId                             *int                               `gorm:"column:sales_order_line_status_id;size:30" json:"sales_order_line_status_id"` // FK to mtr_approval_status in general-service
	ItemSubstituteId                                   *int                               `gorm:"column:item_substitute_id;size:30" json:"item_substitute_id"`
	SubstituteItem                                     *masteritementities.ItemSubstitute `gorm:"foreignKey:ItemSubstituteId;references:ItemSubstituteId"`
	ItemId                                             int                                `gorm:"column:item_id;size:30;not null" json:"item_id"`
	Item                                               masteritementities.Item            `gorm:"foreignKey:ItemId;refenreces:ItemId"`
	QuantityDemand                                     float64                            `gorm:"column:quantity_demand;not null" json:"quantity_demand"`
	IsAvailable                                        bool                               `gorm:"column:is_available;not null" json:"is_available"`
	QuantitySupply                                     float64                            `gorm:"column:quantity_supply;not null" json:"quantity_supply"`
	QuantityPick                                       float64                            `gorm:"column:quantity_pick;not null" json:"quantity_pick"`
	UomId                                              *int                               `gorm:"column:uom_id;size:30;not null" json:"uom_id"`
	Uom                                                *masteritementities.Uom            `gorm:"foreignKey:UomId;references:UomId"`
	Price                                              float64                            `gorm:"column:price;not null" json:"price"`
	PriceEffectiveDate                                 *time.Time                         `gorm:"column:price_effective_date;type:datetime" json:"price_effective_date"`
	DiscountPercent                                    float64                            `gorm:"column:discount_percent;not null" json:"discount_percent"`
	DiscountAmount                                     float64                            `gorm:"column:discount_amount;not null" json:"discount_amount"`
	DiscountRequestPercent                             float64                            `gorm:"column:discount_request_percent;not null" json:"discount_request_percent"`
	DiscountRequestAmount                              float64                            `gorm:"column:discount_request_amount;not null" json:"discount_request_amount"`
	Remark                                             string                             `gorm:"column:remark;size:256	;not null" json:"remark"`
	ApprovalRequestNumber                              *int                               `gorm:"column:approval_reqeust_number;size:30" json:"approval_request_number"` // FK to trx_approval_request_source in ?
	ApprovalRemark                                     string                             `gorm:"column:approval_remark;size:256;not null" json:"approval_remark"`
	VehicleSalesOrderSystemNumber                      *int                               `gorm:"column:vehicle_sales_order_system_number;size:30" json:"vehicle_sales_order_system_number"`                                               // FK to trx_vehicle_sales_order in sales-service
	VehicleSalesOrderDetailSystemNumber                *int                               `gorm:"column:vehicle_sales_order_detail_system_number;size:30" json:"vehicle_sales_order_detail_system_number"`                                 // FK to trx_vehicle_sales_order_detail in sales-service
	VehicleSalesOrderAdditionalAccessoriesSystemNumber *int                               `gorm:"column:vehicle_sales_order_additional_accessories_system_number;size:30" json:"vehicle_sales_order_additional_accessories_system_number"` // FK to trx_vehicle_sales_order_additional_accessories in sales-service
	ItemPackageId                                      *int                               `gorm:"column:item_package_id;size:30;unique" json:"item_package_id"`
	ItemPackage                                        *masteritementities.ItemPackage    `gorm:"foreignKey:ItemPackageId;references:ItemPackageId"`
	PurchaseOrderSystemNumber                          int                                `gorm:"column:purchase_order_system_number;size:30;not null" json:"purchase_order_system_number"` // FK to trx_unit_purchase_order in sales-service
	ReferenceItemId                                    *int                               `gorm:"column:reference_item_id;size:30" json:"reference_item_id"`
	PriceListId                                        *int                               `gorm:"column:price_list_id;size:30;not null" json:"price_list_id"`
	PriceList                                          *masteritementities.ItemPriceList  `gorm:"foreignKey:PriceListId;references:PriceListId"`
	CurrencyId                                         *int                               `gorm:"column:currency_id;size:30" json:"currency_id"`
	CurrencyAmount                                     *float64                           `gorm:"column:currency_amount" json:"currency_amount"`
}

func (*SalesOrderDetail) TableName() string {
	return TableNameSalesOrderDetail
}
