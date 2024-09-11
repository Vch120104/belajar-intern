package transactionworkshopentities

import "time"

var CreateWorkOrderServiceTimeTable = "trx_work_order_service_time"

type WorkOrderServiceTime struct {
	IsActive                         bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CompanyId                        int       `gorm:"column:company_id;size:30;not null" json:"company_id"`
	TechnicianAllocationSystemNumber int       `gorm:"column:technician_allocation_system_number;size:30;not null" json:"technician_allocation_system_number"`
	TechnicianAllocationLine         int       `gorm:"column:technician_allocation_line;size:30;not null" json:"technician_allocation_line"`
	WorkOrderSystemNumber            int       `gorm:"column:work_order_system_number;size:30;not null" json:"work_order_system_number"`
	WorkOrderDocumentNumber          string    `gorm:"column:work_order_document_number;type:varchar(25);not null" json:"work_order_document_number"`
	WorkOrderLineId                  int       `gorm:"column:work_order_line_id;size:30;not null" json:"work_order_line_id"`
	WorkOrderDate                    time.Time `gorm:"column:work_order_date;type:datetimeoffset" json:"work_order_date"`
	OperationItemId                  int       `gorm:"column:operation_item_id;size:30;" json:"operation_item_id"`
	OperationItemCode                string    `gorm:"column:operation_item_code;size:30;" json:"operation_item_code"`
	StartDatetime                    time.Time `gorm:"column:start_datetime;type:datetimeoffset" json:"start_datetime"`
	EndDatetime                      time.Time `gorm:"column:end_datetime;type:datetimeoffset" json:"end_datetime"`
	SequenceNo                       int       `gorm:"column:sequence_no;size:30" json:"sequence_no"`
}

func (*WorkOrderServiceTime) TableName() string {
	return CreateWorkOrderServiceTimeTable
}
