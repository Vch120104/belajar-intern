package transactionworkshopentities

var CreateWorkOrderMasterStatusTable = "mtr_work_order_status"

type WorkOrderMasterStatus struct {
	WorkOrderStatusId   int    `gorm:"column:work_order_status_id;size:30;not null;primaryKey" json:"work_order_status_id"`
	WorkOrderStatusCode string `gorm:"column:work_order_status_code;size:30;" json:"work_order_status_code"`
	WorkOrderStatusDesc string `gorm:"column:work_order_status_description;size:30;" json:"work_order_status_description"`
}

func (*WorkOrderMasterStatus) TableName() string {
	return CreateWorkOrderMasterStatusTable
}
