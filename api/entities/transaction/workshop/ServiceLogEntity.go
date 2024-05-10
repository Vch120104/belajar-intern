package transactionworkshopentities

import (
	transactionentities "after-sales/api/entities/transaction"
	"time"
)

var CreateServiceLogTable = "trx_service_log"

type ServiceLog struct {
	ServiceLogSystemNumber           int                                    `gorm:"column:service_log_system_number;size:30;not null;primaryKey"        json:"service_log_system_number"`
	TechnicianAllocationLine         int                                    `gorm:"column:technician_allocation_line;size:30;not null"        json:"technician_allocation_line"`
	TechnicianAllocationSystemNumber int                                    `gorm:"column:technician_allocation_system_number;size:30;not null"        json:"technician_allocation_system_number"`
	CompanyId                        int                                    `gorm:"column:company_id;size:30;not null"        json:"company_id"`
	WorkOrderSystemNumber            int                                    `gorm:"column:work_order_system_number;size:30;not null"        json:"work_order_system_number"`
	WorkOrderOperationId             int                                    `gorm:"column:work_order_operation_id;size:30;not null"        json:"work_order_operation_id"`
	WorkOrderOperation               transactionentities.WorkOrderOperation `gorm:"references:work_order_operation_id" json:"work_order_operation"`
	TechnicianId                     int                                    `gorm:"column:technician_id;size:30;not null"        json:"technician_id"`
	ShiftScheduleId                  int                                    `gorm:"column:shift_schedule_id;size:30;null"        json:"shift_schedule_id"`
	ServiceStatusId                  int                                    `gorm:"column:service_status_id;size:30;not null"        json:"service_status_id"`
	ServiceReasonId                  int                                    `gorm:"column:service_reason_id;size:30;null"        json:"service_reason_id"`
	StartDatetime                    time.Time                              `gorm:"column:start_datetime;null"        json:"start_datetime"`
	EndDatetime                      time.Time                              `gorm:"column:end_datetime;null"        json:"end_datetime"`
	ActualTime                       float64                                `gorm:"column:actual_time;not null"        json:"actual_time"`
	EstimatedPendingTime             float64                                `gorm:"column:estimated_pending_time;not null"        json:"estimated_pending_time"`
	PendingTime                      float64                                `gorm:"column:pending_time;not null"        json:"pending_time"`
	Remark                           string                                 `gorm:"column:remark;not null"        json:"remark"`
	ActualStartTime                  float64                                `gorm:"column:actual_start_time;null"        json:"actual_start_time"`
}

func (*ServiceLog) TableName() string {
	return CreateServiceLogTable
}
