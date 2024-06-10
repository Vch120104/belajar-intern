package transactionworkshopentities

import "time"

const TableNameWorkOrderServiceVehicle = "trx_work_order_service_vehicle"

type WorkOrderServiceVehicle struct {
	WorkOrderServiceVehicleId int       `gorm:"column:work_order_service_vehicle_id;size:30;primary_key" json:"work_order_service_vehicle_id"`
	WorkOrderDocumentNumber   string    `gorm:"column:work_order_document_number;" json:"work_order_document_number"`
	WorkOrderSystemNumber     int       `gorm:"column:work_order_system_number;size:30;" json:"work_order_system_number"`
	WorkOrderServiceId        int       `gorm:"column:work_order_service_id;size:30;" json:"work_order_service_id"`
	WorkOrderServiceDate      time.Time `gorm:"column:work_order_service_date;" json:"work_order_service_date"`
	WorkOrderServiceRemark    string    `gorm:"column:work_order_service_remark;" json:"work_order_service_remark"`
}

func (*WorkOrderServiceVehicle) TableName() string {
	return TableNameWorkOrderServiceVehicle
}
