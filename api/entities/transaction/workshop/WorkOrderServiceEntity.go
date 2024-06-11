package transactionworkshopentities

import "time"

const TableNameWorkOrderService = "trx_work_order_service_request"

type WorkOrderService struct {
	WorkOrderServiceId      int       `gorm:"column:work_order_service_id;size:30;primary_key" json:"work_order_service_id"`
	WorkOrderSystemNumber   int       `gorm:"column:work_order_system_number;size:30;" json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `gorm:"column:work_order_document_number;" json:"work_order_document_number"`
	WorkOrderServiceType    int       `gorm:"column:work_order_service_type;size:30;" json:"work_order_service_type"`
	WorkOrderServiceStatus  int       `gorm:"column:work_order_service_status;size:30;" json:"work_order_service_status"`
	WorkOrderServiceDate    time.Time `gorm:"column:work_order_service_date;" json:"work_order_service_date"`
	WorkOrderServiceRemark  string    `gorm:"column:work_order_service_remark;" json:"work_order_service_remark"`
}

func (*WorkOrderService) TableName() string {
	return TableNameWorkOrderService
}
