package transactionworkshopentities

var CreateWorkOrderHistoryRequestTable = "trx_work_order_history_request"

type WorkOrderHistoryRequest struct {
	WorkOrderHistoryRequestId   int    `gorm:"column:work_order_history_request_system_number;size:30;not null;primaryKey" json:"work_order_history_request_system_number"`
	WorkOrderSystemNumber       int    `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	WorkOrderDocumentNumber     int    `gorm:"column:work_order_document_number;size:30" json:"work_order_document_number"`
	WorkOrderServiceRequestLine int    `gorm:"column:work_order_service_request_line;size:30" json:"work_order_service_request_line"`
	WorkOrderServiceRequest     string `gorm:"column:work_order_service_request;size:100" json:"work_order_service_request"`
}

func (WorkOrderHistoryRequest) TableName() string {
	return CreateWorkOrderHistoryRequestTable
}
