package transactionworkshopentities

import "time"

var CreateWorkOrderQualityControlTable = "trx_work_order_quality_control"

type WorkOrderQualityControl struct {
	WorkOrderQualityControlId       int       `gorm:"column:work_order_quality_control_system_number;size:30;not null;primaryKey" json:"work_order_quality_control_system_number"`
	WorkOrderSystemNumber           int       `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	WorkOrderQualityControlStatusID int       `gorm:"column:work_order_quality_control_status_id;size:30" json:"work_order_quality_control_status_id"`
	WorkOrderStartDateTime          time.Time `gorm:"column:work_order_start_date_time;type:datetimeoffset" json:"work_order_start_date_time"`
	WorkOrderEndDateTime            time.Time `gorm:"column:work_order_end_date_time;type:datetimeoffset" json:"work_order_end_date_time"`
	WorkOrderActualTime             float32   `gorm:"column:work_order_actual_time;" json:"work_order_actual_time"`
}

func (WorkOrderQualityControl) TableName() string {
	return CreateWorkOrderQualityControlTable
}
