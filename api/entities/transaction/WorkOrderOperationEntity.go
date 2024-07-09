package transactionentities

import (
	masteroperationentities "after-sales/api/entities/master/operation"
)

var CreateWorkOrderOperationTable = "work_order_operation"

type WorkOrderOperation struct {
	WorkOrderOperationId  int                                           `gorm:"column:work_order_operation_id;size:30;not null;primaryKey"        json:"work_order_operation_id"`
	WorkOrderSystemNumber int                                           `gorm:"column:work_order_system_number;null;size:30;"        json:"work_order_system_number"`
	OperationId           int                                           `gorm:"column:operation_id;null;size:30;"        json:"operation_id"`
	Operation             masteroperationentities.OperationModelMapping `gorm:"foreignKey:OperationId;references:operation_model_mapping_id" json:"operation_model_mapping"`
	TechnicianId          int                                           `gorm:"column:technician_id;null;size:30;"        json:"technician_id"` //mtr_user_details
}

func (*WorkOrderOperation) TableName() string {
	return CreateWorkOrderOperationTable
}
