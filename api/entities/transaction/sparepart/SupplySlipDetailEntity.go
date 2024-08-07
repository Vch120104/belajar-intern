package transactionsparepartentities

import (
	masteritementities "after-sales/api/entities/master/item"
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	transactionentities "after-sales/api/entities/transaction"
)

var CreateSupplySlipDetailTable = "trx_supply_slip_detail"

type SupplySlipDetail struct {
	IsActive                          bool                                      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SupplyDetailSystemNumber          int                                       `gorm:"column:supply_detail_system_number;;size:30;not null;primaryKey;size:30" json:"supply_detail_system_number"`
	SupplySystemNumber                int                                       `gorm:"column:supply_system_number;not null;size:30;" json:"supply_system_number"`
	Supply                            SupplySlip                                `gorm:"foreignKey:SupplySystemNumber;references:SupplySystemNumber"`
	SupplySystemLineNumber            int                                       `gorm:"column:supply_system_line_number;null;size:30;" json:"supply_system_line_number"`
	LocationId                        int                                       `gorm:"column:location_id;size:30;" json:"location_id"`
	Location                          masterwarehouseentities.WarehouseLocation `gorm:"foreignKey:LocationId;references:WarehouseLocationId"`
	WorkOrderOperationId              int                                       `gorm:"column:work_order_operation_id;size:30;" json:"work_order_operation_id"`
	WorkOrderOperation                transactionentities.WorkOrderOperation    `gorm:"foreignKey:WorkOrderOperationId;references:WorkOrderOperationId"`
	WorkOrderItemId                   int                                       `gorm:"column:work_order_item_id;size:30;" json:"work_order_item_id"`
	WorkOrderItem                     transactionentities.WorkOrderItem         `gorm:"foreignKey:WorkOrderItemId;references:WorkOrderItemId"`
	UnitOfMeasurementId               int                                       `gorm:"column:unit_of_measurement_id;size:30;" json:"unit_of_measurement_id"`
	UnitOfMeasurement                 masteritementities.Uom                    `gorm:"foreignKey:UnitOfMeasurementId;references:UomId"`
	QuantitySupply                    float32                                   `gorm:"column:quantity_supply" json:"quantity_supply"`
	QuantityReturn                    float32                                   `gorm:"column:quantity_return" json:"quantity_return"`
	QuantityDemand                    float32                                   `gorm:"column:quantity_demand" json:"quantity_demand"`
	CostOfGoodsSold                   float32                                   `gorm:"column:cost_of_goods_sold" json:"cost_of_goods_sold"`
	PurchaseRequestSystemNumber       int                                       `gorm:"column:purchase_request_system_number;size:30;" json:"purchase_request_system_number"`
	PurchaseRequestLineNumber         int                                       `gorm:"column:purchase_request_line_number;size:30;" json:"purchase_request_line_number"`
	PurchaseRequestDetailSystemNumber int                                       `gorm:"column:purchase_request_detail_system_number;size:50;" json:"purchase_request_detail_system_number"`
	PurchaseRequestDetail             PurchaseRequestDetail                     `gorm:"foreignKey:PurchaseRequestDetailSystemNumber;references:PurchaseRequestDetailSystemNumber"`
	WorkOrderSystemNumber             int                                       `gorm:"column:work_order_system_number;size:30;" json:"work_order_system_number"`
	WorkOrderLineNumber               int                                       `gorm:"column:work_order_line_number;size:30;" json:"work_order_line_number"`
	WorkOrderDetailId                 int                                       `gorm:"column:work_order_detail_id;size:30;" json:"work_order_detail_id"`
	WarehouseGroupId                  int                                       `gorm:"column:warehouse_group_id;size:30;" json:"warehouse_group_id"`
	WarehouseGroup                    masterwarehouseentities.WarehouseGroup    `gorm:"foreignKey:WarehouseGroupId;references:WarehouseGroupId"`
	WarehouseId                       int                                       `gorm:"column:warehouse_id;size:30;" json:"warehouse_id"`
	// Warehouse                         masterwarehouseentities.WarehouseMaster   `gorm:"foreignKey:WarehouseId;references:WarehouseWarehouseId"`
	QuantityTotal                     float32                                   `gorm:"column:quantity_total;size;" json:"quantity_total"`
}

func (*SupplySlipDetail) TableName() string {
	return CreateSupplySlipDetailTable
}
