package transactionworkshopentities

var CreateWorkOrderMasterLineTypeTable = "mtr_work_order_line_type"

type WorkOrderMasterLineType struct {
	WorkOrderLineTypeId          int    `gorm:"column:line_type_id;size:30;not null;primaryKey" json:"line_type_id"`
	WorkOrderLineTypeCode        string `gorm:"column:line_type_code;size:30;" json:"line_type_code"`
	WorkOrderLineTypeDescription string `gorm:"column:line_type_description;size:50;" json:"line_type_description"`
	IsActive                     bool   `gorm:"column:is_active;size:1;" json:"is_active"`
}

func (*WorkOrderMasterLineType) TableName() string {
	return CreateWorkOrderMasterLineTypeTable
}
