package transactionsparepartentities

import transactionentities "after-sales/api/entities/transaction"

var CreateSupplySlipDetailTable = "trx_supply_slip_detail"

type SupplySlipDetail struct {
	IsActive                          bool                              `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SupplySlipDetailSystemNumber      int                               `gorm:"column:supply_slip_detail_system_number;not null;primaryKey;size:30" json:"supply_slip_detail_system_number"`
	SupplySystemNumbers               int                               `gorm:"column:supply_system_number;not null" json:"supply_system_number"`
	SupplySlip                        SupplySlip                        `gorm:"foreignKey:SupplySystemNumbers;references:SupplySystemNumber" json:"supply_slip"`
	SupplySystemLineNumber            int                               `gorm:"column:supply_system_line_number;not null" json:"supply_system_line_number"`
	LocationId                        int                               `gorm:"column:location_id" json:"location_id"`
	WorkOrderItemId                   int                               `gorm:"column:work_order_item_id" json:"work_order_item_id"`
	WorkOrderItem                     transactionentities.WorkOrderItem `gorm:"references:WorkOrderItemId" json:"work_order_item"`
	UnitOutMeasurementId              int                               `gorm:"column:unit_out_measurement_id" json:"unit_out_measurement_id"`
	QuantitySupply                    float32                           `gorm:"column:quantity_supply" json:"quantity_supply"`
	QuantityReturn                    float32                           `gorm:"column:quantity_return" json:"quantity_return"`
	QuantityDemand                    float32                           `gorm:"column:quantity_demand" json:"quantity_demand"`
	CostOfGoodsSold                   float32                           `gorm:"column:cost_of_goods_sold" json:"cost_of_goods_sold"`
	PurchaseRequestSystemNumber       int                               `gorm:"column:purchase_request_system_number" json:"purchase_request_system_number"`
	PurchaseRequestSystemNumberDetail int                               `gorm:"column:purchase_request_system_number_detail" json:"purchase_request_system_number_detail"`
	WorkOrderSystemNumber             int                               `gorm:"column:work_order_system_number" json:"work_order_system_number"`
	WorkOrderLineNumberId             int                               `gorm:"column:work_order_line_number_id" json:"work_order_line_number_id"`
	WarehouseGroupId                  int                               `gorm:"column:warehouse_group_id" json:"warehouse_group_id"`
	WarehouseId                       int                               `gorm:"column:warehouse_id" json:"warehouse_id"`
	QuantityTotal                     int                               `gorm:"column:quantity_total" json:"quantity_total"`
}

func (*SupplySlipDetail) TableName() string {
	return CreateSupplySlipDetailTable
}
