package transactionentities

import (
	masteroperationentities "after-sales/api/entities/master/operation"
)

var CreateWorkOrderOperationTable = "work_order_operation"

type WorkOrderOperation struct {
	WorkOrderOperationId  int32                                         `gorm:"column:work_order_operation_id;not null;primaryKey"        json:"work_order_operation_id"`
	WorkOrderSystemNumber int32                                         `gorm:"column:work_order_system_number;null"        json:"work_order_system_number"`
	OperationId           int32                                         `gorm:"column:operation_id;null"        json:"operation_id"`
	Operation             masteroperationentities.OperationModelMapping `gorm:"foreignKey:OperationId;references:operation_model_mapping_id" json:"operation_model_mapping"`
	TechnicianId          int32                                         `gorm:"column:technician_id;null"        json:"technician_id"` //mtr_user_details
}

func (*WorkOrderOperation) TableName() string {
	return CreateWorkOrderOperationTable
}
