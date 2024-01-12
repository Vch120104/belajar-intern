package transactionsparepartentities

import transactionentities "after-sales/api/entities/transaction"

var CreateSupplySlipDetailTable = "trx_supply_slip_detail"

type SupplySlipDetail struct {
	IsActive                          bool                              `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SupplySlipDetailSystemNumber      int32                             `gorm:"column:supply_slip_detail_system_number;not null;primaryKey"        json:"supply_slip_detail_system_number"`
	SupplySystemNumbers               int32                             `gorm:"column:supply_system_number;not null"        json:"supply_system_number"`
	SupplySlip                        SupplySlip                        `gorm:"foreignKey:SupplySystemNumbers;references:supply_system_number" json:"supply_slip"`
	SupplySystemLineNumber            int32                             `gorm:"column:supply_system_line_number;not null"        json:"supply_system_line_number"`
	LocationId                        int32                             `gorm:"column:location_id;null"        json:"location_id"`
	WorkOrderItemId                   int32                             `gorm:"column:work_order_item_id;null"        json:"work_order_item_id"`
	WorkOrderItem                     transactionentities.WorkOrderItem `gorm:"references:work_order_item_id" json:"work_order_item"`
	UnitOutMeasurementId              int32                             `gorm:"column:unit_out_measurement_id;null"        json:"unit_out_measurement_id"`
	QuantitySupply                    float64                           `gorm:"column:quantity_supply;null"        json:"quantity_supply"`
	QuantityReturn                    float64                           `gorm:"column:quantity_return;null"        json:"quantity_return"`
	QuantityDemand                    float64                           `gorm:"column:quantity_demand;null"        json:"quantity_demand"`
	CostOfGoodsSold                   float64                           `gorm:"column:cost_of_goods_sold;null"        json:"cost_of_goods_sold"`
	PurchaseRequestSystemNumber       int32                             `gorm:"column:purchase_request_system_number;null"        json:"purchase_request_system_number"`
	PurchaseRequestSystemNumberDetail int32                             `gorm:"column:purchase_request_system_number_detail;null"        json:"purchase_request_system_number_detail"`
	WorkOrderSystemNumber             int32                             `gorm:"column:work_order_system_number;null"        json:"work_order_system_number"`
	WorkOrderLineNumberId             int32                             `gorm:"column:work_order_line_number_id;null"        json:"work_order_line_number_id"`
	WarehouseGroupId                  int32                             `gorm:"column:warehouse_group_id;null"        json:"warehouse_group_id"`
	WarehouseId                       int32                             `gorm:"column:warehouse_id;null"        json:"warehouse_id"`
	QuantityTotal                     int32                             `gorm:"column:quantity_total;null"        json:"quantity_total"`
}

func (*SupplySlipDetail) TableName() string {
	return CreateSupplySlipDetailTable
}
