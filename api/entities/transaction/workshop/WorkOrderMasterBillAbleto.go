package transactionworkshopentities

var CreateWorkOrderMasterBillAbleto = "mtr_work_order_bill_able_to"

type WorkOrderMasterBillAbleto struct {
	WorkOrderBillabletoId   int    `gorm:"column:billable_to_id;size:30;not null;primaryKey" json:"billable_to_id"`
	WorkOrderBillabletoName string `gorm:"column:billable_to_name;size:30;" json:"billable_to_name"`
	WorkOrderBillabletoCode string `gorm:"column:billable_to_code;size:30;" json:"billable_to_code"`
}

func (*WorkOrderMasterBillAbleto) TableName() string {
	return CreateWorkOrderMasterBillAbleto
}
