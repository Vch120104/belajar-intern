package transactionworkshopentities

var CreateWorkOrderMasterTable = "mtr_work_order_site_type"

type WorkOrderMaster struct {
	WorkOrderTypeId   int    `gorm:"column:work_order_site_type_id;size:30;not null;primaryKey" json:"work_order_site_type_id"`
	WorkOrderTypeCode string `gorm:"column:work_order_site_type_code;size:30;" json:"work_order_site_type_code"`
	WorkOrderTypeDesc string `gorm:"column:work_order_site_type_description;size:30;" json:"work_order_site_type_description"`
}

func (*WorkOrderMaster) TableName() string {
	return CreateWorkOrderMasterTable
}
