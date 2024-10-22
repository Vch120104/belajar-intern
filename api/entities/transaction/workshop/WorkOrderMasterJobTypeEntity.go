package transactionworkshopentities

var CreateWorkOrderMasterJobTypeTable = "mtr_work_order_job_type"

type WorkOrderMasterJobType struct {
	WorkOrderJobTypeId          int    `gorm:"column:job_type_id;size:30;not null;primaryKey" json:"job_type_id"`
	WorkOrderJobTypeCode        string `gorm:"column:job_type_code;size:30;" json:"job_type_code"`
	WorkOrderJobTypeDescription string `gorm:"column:job_type_description;size:50;" json:"job_type_description"`
	IsActive                    bool   `gorm:"column:is_active;size:1;" json:"is_active"`
}

func (*WorkOrderMasterJobType) TableName() string {
	return CreateWorkOrderMasterJobTypeTable
}
