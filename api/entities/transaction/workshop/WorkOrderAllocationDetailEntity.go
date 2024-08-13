package transactionworkshopentities

import "time"

const TableNameWorkOrderAllocationDetail = "trx_work_order_allocation_detail"

type WorkOrderAllocationDetail struct {
	TechnicianId            int       `gorm:"column:technician_id;size:30" json:"technician_id"`
	TechnicianName          string    `gorm:"column:technician_name;" json:"technician_name"`
	WorkOrderSystemNumber   int       `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `gorm:"column:work_order_document_number;" json:"work_order_document_number"`
	ShiftCode               string    `gorm:"column:shift_code;" json:"shift_code"`
	ServiceStatus           string    `gorm:"column:service_status;" json:"service_status"`
	StartTime               time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime                 time.Time `gorm:"column:end_time" json:"end_time"`
	ServiceLogId            int       `gorm:"column:service_log_system_number" json:"service_log_system_number"`
}

func (*WorkOrderAllocationDetail) TableName() string {
	return TableNameWorkOrderAllocationDetail
}
