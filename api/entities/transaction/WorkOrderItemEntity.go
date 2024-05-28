package transactionentities

import (
	masteritementities "after-sales/api/entities/master/item"
)

var CreateWorkOrderItemTable = "work_order_item"

type WorkOrderItem struct {
	IsActive              bool                    `gorm:"column:is_active;not null;default:true" json:"is_active"`
	WorkOrderItemId       int                     `gorm:"column:work_order_item_id;size:30;not null;primaryKey"        json:"work_order_item_id"`
	WorkOrderSystemNumber int                     `gorm:"column:work_order_system_number;size:30;null"        json:"work_order_system_number"`
	ItemIds               int                     `gorm:"column:item_id;null;size:30;"        json:"item_id"`
	Item                  masteritementities.Item `gorm:"foreignKey:ItemIds;references:item_id" json:"item"`
	TechnicianId          int                     `gorm:"column:technician_id;null;size:30;"        json:"technician_id"` //mtr_user_details
}

func (*WorkOrderItem) TableName() string {
	return CreateWorkOrderItemTable
}
