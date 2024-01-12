package transactionentities

import (
	masteritementities "after-sales/api/entities/master/item"
)

var CreateWorkOrderItemTable = "work_order_item"

type WorkOrderItem struct {
	IsActive              bool                    `gorm:"column:is_active;not null;default:true" json:"is_active"`
	WorkOrderItemId       int32                   `gorm:"column:work_order_item_id;not null;primaryKey"        json:"work_order_item_id"`
	WorkOrderSystemNumber int32                   `gorm:"column:work_order_system_number;null"        json:"work_order_system_number"`
	ItemIds               int32                   `gorm:"column:item_id;null"        json:"item_id"`
	Item                  masteritementities.Item `gorm:"foreignKey:ItemIds;references:item_id" json:"item"`
	TechnicianId          int32                   `gorm:"column:technician_id;null"        json:"technician_id"` //mtr_user_details
}

func (*WorkOrderItem) TableName() string {
	return CreateWorkOrderItemTable
}
