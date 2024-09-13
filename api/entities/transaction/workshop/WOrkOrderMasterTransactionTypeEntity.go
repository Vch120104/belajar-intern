package transactionworkshopentities

var CreateWorkOrderMasterTrxTypeTable = "mtr_work_order_transaction_type"

type WorkOrderMasterTrxType struct {
	WorkOrderTrxTypeId          int    `gorm:"column:work_order_transaction_type_id;size:30;not null;primaryKey" json:"work_order_transaction_type_id"`
	WorkOrderTrxTypeCode        string `gorm:"column:work_order_transaction_type_code;size:30;" json:"work_order_transaction_type_code"`
	WorkOrderTrxTypeDescription string `gorm:"column:work_order_transaction_type_description;size:50;" json:"work_order_transaction_type_description"`
	WorkOrderTrxTypePrefix      string `gorm:"column:work_order_transaction_type_prefix;size:30;" json:"work_order_transaction_type_prefix"`
	IsActive                    bool   `gorm:"column:is_active;size:1;" json:"is_active"`
}

func (*WorkOrderMasterTrxType) TableName() string {
	return CreateWorkOrderMasterTrxTypeTable
}
