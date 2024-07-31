package transactionworkshopentities

const TableNameWorkOrderRequestDescription = "trx_work_order_request_description"

type WorkOrderRequestDescription struct {
	WorkOrderRequestId          int    `gorm:"column:work_order_request_id;size:30;primaryKey" json:"work_order_request_id"`
	WorkOrderSystemNumber       int    `gorm:"column:work_order_system_number;size:30;" json:"work_order_system_number"`
	WorkOrderServiceRequestLine int    `gorm:"column:work_order_service_request_line;size:30;" json:"work_order_service_request_line"`
	WorkOrderServiceRequest     string `gorm:"column:work_order_service_request;default:null;size:50;" json:"work_order_service_request"`
}

func (*WorkOrderRequestDescription) TableName() string {
	return TableNameWorkOrderRequestDescription
}
