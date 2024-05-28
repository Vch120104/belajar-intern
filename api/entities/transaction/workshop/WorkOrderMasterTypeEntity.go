package transactionworkshopentities

var CreateWorkOrderMasterTypeTable = "mtr_work_order_type"

type WorkOrderMasterType struct {
	WorkOrderTypeId   int    `gorm:"column:work_order_type_id;size:30;not null;primaryKey" json:"work_order_type_id"`
	WorkOrderTypeCode string `gorm:"column:work_order_type_code;size:30;" json:"work_order_type_code"`
	WorkOrderTypeDesc string `gorm:"column:work_order_type_description;size:30;" json:"work_order_type_description"`
	IsActive          bool   `gorm:"column:is_active;size:1;" json:"is_active"`
}

func (*WorkOrderMasterType) TableName() string {
	return CreateWorkOrderMasterTypeTable
}
